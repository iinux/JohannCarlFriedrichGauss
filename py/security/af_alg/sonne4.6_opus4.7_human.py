#!/usr/bin/env python3
"""
使用 Linux socket.AF_ALG 接口进行加解密示例
AF_ALG 是 Linux 内核提供的用户态加密 API，通过 socket 调用内核加密子系统

支持的操作:
  - AES-CBC 加密 / 解密
  - AES-GCM 认证加密 / 解密
  - HMAC-SHA256 消息认证码
  - SHA-256 哈希

需要 Linux 内核 >= 3.6，且需要 root 或 CAP_NET_ADMIN 权限（部分发行版不需要）
"""

import socket
import os
import struct
import sys
from typing import Tuple
import traceback

# AF_ALG 相关常量
SOL_ALG = 279
ALG_SET_KEY = 1
ALG_SET_IV = 2
ALG_SET_OP = 3
ALG_SET_AEAD_AUTHSIZE = 4

ALG_OP_DECRYPT = 0
ALG_OP_ENCRYPT = 1

# cmsg 辅助函数
def _cmsg_afalg_iv(iv: bytes) -> bytes:
    """构造 IV 的控制消息"""
    # struct af_alg_iv { __u32 ivlen; __u8 iv[]; }
    return struct.pack("I", len(iv)) + iv

def alg_encrypt_decrypt(alg_type: str, alg_name: str, key: bytes,
                         data: bytes, iv: bytes, op: int) -> bytes:
    """
    通用 AF_ALG 加解密函数

    :param alg_type: 算法类型，如 "skcipher"、"aead"
    :param alg_name: 算法名称，如 "cbc(aes)"
    :param key:      密钥
    :param data:     待处理数据
    :param iv:       初始向量（无则传 b''）
    :param op:       ALG_OP_ENCRYPT 或 ALG_OP_DECRYPT
    :return:         处理后的数据
    """
    # 1. 创建 AF_ALG 绑定 socket
    sock = socket.socket(socket.AF_ALG, socket.SOCK_SEQPACKET, 0)
    sock.bind((alg_type, alg_name))

    # 2. 设置密钥
    sock.setsockopt(SOL_ALG, ALG_SET_KEY, key)

    # 3. accept() 得到操作 socket
    op_sock, _ = sock.accept()

    # 5. 发送数据
    op_sock.sendmsg_afalg([data], op=op, iv=iv)

    # 6. 接收结果
    result = op_sock.recv(len(data))

    op_sock.close()
    sock.close()
    return result


# ─────────────────────────────────────────────
# AES-CBC
# ─────────────────────────────────────────────

def pkcs7_pad(data: bytes, block_size: int = 16) -> bytes:
    pad_len = block_size - (len(data) % block_size)
    return data + bytes([pad_len] * pad_len)

def pkcs7_unpad(data: bytes) -> bytes:
    pad_len = data[-1]
    return data[:-pad_len]

def aes_cbc_encrypt(key: bytes, iv: bytes, plaintext: bytes) -> bytes:
    padded = pkcs7_pad(plaintext)
    return alg_encrypt_decrypt("skcipher", "cbc(aes)", key, padded, iv, ALG_OP_ENCRYPT)

def aes_cbc_decrypt(key: bytes, iv: bytes, ciphertext: bytes) -> bytes:
    padded = alg_encrypt_decrypt("skcipher", "cbc(aes)", key, ciphertext, iv, ALG_OP_DECRYPT)
    return pkcs7_unpad(padded)


# ─────────────────────────────────────────────
# HMAC-SHA256
# ─────────────────────────────────────────────

def hmac_sha256(key: bytes, data: bytes) -> bytes:
    sock = socket.socket(socket.AF_ALG, socket.SOCK_SEQPACKET, 0)
    sock.bind(("hash", "hmac(sha256)"))
    sock.setsockopt(SOL_ALG, ALG_SET_KEY, key)
    op_sock, _ = sock.accept()
    op_sock.sendall(data)
    result = op_sock.recv(32)
    op_sock.close()
    sock.close()
    return result


# ─────────────────────────────────────────────
# SHA-256 哈希
# ─────────────────────────────────────────────

def sha256(data: bytes) -> bytes:
    sock = socket.socket(socket.AF_ALG, socket.SOCK_SEQPACKET, 0)
    sock.bind(("hash", "sha256"))
    op_sock, _ = sock.accept()
    op_sock.sendall(data)
    result = op_sock.recv(32)
    op_sock.close()
    sock.close()
    return result


# ─────────────────────────────────────────────
# 演示
# ─────────────────────────────────────────────

def demo_aes_cbc():
    print("=" * 50)
    print("AES-CBC 加解密")
    print("=" * 50)

    key = os.urandom(32)   # AES-256
    iv  = os.urandom(16)
    plaintext = b"Hello, AF_ALG! This is a test message for AES-CBC encryption."

    print(f"明文:       {plaintext.decode()}")
    print(f"密钥(hex):  {key.hex()}")
    print(f"IV(hex):    {iv.hex()}")

    ciphertext = aes_cbc_encrypt(key, iv, plaintext)
    print(f"密文(hex):  {ciphertext.hex()}")

    decrypted = aes_cbc_decrypt(key, iv, ciphertext)
    print(f"解密结果:   {decrypted.decode()}")
    print(f"验证:       {'✓ 成功' if decrypted == plaintext else '✗ 失败'}")
    print()

def demo_hmac():
    print("=" * 50)
    print("HMAC-SHA256")
    print("=" * 50)

    key  = os.urandom(32)
    data = b"Message to authenticate"

    mac = hmac_sha256(key, data)
    print(f"数据:       {data.decode()}")
    print(f"密钥(hex):  {key.hex()}")
    print(f"MAC(hex):   {mac.hex()}")

    # 篡改数据验证
    tampered = b"Message to authenticate!"
    mac2 = hmac_sha256(key, tampered)
    print(f"篡改后MAC:  {mac2.hex()}")
    print(f"MAC一致:    {'✓' if mac == mac2 else '✗ 数据已被篡改'}")
    print()

def demo_sha256():
    print("=" * 50)
    print("SHA-256 哈希")
    print("=" * 50)

    data = b"Hello, world!"
    digest = sha256(data)
    print(f"数据:       {data.decode()}")
    print(f"SHA256:     {digest.hex()}")

    # 对比标准库结果
    import hashlib
    expected = hashlib.sha256(data).digest()
    print(f"标准库结果: {expected.hex()}")
    print(f"验证:       {'✓ 一致' if digest == expected else '✗ 不一致'}")
    print()


if __name__ == "__main__":
    if not sys.platform.startswith("linux"):
        print("错误: AF_ALG 仅支持 Linux 系统")
        sys.exit(1)

    print("Linux AF_ALG 内核加密接口演示")
    print("通过 socket(AF_ALG) 调用内核加密子系统\n")

    try:
        demo_aes_cbc()
        demo_hmac()
        demo_sha256()
    except OSError as e:
        traceback.print_exc()
        print(f"错误: {e}")
        print("提示: 部分操作可能需要 root 权限，或内核未编译对应加密模块")
        print("      可用 'grep -i aes /proc/crypto' 查看支持的算法")
        sys.exit(1)
