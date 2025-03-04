// Copyright 2025 R5
// This file is part of the R5 Core library.
//
// This software is provided "as is", without warranty of any kind,
// express or implied, including but not limited to the warranties
// of merchantability, fitness for a particular purpose and
// noninfringement. In no event shall the authors or copyright
// holders be liable for any claim, damages, or other liability,
// whether in an action of contract, tort or otherwise, arising
// from, out of or in connection with the software or the use or
// other dealings in the software.

package ethash

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"
	"runtime"
	"time"

	"github.com/r5-codebase/r5-core/common"
	"github.com/r5-codebase/r5-core/common/math"
	"github.com/r5-codebase/r5-core/consensus"
	"github.com/r5-codebase/r5-core/consensus/misc"
	"github.com/r5-codebase/r5-core/core/state"
	"github.com/r5-codebase/r5-core/core/types"
	"github.com/r5-codebase/r5-core/params"
	"github.com/r5-codebase/r5-core/rlp"
	"github.com/r5-codebase/r5-core/trie"
	"golang.org/x/crypto/sha3"
)

// R5 proof-of-work protocol parameters.
var (
	// Maximum number of uncle blocks allowed per block.
	maxUncles = 2

	// Allowed future block time updated to 7 seconds.
	allowedFutureBlockTimeSeconds = int64(7)

	// Difficulty adjustment calculators (tuned for a 7-second target, using division by 2).
	calcDifficultyEip5133        = makeDifficultyCalculator()
	calcDifficultyEip4345        = makeDifficultyCalculator()
	calcDifficultyEip3554        = makeDifficultyCalculator()
	calcDifficultyEip2384        = makeDifficultyCalculator()
	calcDifficultyConstantinople = makeDifficultyCalculator()
	calcDifficultyByzantium      = makeDifficultyCalculator()
)

// Error values used for block validation.
var (
	errOlderBlockTime    = errors.New("timestamp older than parent")
	errTooManyUncles     = errors.New("too many uncles")
	errInvalidDifficulty = errors.New("non-positive difficulty")
	errInvalidMixDigest  = errors.New("invalid mix digest")
	errInvalidPoW        = errors.New("invalid proof-of-work")
)

// Author implements consensus.Engine, returning the header's coinbase as the
// proof-of-work verified author of the block.
func (r5 *Ethash) Author(header *types.Header) (common.Address, error) {
	return header.Coinbase, nil
}

// VerifyHeader checks whether a header conforms to the consensus rules of the
// modified R5 consensus engine.
func (r5 *Ethash) VerifyHeader(chain consensus.ChainHeaderReader, header *types.Header, seal bool) error {
	if r5.config.PowMode == ModeFullFake {
		return nil
	}
	number := header.Number.Uint64()
	if chain.GetHeader(header.Hash(), number) != nil {
		return nil
	}
	parent := chain.GetHeader(header.ParentHash, number-1)
	if parent == nil {
		return consensus.ErrUnknownAncestor
	}
	return r5.verifyHeader(chain, header, parent, false, seal, time.Now().Unix())
}

// VerifyHeaders concurrently verifies a batch of headers.
func (r5 *Ethash) VerifyHeaders(chain consensus.ChainHeaderReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error) {
	if r5.config.PowMode == ModeFullFake || len(headers) == 0 {
		abort, results := make(chan struct{}), make(chan error, len(headers))
		for i := 0; i < len(headers); i++ {
			results <- nil
		}
		return abort, results
	}

	workers := runtime.GOMAXPROCS(0)
	if len(headers) < workers {
		workers = len(headers)
	}

	var (
		inputs  = make(chan int)
		done    = make(chan int, workers)
		errors  = make([]error, len(headers))
		abort   = make(chan struct{})
		unixNow = time.Now().Unix()
	)
	for i := 0; i < workers; i++ {
		go func() {
			for index := range inputs {
				errors[index] = r5.verifyHeaderWorker(chain, headers, seals, index, unixNow)
				done <- index
			}
		}()
	}

	errorsOut := make(chan error, len(headers))
	go func() {
		defer close(inputs)
		var (
			in, out = 0, 0
			checked = make([]bool, len(headers))
			inputs  = inputs
		)
		for {
			select {
			case inputs <- in:
				if in++; in == len(headers) {
					inputs = nil
				}
			case index := <-done:
				for checked[index] = true; checked[out]; out++ {
					errorsOut <- errors[out]
					if out == len(headers)-1 {
						return
					}
				}
			case <-abort:
				return
			}
		}
	}()
	return abort, errorsOut
}

func (r5 *Ethash) verifyHeaderWorker(chain consensus.ChainHeaderReader, headers []*types.Header, seals []bool, index int, unixNow int64) error {
	var parent *types.Header
	if index == 0 {
		parent = chain.GetHeader(headers[0].ParentHash, headers[0].Number.Uint64()-1)
	} else if headers[index-1].Hash() == headers[index].ParentHash {
		parent = headers[index-1]
	}
	if parent == nil {
		return consensus.ErrUnknownAncestor
	}
	return r5.verifyHeader(chain, headers[index], parent, false, seals[index], unixNow)
}

// VerifyUncles checks that the block's uncle headers conform to R5 consensus rules.
// Although uncle blocks are allowed (up to 2 per block), no rewards are granted for them.
func (r5 *Ethash) VerifyUncles(chain consensus.ChainReader, block *types.Block) error {
	if len(block.Uncles()) > maxUncles {
		return errTooManyUncles
	}
	return nil
}

// verifyHeader validates a block header against R5 consensus rules.
func (r5 *Ethash) verifyHeader(chain consensus.ChainHeaderReader, header, parent *types.Header, uncle bool, seal bool, unixNow int64) error {
	if uint64(len(header.Extra)) > params.MaximumExtraDataSize {
		return fmt.Errorf("extra-data too long: %d > %d", len(header.Extra), params.MaximumExtraDataSize)
	}
	if !uncle {
		if header.Time > uint64(unixNow+allowedFutureBlockTimeSeconds) {
			return consensus.ErrFutureBlock
		}
	}
	if header.Time <= parent.Time {
		return errOlderBlockTime
	}
	expected := r5.CalcDifficulty(chain, header.Time, parent)
	if expected.Cmp(header.Difficulty) != 0 {
		return fmt.Errorf("invalid difficulty: have %v, want %v", header.Difficulty, expected)
	}
	if header.GasLimit > params.MaxGasLimit {
		return fmt.Errorf("invalid gasLimit: have %v, max %v", header.GasLimit, params.MaxGasLimit)
	}
	if header.GasUsed > header.GasLimit {
		return fmt.Errorf("invalid gasUsed: have %d, gasLimit %d", header.GasUsed, header.GasLimit)
	}
	if !chain.Config().IsLondon(header.Number) {
		if header.BaseFee != nil {
			return fmt.Errorf("invalid baseFee before fork: have %d, expected 'nil'", header.BaseFee)
		}
		if err := misc.VerifyGaslimit(parent.GasLimit, header.GasLimit); err != nil {
			return err
		}
	} else if err := misc.VerifyEip1559Header(chain.Config(), parent, header); err != nil {
		return err
	}
	if diff := new(big.Int).Sub(header.Number, parent.Number); diff.Cmp(big.NewInt(1)) != 0 {
		return consensus.ErrInvalidNumber
	}
	if chain.Config().IsShanghai(header.Time) {
		return fmt.Errorf("R5 does not support shanghai fork")
	}
	if chain.Config().IsCancun(header.Time) {
		return fmt.Errorf("R5 does not support cancun fork")
	}
	if seal {
		if err := r5.verifySeal(chain, header, false); err != nil {
			return err
		}
	}
	if err := misc.VerifyDAOHeaderExtraData(chain.Config(), header); err != nil {
		return err
	}
	return nil
}

// CalcDifficulty returns the difficulty that a new block should have.
func (r5 *Ethash) CalcDifficulty(chain consensus.ChainHeaderReader, time uint64, parent *types.Header) *big.Int {
	return CalcDifficulty(chain.Config(), time, parent)
}

// CalcDifficulty computes the new difficulty based on the parent's difficulty and timestamp.
func CalcDifficulty(config *params.ChainConfig, time uint64, parent *types.Header) *big.Int {
	next := new(big.Int).Add(parent.Number, big1)
	switch {
	case config.IsGrayGlacier(next):
		return calcDifficultyEip5133(time, parent)
	case config.IsArrowGlacier(next):
		return calcDifficultyEip4345(time, parent)
	case config.IsLondon(next):
		return calcDifficultyEip3554(time, parent)
	case config.IsMuirGlacier(next):
		return calcDifficultyEip2384(time, parent)
	case config.IsConstantinople(next):
		return calcDifficultyConstantinople(time, parent)
	case config.IsByzantium(next):
		return calcDifficultyByzantium(time, parent)
	case config.IsHomestead(next):
		return calcDifficultyHomestead(time, parent)
	default:
		return calcDifficultyFrontier(time, parent)
	}
}

// Consolidated constants to avoid repeated allocations.
var (
	expDiffPeriod = big.NewInt(100000)
	big1          = big.NewInt(1)
	big2          = big.NewInt(2)
	bigMinus99    = big.NewInt(-99)
)

// makeDifficultyCalculator creates a difficulty calculator tuned for a ~7‑second target.
// It uses a divisor of 2 instead of 3.
func makeDifficultyCalculator() func(time uint64, parent *types.Header) *big.Int {
	// Set target divisor to 2 for a ~7-second block target.
	bigTargetDivisor := big.NewInt(2)
	return func(time uint64, parent *types.Header) *big.Int {
		bigTime := new(big.Int).SetUint64(time)
		bigParentTime := new(big.Int).SetUint64(parent.Time)

		x := new(big.Int)
		y := new(big.Int)
		x.Sub(bigTime, bigParentTime)
		x.Div(x, bigTargetDivisor)
		if parent.UncleHash == types.EmptyUncleHash {
			x.Sub(big1, x)
		} else {
			x.Sub(big2, x)
		}
		if x.Cmp(bigMinus99) < 0 {
			x.Set(bigMinus99)
		}
		y.Div(parent.Difficulty, params.DifficultyBoundDivisor)
		x.Mul(y, x)
		x.Add(parent.Difficulty, x)
		if x.Cmp(params.MinimumDifficulty) < 0 {
			x.Set(params.MinimumDifficulty)
		}
		return x
	}
}

// calcDifficultyHomestead computes the difficulty using Homestead rules tuned for a ~7‑second target.
func calcDifficultyHomestead(time uint64, parent *types.Header) *big.Int {
	bigTime := new(big.Int).SetUint64(time)
	bigParentTime := new(big.Int).SetUint64(parent.Time)

	// Use a divisor of 2 instead of 10.
	bigTargetDivisor := big.NewInt(2)

	x := new(big.Int)
	y := new(big.Int)
	x.Sub(bigTime, bigParentTime)
	x.Div(x, bigTargetDivisor)
	x.Sub(big1, x)
	if x.Cmp(bigMinus99) < 0 {
		x.Set(bigMinus99)
	}
	y.Div(parent.Difficulty, params.DifficultyBoundDivisor)
	x.Mul(y, x)
	x.Add(parent.Difficulty, x)
	if x.Cmp(params.MinimumDifficulty) < 0 {
		x.Set(params.MinimumDifficulty)
	}
	periodCount := new(big.Int).Add(parent.Number, big1)
	periodCount.Div(periodCount, expDiffPeriod)
	if periodCount.Cmp(big1) > 0 {
		y.Sub(periodCount, big2)
		y.Exp(big2, y, nil)
		x.Add(x, y)
	}
	return x
}

// calcDifficultyFrontier computes the difficulty using Frontier rules tuned for a ~7‑second target.
func calcDifficultyFrontier(time uint64, parent *types.Header) *big.Int {
	diff := new(big.Int)
	adjust := new(big.Int).Div(parent.Difficulty, params.DifficultyBoundDivisor)
	bigTime := new(big.Int).SetUint64(time)
	bigParentTime := new(big.Int).SetUint64(parent.Time)

	// Use a hardcoded 7-second target for comparison.
	targetBlockTime := big.NewInt(7)
	if bigTime.Sub(bigTime, bigParentTime).Cmp(targetBlockTime) < 0 {
		diff.Add(parent.Difficulty, adjust)
	} else {
		diff.Sub(parent.Difficulty, adjust)
	}
	if diff.Cmp(params.MinimumDifficulty) < 0 {
		diff.Set(params.MinimumDifficulty)
	}
	periodCount := new(big.Int).Add(parent.Number, big1)
	periodCount.Div(periodCount, expDiffPeriod)
	if periodCount.Cmp(big1) > 0 {
		expDiff := periodCount.Sub(periodCount, big2)
		expDiff.Exp(big2, expDiff, nil)
		diff.Add(diff, expDiff)
		diff = math.BigMax(diff, params.MinimumDifficulty)
	}
	return diff
}

// Exported for fuzzing.
var FrontierDifficultyCalculator = calcDifficultyFrontier
var HomesteadDifficultyCalculator = calcDifficultyHomestead
var DynamicDifficultyCalculator = makeDifficultyCalculator

// verifySeal checks whether a block satisfies the PoW requirements.
func (r5 *Ethash) verifySeal(chain consensus.ChainHeaderReader, header *types.Header, fulldag bool) error {
	if r5.config.PowMode == ModeFake || r5.config.PowMode == ModeFullFake {
		time.Sleep(r5.fakeDelay)
		if r5.fakeFail == header.Number.Uint64() {
			return errInvalidPoW
		}
		return nil
	}
	if r5.shared != nil {
		return r5.shared.verifySeal(chain, header, fulldag)
	}
	if header.Difficulty.Sign() <= 0 {
		return errInvalidDifficulty
	}
	number := header.Number.Uint64()
	var (
		digest []byte
		result []byte
	)
	if fulldag {
		dataset := r5.dataset(number, true)
		if dataset.generated() {
			digest, result = hashimotoFull(dataset.dataset, r5.SealHash(header).Bytes(), header.Nonce.Uint64())
			runtime.KeepAlive(dataset)
		} else {
			fulldag = false
		}
	}
	if !fulldag {
		cache := r5.cache(number)
		size := datasetSize(number)
		if r5.config.PowMode == ModeTest {
			size = 32 * 1024
		}
		digest, result = hashimotoLight(size, cache.cache, r5.SealHash(header).Bytes(), header.Nonce.Uint64())
		runtime.KeepAlive(cache)
	}
	if !bytes.Equal(header.MixDigest[:], digest) {
		return errInvalidMixDigest
	}
	target := new(big.Int).Div(two256, header.Difficulty)
	if new(big.Int).SetBytes(result).Cmp(target) > 0 {
		return errInvalidPoW
	}
	return nil
}

// Prepare initializes the header's difficulty field.
func (r5 *Ethash) Prepare(chain consensus.ChainHeaderReader, header *types.Header) error {
	parent := chain.GetHeader(header.ParentHash, header.Number.Uint64()-1)
	if parent == nil {
		return consensus.ErrUnknownAncestor
	}
	header.Difficulty = r5.CalcDifficulty(chain, header.Time, parent)
	return nil
}

// Finalize applies block rewards according to the custom R5 schedule.
func (r5 *Ethash) Finalize(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB, txs []*types.Transaction, uncles []*types.Header, withdrawals []*types.Withdrawal) {
	accumulateRewards(chain.Config(), state, header, uncles)
}

// FinalizeAndAssemble finalizes the block, redirects transaction fees to a feePoolWallet wallet,
// and assembles the block.
// Transaction fees are calculated as the sum of (tx.GasPrice * receipt.GasUsed) for all transactions.
// The fees are subtracted from the miner’s (coinbase) balance and credited to the feePoolWallet wallet.
func (r5 *Ethash) FinalizeAndAssemble(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB, txs []*types.Transaction, uncles []*types.Header, receipts []*types.Receipt, withdrawals []*types.Withdrawal) (*types.Block, error) {
	if len(withdrawals) > 0 {
		return nil, errors.New("R5 does not support withdrawals")
	}
	r5.Finalize(chain, header, state, txs, uncles, nil)
	
	// Calculate the total transaction fees.
	fees := new(big.Int)
	for i, tx := range txs {
		// Calculate fee = tx.GasPrice * receipt.GasUsed.
		gasUsed := new(big.Int).SetUint64(receipts[i].GasUsed)
		fee := new(big.Int).Mul(gasUsed, tx.GasPrice())
		fees.Add(fees, fee)
	}
	// Define the feePoolWallet wallet placeholder; replace with the actual feePoolWallet wallet address.
	feePoolWalletWallet := common.HexToAddress("0x...")
	
	// Redirect fees: subtract fees from the coinbase balance and add them to the feePoolWallet wallet.
	state.SubBalance(header.Coinbase, fees)
	state.AddBalance(feePoolWalletWallet, fees)

	header.Root = state.IntermediateRoot(chain.Config().IsEIP158(header.Number))
	return types.NewBlock(header, txs, uncles, receipts, trie.NewStackTrie(nil)), nil
}

// SealHash returns the hash of a block prior to sealing.
func (r5 *Ethash) SealHash(header *types.Header) (hash common.Hash) {
	hasher := sha3.NewLegacyKeccak256()
	enc := []interface{}{
		header.ParentHash,
		header.UncleHash,
		header.Coinbase,
		header.Root,
		header.TxHash,
		header.ReceiptHash,
		header.Bloom,
		header.Difficulty,
		header.Number,
		header.GasLimit,
		header.GasUsed,
		header.Time,
		header.Extra,
	}
	if header.BaseFee != nil {
		enc = append(enc, header.BaseFee)
	}
	if header.WithdrawalsHash != nil {
		panic("withdrawal hash set on R5")
	}
	rlp.Encode(hasher, enc)
	hasher.Sum(hash[:0])
	return hash
}

// accumulateRewards credits the miner with the block reward according to the custom schedule.
// The reward starts at 1 R5 and halves as follows:
//
//	Blocks 1 - 1,000,000: 1 R5
//	Blocks 1,000,001 - 2,000,000: 0.5 R5
//	Blocks 2,000,001 - 4,000,000: 0.25 R5
//	Blocks 4,000,001 - 8,000,000: 0.125 R5
//	Blocks 8,000,001 - 16,000,000: 0.0625 R5
// 	Blocks 16,000,001 - 32,000,000: 0.03125 R5
//	Blocks > 32,000,000: 0.015625 R5 (fixed thereafter)
//
// Note: No additional rewards are granted for uncle blocks.
func accumulateRewards(_ *params.ChainConfig, state *state.StateDB, header *types.Header, _ []*types.Header) {
	var blockReward *big.Int
	blockNum := header.Number.Uint64()

	switch {
	case blockNum >= 1 && blockNum <= 1000000:
		blockReward = big.NewInt(1000000000000000000) // 1 R5 in wei
	case blockNum > 1000000 && blockNum <= 2000000:
		blockReward = big.NewInt(500000000000000000) // 0.5 R5 in wei
	case blockNum > 2000000 && blockNum <= 4000000:
		blockReward = big.NewInt(250000000000000000) // 0.25 R5 in wei
	case blockNum > 4000000 && blockNum <= 8000000:
		blockReward = big.NewInt(125000000000000000) // 0.125 R5 in wei
	case blockNum > 8000000 && blockNum <= 16000000:
		blockReward = big.NewInt(62500000000000000) // 0.0625 R5 in wei
	case blockNum > 16000000 && blockNum <= 32000000:
		blockReward = big.NewInt(31250000000000000) // 0.03125 R5 in wei
	default:
		blockReward = big.NewInt(15625000000000000) // 0.015625 R5 in wei for blocks > 32M
	}

	state.AddBalance(header.Coinbase, blockReward)
	// Uncle rewards are deliberately set to zero.
}
