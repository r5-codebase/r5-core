#############################################################################
# Generated by PAGE version 8.0
#  in conjunction with Tcl version 8.6
#  Feb 26, 2025 08:08:26 PM +1300  platform: Windows NT
set vTcl(timestamp) ""
if {![info exists vTcl(borrow)]} {
    ::vTcl::MessageBox -title Error -message  "You must open project files from within PAGE."
    exit}


set vTcl(actual_gui_font_dft_desc)  TkDefaultFont
set vTcl(actual_gui_font_dft_name)  TkDefaultFont
set vTcl(actual_gui_font_text_desc)  TkTextFont
set vTcl(actual_gui_font_text_name)  TkTextFont
set vTcl(actual_gui_font_fixed_desc)  TkFixedFont
set vTcl(actual_gui_font_fixed_name)  TkFixedFont
set vTcl(actual_gui_font_menu_desc)  TkMenuFont
set vTcl(actual_gui_font_menu_name)  TkMenuFont
set vTcl(actual_gui_font_tooltip_desc)  TkDefaultFont
set vTcl(actual_gui_font_tooltip_name)  TkDefaultFont
set vTcl(actual_gui_font_treeview_desc)  TkDefaultFont
set vTcl(actual_gui_font_treeview_name)  TkDefaultFont
########################################### 
set vTcl(actual_gui_bg) #d9d9d9
set vTcl(actual_gui_fg) #000000
set vTcl(actual_gui_analog) #ececec
set vTcl(actual_gui_menu_analog) #ececec
set vTcl(actual_gui_menu_bg) #d9d9d9
set vTcl(actual_gui_menu_fg) #000000
set vTcl(complement_color) gray40
set vTcl(analog_color_p) #c3c3c3
set vTcl(analog_color_m) beige
set vTcl(tabfg1) black
set vTcl(tabfg2) white
set vTcl(actual_gui_menu_active_bg)  #ececec
set vTcl(actual_gui_menu_active_fg)  #000000
########################################### 
set vTcl(pr,autoalias) 1
set vTcl(pr,relative_placement) 1
set vTcl(mode) Relative
set vTcl(project_theme) vista



proc vTclWindow.top1 {base} {
    global vTcl
    if {$base == ""} {
        set base .top1
    }
    if {[winfo exists $base]} {
        wm deiconify $base; return
    }
    set top $base
    set target $base
    ###################
    # CREATING WIDGETS
    ###################
    vTcl::widgets::core::toplevel::createCmd $top -class Toplevel \
        -menu $top.m49 -background SystemButtonFace -highlightbackground SystemButtonFace \
        -highlightcolor SystemWindowText 
    wm focusmodel $top passive
    wm geometry $top 758x571+448+159
    update
    # set in toplevel.wgt.
    global vTcl
    global img_list
    set vTcl(save,dflt,origin) 0
    wm maxsize $top 3844 1041
    wm minsize $top 120 1
    wm overrideredirect $top 0
    wm resizable $top 1 1
    wm deiconify $top
    set toptitle "R5 Desktop Wallet"
    wm title $top $toptitle
    namespace eval ::widgets::${top}::ClassOption {}
    set ::widgets::${top}::ClassOption(-toptitle) $toptitle
    vTcl:DefineAlias "$top" "MainWindow" vTcl:Toplevel:WidgetProc "" 1
    set vTcl(real_top) {}
    menu "$top.m49" \
        -activebackground #d9d9d9 -activeforeground black \
        -font "-family {Segoe UI} -size 9" -tearoff 0 
    
set site_3_0 $top.m49
    $top.m49 add cascade \
        -menu "$top.m49.men50" -compound left -font "-family {Segoe UI} -size 9" -label "File" 
    menu "$site_3_0.men50" \
        -activebackground #d9d9d9 -activeforeground black \
        -font "-family {Segoe UI} -size 9" -tearoff 0 
    $site_3_0.men50 add command \
        -command "#" -compound left -font "-family {Segoe UI} -size 9" \
        -label "Create Wallet..." 
    $site_3_0.men50 add command \
        -command "#" -compound left -font "-family {Segoe UI} -size 9" \
        -label "Load Wallet..." 
    $site_3_0.men50 add command \
        -command "#" -compound left -font "-family {Segoe UI} -size 9" \
        -label "Close Wallet..." -state disabled 
    $site_3_0.men50 add command \
        -command "#" -compound left -font "-family {Segoe UI} -size 9" \
        -label "Close All Wallets..." -state disabled 
    $site_3_0.men50 add separator \
        
    $site_3_0.men50 add command \
        -command "#" -compound left -font "-family {Segoe UI} -size 9" \
        -label "Settings..." 
    $site_3_0.men50 add separator \
        
    $site_3_0.men50 add command \
        -command "#" -compound left -font "-family {Segoe UI} -size 9" \
        -label "Exit" 
    
set site_3_0 $top.m49
    $top.m49 add cascade \
        -menu "$top.m49.men71" -compound left -font "-family {Segoe UI} -size 9" -label "Node" 
    menu "$site_3_0.men71" \
        -activebackground #d9d9d9 -activeforeground black \
        -font "-family {Segoe UI} -size 9" -tearoff 0 
    $site_3_0.men71 add command \
        -command "#" -compound left -font "-family {Segoe UI} -size 9" \
        -label "Node Console..." 
    $site_3_0.men71 add command \
        -command "#" -compound left -font "-family {Segoe UI} -size 9" \
        -label "JS Console..." 
    
set site_3_0 $top.m49
    $top.m49 add cascade \
        -menu "$top.m49.men52" -compound left -font "-family {Segoe UI} -size 9" -label "Help" 
    menu "$site_3_0.men52" \
        -activebackground #d9d9d9 -activeforeground black \
        -font "-family {Segoe UI} -size 9" -tearoff 0 
    $site_3_0.men52 add command \
        -command "#" -compound left -font "-family {Segoe UI} -size 9" \
        -label "About" 
    ###################
    # SETTING GEOMETRY
    ###################

    vTcl:FireEvent $base <<Ready>>
}

proc vTclWindow.top5 {base} {
    global vTcl
    if {$base == ""} {
        set base .top5
    }
    if {[winfo exists $base]} {
        wm deiconify $base; return
    }
    set top $base
    set target $base
    ###################
    # CREATING WIDGETS
    ###################
    vTcl::widgets::core::toplevel::createCmd $top -class Toplevel \
        -background SystemButtonFace -highlightcolor SystemWindowText 
    wm focusmodel $top passive
    wm geometry $top 661x387+471+300
    update
    # set in toplevel.wgt.
    global vTcl
    global img_list
    set vTcl(save,dflt,origin) 0
    wm maxsize $top 3844 1061
    wm minsize $top 120 1
    wm overrideredirect $top 0
    wm resizable $top 1 1
    wm iconify $top
    set toptitle "R5 Wallet Settings"
    wm title $top $toptitle
    namespace eval ::widgets::${top}::ClassOption {}
    set ::widgets::${top}::ClassOption(-toptitle) $toptitle
    vTcl:DefineAlias "$top" "Settings" vTcl:Toplevel:WidgetProc "" 1
    set vTcl(real_top) {}
    ttk::button "$top.tBu62" \
        -text "Load Settings File..." -compound left 
    vTcl:DefineAlias "$top.tBu62" "TButton3" vTcl:WidgetProc "Settings" 1
    ttk::button "$top.tBu63" \
        -text "Save Settings File..." -compound left 
    vTcl:DefineAlias "$top.tBu63" "TButton4" vTcl:WidgetProc "Settings" 1
    ttk::button "$top.tBu60" \
        -text "Save" -compound left 
    vTcl:DefineAlias "$top.tBu60" "TButton1" vTcl:WidgetProc "Settings" 1
    ttk::button "$top.tBu61" \
        -text "Cancel" -compound left 
    vTcl:DefineAlias "$top.tBu61" "TButton2" vTcl:WidgetProc "Settings" 1
    vTcl::widgets::ttk::notebook::createCmd "$top.tNo66" \
        -width 644 -height 306 
    vTcl:DefineAlias "$top.tNo66" "TNotebook2" vTcl:WidgetProc "Settings" 1
    frame "$top.tNo66.t0" \
        -background SystemButtonFace -highlightcolor SystemWindowText 
    vTcl:DefineAlias "$top.tNo66.t0" "TNotebook2_t1" vTcl:WidgetProc "Settings" 1
    $top.tNo66 add $top.tNo66.t0 \
        -padding 0 -sticky nsew -state normal -text "Application" \
        -compound left -underline -1 
    set site_4_0  $top.tNo66.t0
    frame "$top.tNo66.t1" \
        -background SystemButtonFace -highlightcolor SystemWindowText 
    vTcl:DefineAlias "$top.tNo66.t1" "TNotebook2_t2" vTcl:WidgetProc "Settings" 1
    $top.tNo66 add $top.tNo66.t1 \
        -padding 0 -sticky nsew -state normal -text "Wallet" -compound left \
        -underline -1 
    set site_4_1  $top.tNo66.t1
    frame "$top.tNo66.t2" \
        -background SystemButtonFace -highlightcolor SystemWindowText 
    vTcl:DefineAlias "$top.tNo66.t2" "TNotebook2_t3" vTcl:WidgetProc "Settings" 1
    $top.tNo66 add $top.tNo66.t2 \
        -padding 0 -sticky nsew -state normal -text "Node" -compound left \
        -underline -1 
    set site_4_2  $top.tNo66.t2
    ###################
    # SETTING GEOMETRY
    ###################
    place $top.tBu62 \
        -in $top -x 0 -relx 0.015 -y 0 -rely 0.827 -width 136 -relwidth 0 \
        -height 25 -relheight 0 -anchor nw -bordermode ignore 
    place $top.tBu63 \
        -in $top -x 0 -relx 0.015 -y 0 -rely 0.904 -width 136 -relwidth 0 \
        -height 25 -relheight 0 -anchor nw -bordermode ignore 
    place $top.tBu60 \
        -in $top -x 0 -relx 0.862 -y 0 -rely 0.904 -width 76 -relwidth 0 \
        -height 25 -relheight 0 -anchor nw -bordermode ignore 
    place $top.tBu61 \
        -in $top -x 0 -relx 0.741 -y 0 -rely 0.904 -width 76 -relwidth 0 \
        -height 25 -relheight 0 -anchor nw -bordermode ignore 
    place $top.tNo66 \
        -in $top -x 0 -relx 0.015 -y 0 -rely 0.026 -width 0 -relwidth 0.974 \
        -height 0 -relheight 0.791 -anchor nw -bordermode ignore 

    vTcl:FireEvent $base <<Ready>>
}

proc vTclWindow.top6 {base} {
    global vTcl
    if {$base == ""} {
        set base .top6
    }
    if {[winfo exists $base]} {
        wm deiconify $base; return
    }
    set top $base
    set target $base
    ###################
    # CREATING WIDGETS
    ###################
    vTcl::widgets::core::toplevel::createCmd $top -class Toplevel \
        -background SystemButtonFace -highlightcolor SystemWindowText 
    wm focusmodel $top passive
    wm geometry $top 600x450+470+263
    update
    # set in toplevel.wgt.
    global vTcl
    global img_list
    set vTcl(save,dflt,origin) 0
    wm maxsize $top 3844 1061
    wm minsize $top 120 1
    wm overrideredirect $top 0
    wm resizable $top 1 1
    wm iconify $top
    set toptitle "JS Console"
    wm title $top $toptitle
    namespace eval ::widgets::${top}::ClassOption {}
    set ::widgets::${top}::ClassOption(-toptitle) $toptitle
    vTcl:DefineAlias "$top" "JSConsole" vTcl:Toplevel:WidgetProc "" 1
    set vTcl(real_top) {}
    ###################
    # SETTING GEOMETRY
    ###################

    vTcl:FireEvent $base <<Ready>>
}

proc vTclWindow.top7 {base} {
    global vTcl
    if {$base == ""} {
        set base .top7
    }
    if {[winfo exists $base]} {
        wm deiconify $base; return
    }
    set top $base
    set target $base
    ###################
    # CREATING WIDGETS
    ###################
    vTcl::widgets::core::toplevel::createCmd $top -class Toplevel \
        -background SystemButtonFace -highlightcolor SystemWindowText 
    wm focusmodel $top passive
    wm geometry $top 600x450+530+241
    update
    # set in toplevel.wgt.
    global vTcl
    global img_list
    set vTcl(save,dflt,origin) 0
    wm maxsize $top 3844 1061
    wm minsize $top 120 1
    wm overrideredirect $top 0
    wm resizable $top 1 1
    wm iconify $top
    set toptitle "Node Console"
    wm title $top $toptitle
    namespace eval ::widgets::${top}::ClassOption {}
    set ::widgets::${top}::ClassOption(-toptitle) $toptitle
    vTcl:DefineAlias "$top" "NodeConsole" vTcl:Toplevel:WidgetProc "" 1
    set vTcl(real_top) {}
    ###################
    # SETTING GEOMETRY
    ###################

    vTcl:FireEvent $base <<Ready>>
}

proc vTclWindow.top8 {base} {
    global vTcl
    if {$base == ""} {
        set base .top8
    }
    if {[winfo exists $base]} {
        wm deiconify $base; return
    }
    set top $base
    set target $base
    ###################
    # CREATING WIDGETS
    ###################
    vTcl::widgets::core::toplevel::createCmd $top -class Toplevel \
        -background SystemButtonFace -highlightcolor SystemWindowText 
    wm focusmodel $top passive
    wm geometry $top 275x203+508+238
    update
    # set in toplevel.wgt.
    global vTcl
    global img_list
    set vTcl(save,dflt,origin) 0
    wm maxsize $top 3844 1061
    wm minsize $top 120 1
    wm overrideredirect $top 0
    wm resizable $top 1 1
    wm iconify $top
    set toptitle "About R5 Wallet"
    wm title $top $toptitle
    namespace eval ::widgets::${top}::ClassOption {}
    set ::widgets::${top}::ClassOption(-toptitle) $toptitle
    vTcl:DefineAlias "$top" "About" vTcl:Toplevel:WidgetProc "" 1
    set vTcl(real_top) {}
    ###################
    # SETTING GEOMETRY
    ###################

    vTcl:FireEvent $base <<Ready>>
}

proc 36 {args} {return 1}


Window show .
set btop1 ""
if {$vTcl(borrow)} {
    set btop1 .bor[expr int([expr rand() * 100])]
    while {[lsearch $btop1 $vTcl(tops)] != -1} {
        set btop1 .bor[expr int([expr rand() * 100])]
    }
}
set vTcl(btop) $btop1
Window show .top1 $btop1
if {$vTcl(borrow)} {
    $btop1 configure -background plum
}
set btop2 ""
if {$vTcl(borrow)} {
    set btop2 .bor[expr int([expr rand() * 100])]
    while {[lsearch $btop2 $vTcl(tops)] != -1} {
        set btop2 .bor[expr int([expr rand() * 100])]
    }
}
set vTcl(btop) $btop2
Window show .top5 $btop2
if {$vTcl(borrow)} {
    $btop2 configure -background plum
}
set btop3 ""
if {$vTcl(borrow)} {
    set btop3 .bor[expr int([expr rand() * 100])]
    while {[lsearch $btop3 $vTcl(tops)] != -1} {
        set btop3 .bor[expr int([expr rand() * 100])]
    }
}
set vTcl(btop) $btop3
Window show .top6 $btop3
if {$vTcl(borrow)} {
    $btop3 configure -background plum
}
set btop4 ""
if {$vTcl(borrow)} {
    set btop4 .bor[expr int([expr rand() * 100])]
    while {[lsearch $btop4 $vTcl(tops)] != -1} {
        set btop4 .bor[expr int([expr rand() * 100])]
    }
}
set vTcl(btop) $btop4
Window show .top7 $btop4
if {$vTcl(borrow)} {
    $btop4 configure -background plum
}
set btop5 ""
if {$vTcl(borrow)} {
    set btop5 .bor[expr int([expr rand() * 100])]
    while {[lsearch $btop5 $vTcl(tops)] != -1} {
        set btop5 .bor[expr int([expr rand() * 100])]
    }
}
set vTcl(btop) $btop5
Window show .top8 $btop5
if {$vTcl(borrow)} {
    $btop5 configure -background plum
}

