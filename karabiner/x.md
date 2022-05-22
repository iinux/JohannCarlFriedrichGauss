# x11 apps

1. x11-apps: /usr/bin/atobm
1. x11-apps: /usr/bin/bitmap
1. x11-apps: /usr/bin/bmtoa
1. x11-apps: /usr/bin/ico
1. x11-apps: /usr/bin/oclock
1. x11-apps: /usr/bin/rendercheck
1. x11-apps: /usr/bin/transset
1. x11-apps: /usr/bin/x11perf
1. x11-apps: /usr/bin/x11perfcomp
1. x11-apps: /usr/bin/xbiff
1. x11-apps: /usr/bin/xcalc
1. x11-apps: /usr/bin/xclipboard
1. x11-apps: /usr/bin/xclock
1. x11-apps: /usr/bin/xconsole
1. x11-apps: /usr/bin/xcursorgen
1. x11-apps: /usr/bin/xcutsel
1. x11-apps: /usr/bin/xditview
1. x11-apps: /usr/bin/xedit
1. x11-apps: /usr/bin/xeyes
1. x11-apps: /usr/bin/xgc
1. x11-apps: /usr/bin/xload
1. x11-apps: /usr/bin/xlogo
1. x11-apps: /usr/bin/xmag
1. x11-apps: /usr/bin/xman
1. x11-apps: /usr/bin/xmore
1. x11-apps: /usr/bin/xwd
1. x11-apps: /usr/bin/xwud

# 概念

* Wayland 是想要取代 X window
* X window = X11 = X
* Xorg 是 X window 的一种实现
* DE 桌面环境 kde gnome
* DM 桌面管理器，比如 lightdm sddm
* WM 窗口管理器，比如 i3wm dwm 

# shortcut

```
sudo apt install xcape
xcape -e 'Control_L=Escape'
xcape -e 'Caps_Lock=Menu'

xmodmap
xmodmap -pke | less
dumpkeys -f | less
setxkbmap -option caps:ctrl_modifier
```