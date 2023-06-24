import os
import fcntl
import struct
import select
import socket


TUNSETIFF = 0x400454ca
IFF_TUN = 0x0001
IFF_NO_PI = 0x1000


def create_tun_interface():
    tun_fd = os.open('/dev/net/tun', os.O_RDWR)
    ifr = struct.pack('16sH', b'tun0', IFF_TUN | IFF_NO_PI)
    print(ifr)
    fcntl.ioctl(tun_fd, TUNSETIFF, ifr)
    if_name = struct.unpack('16sH', ifr)[0].strip(b'\x00').decode()
    return tun_fd, if_name


def configure_ip_address(if_name, ip_address, netmask):
    os.system(f'ip addr add {ip_address}/{netmask} dev {if_name}')
    os.system(f'ip link set {if_name} up')


def run_tun():
    tun_fd, if_name = create_tun_interface()
    print(if_name)
    configure_ip_address(if_name, '10.0.0.1', 24)

    while True:
        r, _, _ = select.select([tun_fd], [], [])
        if tun_fd in r:
            packet = os.read(tun_fd, 65536)
            # 处理接收到的数据包
            print('Received packet:', packet)


if __name__ == '__main__':
    run_tun()

