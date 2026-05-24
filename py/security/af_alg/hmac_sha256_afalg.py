#!/usr/bin/env python3
"""使用 Linux AF_ALG 接口计算 HMAC-SHA256"""

import socket


def hmac_sha256_afalg(key: bytes, data: bytes) -> bytes:
    """通过内核 AF_ALG 套接字计算 HMAC-SHA256

    Args:
        key:  HMAC 密钥
        data: 待认证的数据

    Returns:
        32 字节的 HMAC-SHA256 摘要
    """
    # 1. 创建 AF_ALG 套接字, 类型为 hash
    tfm = socket.socket(socket.AF_ALG, socket.SOCK_SEQPACKET, 0)
    try:
        # 2. 绑定到具体算法 hmac(sha256)
        tfm.bind(("hash", "hmac(sha256)"))

        # 3. 设置 HMAC 密钥 (SOL_ALG=279, ALG_SET_KEY=1)
        tfm.setsockopt(socket.SOL_ALG, socket.ALG_SET_KEY, key)

        # 4. accept 出一个操作套接字, 用于送数据/取结果
        op, _ = tfm.accept()
        try:
            op.sendall(data)
            return op.recv(32)  # SHA-256 输出 32 字节
        finally:
            op.close()
    finally:
        tfm.close()


def hmac_sha256_afalg_stream(key: bytes, chunks) -> bytes:
    """支持分块送入的版本, chunks 为可迭代字节块"""
    tfm = socket.socket(socket.AF_ALG, socket.SOCK_SEQPACKET, 0)
    try:
        tfm.bind(("hash", "hmac(sha256)"))
        tfm.setsockopt(socket.SOL_ALG, socket.ALG_SET_KEY, key)
        op, _ = tfm.accept()
        try:
            chunk_list = list(chunks)
            for i, chunk in enumerate(chunk_list):
                # MSG_MORE 表示后面还有更多数据
                flags = socket.MSG_MORE if i < len(chunk_list) - 1 else 0
                op.send(chunk, flags)
            return op.recv(32)
        finally:
            op.close()
    finally:
        tfm.close()


if __name__ == "__main__":
    import hmac
    import hashlib

    key = b"my-secret-key"
    data = b"hello, AF_ALG hmac(sha256)!"

    digest_kernel = hmac_sha256_afalg(key, data)
    digest_python = hmac.new(key, data, hashlib.sha256).digest()

    print("AF_ALG :", digest_kernel.hex())
    print("hashlib:", digest_python.hex())
    assert digest_kernel == digest_python, "结果不一致!"
    print("OK: 内核与 hashlib 结果一致")

    # 分块测试
    digest_stream = hmac_sha256_afalg_stream(
        key, [b"hello, ", b"AF_ALG ", b"hmac(sha256)!"]
    )
    assert digest_stream == digest_python
    print("OK: 分块结果也一致 ->", digest_stream.hex())
