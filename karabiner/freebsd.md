# serial

`tip ucom1 -115200`

`cu -s 115200 -l /dev/cuaU0`

`man uchcom`

---

```
/var/spool/lock/LCK..cuaU0: No such file or directory
Can't open lock file.
all ports busy
```

need root

---

https://docs.freebsd.org/zh-cn/books/handbook/serialcomms/

---

/boot/loader.conf
```
uchcom_load="YES"
```

# links

1. https://mirrors.ustc.edu.cn/help/freebsd-pkg.html
1. https://docs.freebsd.org/doc/13.0-RELEASE/usr/local/share/doc/freebsd/
1. http://shouce.jb51.net/freebsd-handbook/
1. https://handbook.freebsdcn.org/di-7-zhang-duo-mei-ti/7.2.-she-zhi-sheng-ka
1. https://docs.freebsd.org/zh-cn/books/handbook/x11/
1. https://docs.freebsd.org/doc/6.2-RELEASE/usr/share/doc/handbook/network-dhcp.html
1. https://srobb.net/fbsdquickwireless.html
1. https://blog.csdn.net/weixin_34137975/article/details/116832026
1. https://www.cnblogs.com/freedom-try/p/15438079.html
1. https://docs.huihoo.com/freebsd/zh_CN.GB2312/linuxemu-lbc-install.html
1. https://wiki.freebsd.org/Graphics
1. https://community.kde.org/FreeBSD/Setup#Graphics_first
1. https://bugs.freebsd.org/bugzilla/show_bug.cgi?id=231884
1. https://www.freebsd.org/cgi/man.cgi?query=radeon&sektion=4&format=html
1. https://unixsheikh.com/tutorials/how-to-setup-freebsd-with-a-riced-desktop-part-3-i3.html

# remote wireshark

```
ssh hw 'sudo dumpcap -w - -f "not port 22"' > hwwiresharkcli
ssh hw 'sudo tcpdump -w - "not port 22"' > hwwiresharkcli
wireshark -k -i hwwiresharkcli
sudo tcpdump -w - "not port 22" | nc -l 2222
nc hw.iinux.cn 2222 | wireshark -k -i -
```

# other based version

* nomadBSD
* GhostB

# fix chinese error display

https://codeantenna.com/a/hYPw31x7M7

`LANG=zh_CN.UTF-8 wine WePE_32_V1.2.exe`

# highlight

zfs

# pkg command

```
pkg info --pkg-message wireshark
pkg info -D wireshark
pkg ins shadowsocks-libev
pkg ins py38-pip
pkg ins usbutils
pkg ins coreutils
pkg ins gmake
pkg ins wqy-fonts
pkg ins neofetch
pkg ins wifimgr
pkg ins gnome3
pkg ins kde5
pkg install i3-gaps i3lock i3status conky dmenu xterm gnome-screenshot nitrogen xorg
```

# drm

```
pkg ins drm-kmod
sysrc -f /etc/rc.conf kld_list+=radeonkms
sysrc -f /etc/rc.conf kld_list+=amdgpu
```

# for support ntfs

`pkg install fusefs-ntfs`

add to /boot/loader.conf
```
fusefs_load="YES"
```

# sound

```
cat /dev/sndstat
sysctl hw.snd.default_unit=
https://wiki.freebsd.org/Sound
pkg install audio/dsbmixer
pkg install dsbmixer
pkg install audio/gtk-mixer
pkg install gtk-mixer
pkg install audio/mixertui
pkg install mixertui
```

# 蓝牙耳机(未成功,仅参考)
https://docs.freebsd.org/en/books/handbook/advanced-networking/#network-bluetooth
https://docs.freebsd.org/doc/6.0-RELEASE/usr/share/doc/handbook/network-bluetooth.html
https://forums.freebsd.org/threads/using-bluetooth-audio-devices-speaker-headphones-earbuds-with-freebsd.78992/
https://forums.freebsd.org/threads/bluetooth-audio-how-to-connect-and-use-bluetooth-headphones-on-freebsd.82671/

# commnad

* `pkg stats`
* `brandelf -t Linux <elf_file>` 给elf文件指定 Linux type ，以便可以使用Linux兼容层运行
* `pkg ins linux_base-c7` 安装Linux兼容层
* `pciconf -lvbce command output`
* `devinfo -vr command output`
* `sysctl hw.model`
* `swapinfo`
* `gmd5sum` gnu 版的计算md5

# stat

## sockstat

## vmstat
procs:
* r-->在运行的进程数
* b-->在等待io的进程数(等待i/o,paging等等)
* w-->可以进入运行队列但被替换的进程

memoy（以k为单位，包括虚拟内存和真实内存，正在运行或最近20秒在运行的进程所用的虚拟内存将被视为active）
* avm-->活动的虚拟内存
* free-->空闲的内存

pages（统计错误页和活动页，每5秒平均一下，以秒为单位给出数值）
* flt-->错误页总数
* re-->回收的页面
* pi-->进入页面数
* po-->出页面数
* fr-->空余的页面数
* sr-->每秒通过时钟算法扫描的页面

disk 显示每秒的磁盘操作（磁盘名字的前两个字母加数字，默认只显示两个磁盘，如果有多的，可以加-n来增加数字或在命令行下把磁盘名都填上。）

fault 显示每秒的中断数
* in-->设备中断
* sy-->系统中断
* cy-->cpu交换

from 
https://blog.51cto.com/sasyun/1530841

## fstat
## gstat
## pstat
## iostat
## netstat
## nfsstat
## systat

# 查看cpu多少核
`sysctl kern.smp.cpus`

# wlan
* `sysctl -n net.wlan.devices`
* `wpa_passphrase mywpa 12345678 > wpa_supplicant.conf`
* `bsdconfig netdev`
* `kldstat`
* `kldload`

sudo vim /boot/loader.conf
```
if_ath_load="YES"
wlan_scan_ap_load="YES"
wlan_scan_sta_load="YES"
wlan_wep_load="YES"
wlan_ccmp_load="YES"
wlan_tkip_load="YES"
autoboot_delay="3"
```

# make driver from windows driver (fail)

```
ndisgen netr28ux.inf netr28ux.sys
ndiscvt
sudo usbconfig
/usr/src/sys/dev/usb/usbdevs
```

# input method
.xinitrc
```
export GTK_IM_MODULE=fcitx
export QT_IM_MODULE=fcitx
export XMODIFIERS="@im=fcitx"
```

# startx
.xinitrc
```
/usr/local/bin/gnome-session
exec startplasma-x11
```

# net
```
pkg install bind-tools # nslookup
route -n add -inet6 default 2a00:b700::1
ifconfig vtnet0 inet6 2a00:b700::b:138
```

/etc/netstart

`sudo sysrc dbus_enable`

/etc/rc.conf
```
ifconfig_fxp0="DHCP"
hald_enable="YES"
dbus_enable="YES"
ipv6_enable="YES"
ipv6_ifconfig_fxp0="2001:470:1e04:5ea::10"
ipv6_defaultrouter="2001:470:1e04:5ea::1"
wlans_ath0="wlan0"
ifconfig_wlan0="WPA SYNCDHCP"
gnome_enable="YES"
ifconfig_em0="inet 192.168.1.100 netmask 255.255.255.0"
defaultrouter="192.168.1.1"
```
`service netif restart`

`service routing restart`

要在FreeBSD上为网络接口添加多个IP地址，请在/etc/rc.conf文件中添加以下行。
ifconfig_em0_alias0="192.168.1.5 netmask 255.255.255.255"

`/etc/rc.d/network_ipv6 restart`

# install gnome need
```
mkdir -p /var/lib/dbus
dbus-uuidgen  > /var/lib/dbus/machine-id
```

# bbr
https://www.cnblogs.com/rn7s2/p/15212124.html

# docker
https://grass.show/post/use-docker-on-freebsd/

`sudo pkg install virtualbox-ose-nox11`

/boot/loader.conf  
```
vboxdrv_load="YES"
```
`sudo pw groupmod vboxusers -m grass` # grass是我的用户名
/etc/rc.conf   
```
vboxnet_enable="YES"
```

`sudo pkg install docker-machine`

`docker-machine create -d virtualbox default` # 这步及后面没有成功

`sudo pkg install docker`

`eval (docker-machine env default)`

`docker run hello-world`
https://www.jianshu.com/p/9b861a805ee3
```
export DOCKER_HOST=tcp://192.168.59.103:2376
export DOCKER_CERT_PATH=/Users/wangshuang/.boot2docker/certs/boot2docker-vm
export DOCKER_TLS_VERIFY=1
```

# cmake include dir
```
target_include_directories(test PRIVATE ${YOUR_DIRECTORY})
```
https://stackoverflow.com/questions/13703647/how-to-properly-add-include-directories-with-cmake

# sudoers

```
/usr/local/etc/sudoers
visudo
uname -a
```

# gpu

```
gpu monitor
sudo apt install radeontop
apt install nvtop
sudo apt install nvidia-smi
sudo apt install intel-gpu-tools
pip install gpustat // for nvidia
sudo apt install mesa-utils
glxgears
sudo apt install glmark2
```
https://www.sohu.com/a/439433413_495675


# pkg alias
If you can't remember the options create an alias for it. 

Edit `/usr/local/etc/pkg.conf `and add to the ALIAS section:
Code:
```
  message: info -D
  ```
Now you can use `pkg message <pkgname>`

# protsnap
* `portsnap fetch`
* `portsnap extract`
* `portsnap update`

# swap
1.Create the swap file:

`dd if=/dev/zero of=/usr/swap0 bs=1m count=64`

2.Set the proper permissions on the new file:

`chmod 0600 /usr/swap0`

3.Inform the system about the swap file by adding a line to /etc/fstab:
```
md99	none	swap	sw,file=/usr/swap0,late	0	0
```
The md(4) device md99 is used, leaving lower device numbers available for interactive use.

4.Swap space will be added on system startup. To add swap space immediately, use swapon(8):

`swapon -aL`

refer https://www.freebsd.org/doc/handbook/adding-swap-space.html

# pf
sysrc pf_enable=yes

refer https://docs.freebsd.org/en/books/handbook/firewalls/#firewalls-pf

tcpdump -n -e -ttt -r /var/log/pflog

# remote x server

```bash
lsof | grep -i xorg
socat -d -d TCP-LISTEN:6000,fork UNIX-CONNECT:/tmp/.X11-unix/X0

DISPLAY=127.0.0.1:0 xclock

xhost +
DISPLAY=192.168.1.5:0 xclock

xauth list
```


