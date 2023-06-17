import os
import socket
import struct
import time

QCLASS_IN = 1
QTYPE_A = 1
QTYPE_AAAA = 28


def compat_chr(d):
    if bytes == str:
        return chr(d)
    return bytes([d])


def build_address(address):
    address = address.strip(b'.')
    labels = address.split(b'.')
    results = []
    for label in labels:
        l = len(label)
        if l > 63:
            return None
        results.append(compat_chr(l))
        results.append(label)
    results.append(b'\0')
    return b''.join(results)


def build_request(address, qtype):
    request_id = os.urandom(2)
    header = struct.pack('!BBHHHH', 1, 0, 1, 0, 0, 0)
    addr = build_address(address)
    qtype_qclass = struct.pack('!HH', qtype, QCLASS_IN)
    return request_id + header + addr + qtype_qclass


my_sock = socket.socket(socket.AF_INET6, socket.SOCK_DGRAM, socket.SOL_UDP)
hostname = 'www.google.com'.encode()
server = '192.168.1.1'
# server = '8.8.8.8'
server = '2001:67c:2b0::4'
req = build_request(hostname, QTYPE_A)
my_sock.setblocking(False)
my_sock.sendto(req, (server, 53))

if __name__ == '__main__':
    time.sleep(1)

