#!/usr/bin/env python3
"""
使用 Linux 内核加密接口 AF_ALG 进行 AES-CBC 加解密。

要求:
    - Linux 内核 >= 2.6.38 (启用 CONFIG_CRYPTO_USER_API_SKCIPHER)
    - Python >= 3.6
"""

import os
import socket


def _pkcs7_pad(data: bytes, block_size: int = 16) -> bytes:
    pad_len = block_size - (len(data) % block_size)
    return data + bytes([pad_len]) * pad_len


def _pkcs7_unpad(data: bytes) -> bytes:
    if not data:
        raise ValueError("空数据无法去填充")
    pad_len = data[-1]
    if pad_len < 1 or pad_len > 16 or data[-pad_len:] != bytes([pad_len]) * pad_len:
        raise ValueError("PKCS7 填充无效")
    return data[:-pad_len]


def af_alg_crypt(op: int, key: bytes, iv: bytes, data: bytes,
                 alg_name: str = "cbc(aes)") -> bytes:
    """
    使用 AF_ALG 执行对称加/解密。

    参数:
        op: socket.ALG_OP_ENCRYPT 或 socket.ALG_OP_DECRYPT
        key: 对称密钥 (AES: 16/24/32 字节)
        iv:  CBC 初始化向量 (16 字节)
        data: 待处理数据 (必须是 block_size 整数倍)
        alg_name: 内核算法名,默认 "cbc(aes)"
    返回:
        加/解密后的字节串
    """
    # 创建算法套接字 (类型: skcipher)
    alg_sock = socket.socket(socket.AF_ALG, socket.SOCK_SEQPACKET, 0)
    try:
        alg_sock.bind(("skcipher", alg_name))
        alg_sock.setsockopt(socket.SOL_ALG, socket.ALG_SET_KEY, key)

        # accept 后获得操作套接字,用 sendmsg_afalg 传 op/iv,然后写入数据
        op_sock, _ = alg_sock.accept()
        try:
            op_sock.sendmsg_afalg([data], op=op, iv=iv)
            return op_sock.recv(len(data))
        finally:
            op_sock.close()
    finally:
        alg_sock.close()


def encrypt(key: bytes, iv: bytes, plaintext: bytes) -> bytes:
    padded = _pkcs7_pad(plaintext, 16)
    return af_alg_crypt(socket.ALG_OP_ENCRYPT, key, iv, padded)


def decrypt(key: bytes, iv: bytes, ciphertext: bytes) -> bytes:
    if len(ciphertext) % 16 != 0:
        raise ValueError("密文长度必须是 16 的倍数")
    padded = af_alg_crypt(socket.ALG_OP_DECRYPT, key, iv, ciphertext)
    return _pkcs7_unpad(padded)


def main() -> None:
    key = os.urandom(32)              # AES-256
    iv = os.urandom(16)
    plaintext = "AF_ALG 内核加密接口测试: hello, 世界!".encode("utf-8")

    print(f"明文 ({len(plaintext)} 字节): {plaintext!r}")
    print(f"密钥 (hex): {key.hex()}")
    print(f"IV   (hex): {iv.hex()}")

    ciphertext = encrypt(key, iv, plaintext)
    print(f"密文 (hex): {ciphertext.hex()}")

    recovered = decrypt(key, iv, ciphertext)
    print(f"解密 ({len(recovered)} 字节): {recovered!r}")

    assert recovered == plaintext, "解密结果与原文不一致!"
    print("OK: 加解密往返一致")


if __name__ == "__main__":
    main()
