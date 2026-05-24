#!/usr/bin/env python3
"""使用 Linux AF_ALG 内核加密接口计算 SHA256"""
import socket
import sys


def af_alg_sha256(data: bytes) -> bytes:
    # 创建 AF_ALG socket
    sock = socket.socket(socket.AF_ALG, socket.SOCK_SEQPACKET, 0)
    try:
        # 绑定到 sha256 hash 算法
        sock.bind(("hash", "sha256"))
        # accept() 返回一个用于本次操作的子 socket
        op, _ = sock.accept()
        try:
            op.sendall(data)
            # SHA256 摘要长度为 32 字节
            return op.recv(32)
        finally:
            op.close()
    finally:
        sock.close()


def af_alg_sha256_file(path: str, chunk_size: int = 65536) -> bytes:
    sock = socket.socket(socket.AF_ALG, socket.SOCK_SEQPACKET, 0)
    try:
        sock.bind(("hash", "sha256"))
        op, _ = sock.accept()
        try:
            with open(path, "rb") as f:
                while True:
                    chunk = f.read(chunk_size)
                    if not chunk:
                        break
                    # MSG_MORE 表示还有后续数据，内核继续累积
                    op.send(chunk, socket.MSG_MORE)
            return op.recv(32)
        finally:
            op.close()
    finally:
        sock.close()


if __name__ == "__main__":
    if len(sys.argv) > 1:
        digest = af_alg_sha256_file(sys.argv[1])
        print(f"{digest.hex()}  {sys.argv[1]}")
    else:
        data = b"hello, AF_ALG"
        digest = af_alg_sha256(data)
        print(f"data:   {data!r}")
        print(f"sha256: {digest.hex()}")

        # 与 hashlib 对比验证
        import hashlib
        expected = hashlib.sha256(data).hexdigest()
        print(f"expect: {expected}")
        assert digest.hex() == expected, "AF_ALG 与 hashlib 结果不一致"
        print("OK")
