#!/usr/bin/env python3
import socket
import os
import struct

def pkcs7_pad(data: bytes, block_size: int = 16) -> bytes:
    pad_len = block_size - (len(data) % block_size)
    return data + bytes([pad_len] * pad_len)

def pkcs7_unpad(data: bytes) -> bytes:
    pad_len = data[-1]
    return data[:-pad_len]

def af_alg_encrypt_aes_cbc(data, key, iv):
    """
    使用内核AF_ALG进行AES-CBC加密
    """
    # 创建AF_ALG套接字
    ALG_TYPE_SKCIPHER = 0x0a
    sock = socket.socket(socket.AF_ALG, socket.SOCK_SEQPACKET, 0)
    
    # 绑定加密算法
    salg_type = 'skcipher'
    salg_name = "cbc(aes)"
    
    sock.bind((salg_type, salg_name))
    
    # 设置密钥
    sock.setsockopt(socket.SOL_ALG, socket.ALG_SET_KEY, key)
    sock1, _ = sock.accept()
    
    # 执行加密
    plaintext = pkcs7_pad(data)
    sock1.sendmsg_afalg([plaintext], op=socket.ALG_OP_ENCRYPT, iv=iv)
    
    # 接收密文
    size = len(plaintext)
    data = b''
    while len(data) < size:
        chunk = sock1.recv(size - len(data))
        if not chunk:
            break
        data += chunk
    sock1.close()
    sock.close()
    return data

def af_alg_sha256(data):
    """
    使用内核AF_ALG计算SHA256哈希
    """
    ALG_TYPE_HASH = 0x09
    sock = socket.socket(socket.AF_ALG, socket.SOCK_SEQPACKET, 0)
    
    # 绑定哈希算法
    salg_type = 'hash'
    salg_name = 'sha256'
    
    sock.bind((salg_type, salg_name))
    sock1, _ = sock.accept()
    
    # 计算哈希
    sock1.sendall(data)
    hash_result = sock1.recv(1024)
    
    sock1.close()
    sock.close()
    return hash_result

def af_alg_hmac_sha256(data, key):
    """
    使用内核AF_ALG计算HMAC-SHA256
    """
    ALG_TYPE_HASH = 0x09
    sock = socket.socket(socket.AF_ALG, socket.SOCK_SEQPACKET, 0)
    
    # 绑定HMAC算法
    salg_type = 'hash'
    salg_name = "hmac(sha256)"
    
    sock.bind((salg_type, salg_name))
    
    # 设置HMAC密钥
    sock.setsockopt(socket.SOL_ALG, socket.ALG_SET_KEY, key)

    sock1, _ = sock.accept()
    
    # 计算HMAC
    sock1.sendall(data)
    hmac_result = sock1.recv(32)
    
    sock1.close()
    sock.close()
    return hmac_result

# 示例使用
if __name__ == "__main__":
    # 注意：需要root权限运行
    if os.geteuid() != 0:
        pass
        #print("需要root权限运行")
        #exit(1)
    
    # 测试AES-CBC加密
    key = b"0123456789abcdef"  # 16字节AES-128密钥
    iv = b"1234567890abcdef"   # 16字节IV
    plaintext = b"Hello, AF_ALG encryption!"
    
    print("测试AF_ALG AES-CBC加密:")
    encrypted = af_alg_encrypt_aes_cbc(plaintext, key, iv)
    print(f"明文: {plaintext}")
    print(f"加密后: {encrypted.hex()}")
    
    # 测试SHA256
    print("\n测试AF_ALG SHA256:")
    hash_val = af_alg_sha256(plaintext)
    print(f"SHA256: {hash_val.hex()}")
    
    # 测试HMAC-SHA256
    print("\n测试AF_ALG HMAC-SHA256:")
    hmac_key = b"hmac-secret-key-32-bytes-long!!"
    hmac_val = af_alg_hmac_sha256(plaintext, hmac_key)
    print(f"HMAC-SHA256: {hmac_val.hex()}")
