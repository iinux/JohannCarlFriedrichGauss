# serial

`tip ucom1 -115200`
`cu -s 115200 -l /dev/cuaU0`

```
/var/spool/lock/LCK..cuaU0: No such file or directory
Can't open lock file.
all ports busy
```

need root

https://docs.freebsd.org/zh-cn/books/handbook/serialcomms/

/boot/loader.conf
uchcom_load="YES"

`man uchcom`
