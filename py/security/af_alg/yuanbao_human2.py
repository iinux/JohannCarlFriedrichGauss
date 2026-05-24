#!/usr/bin/env python3
"""
简化的 AF_ALG 加解密封装
"""

import socket
import os
import hashlib

def pkcs7_pad(data, block_size=16):
    pad_len = block_size - (len(data) % block_size)
    return data + bytes([pad_len]) * pad_len


def pkcs7_unpad(data):
    return data[:-data[-1]]

class SimpleAFAlgCrypto:
    """简化的 AF_ALG 加解密类"""
    
    AF_ALG = 38
    SOL_ALG = 279
    
    @staticmethod
    def encrypt_aes_cbc(plaintext: bytes, key: bytes, iv: bytes) -> bytes:
        """
        AES-CBC 加密
        
        Args:
            plaintext: 明文
            key: 密钥 (16, 24, 32字节对应AES-128/192/256)
            iv: 初始化向量 (16字节)
            
        Returns:
            密文
        """
        
        try:
            # 绑定算法
            
            # 设置密钥
            client = socket.socket(SimpleAFAlgCrypto.AF_ALG, socket.SOCK_SEQPACKET, 0)
            client.bind(("skcipher", "cbc(aes)"))
            client.setsockopt(SimpleAFAlgCrypto.SOL_ALG, socket.ALG_SET_KEY, key)
            
            # 接受连接
            tfm, _ = client.accept()
            
            
            # 加密
            plaintext = pkcs7_pad(plaintext)
            tfm.sendmsg_afalg([plaintext], op=socket.ALG_OP_ENCRYPT, iv=iv)
            
            # 接收密文
            ciphertext = tfm.recv(len(plaintext))
            return ciphertext
            
        finally:
            pass
    
    @staticmethod
    def decrypt_aes_cbc(ciphertext: bytes, key: bytes, iv: bytes) -> bytes:
        """
        AES-CBC 解密
        
        Args:
            ciphertext: 密文
            key: 密钥
            iv: 初始化向量
            
        Returns:
            明文
        """
        
        try:
            # 设置密钥
            client = socket.socket(SimpleAFAlgCrypto.AF_ALG, socket.SOCK_SEQPACKET, 0)
            client.bind(("skcipher", "cbc(aes)"))
            client.setsockopt(SimpleAFAlgCrypto.SOL_ALG, 1, key)
            
            # 接受连接
            tfm, _ = client.accept()
            
            
            # 解密
            tfm.sendmsg_afalg([ciphertext], op=socket.ALG_OP_DECRYPT, iv=iv)
            
            # 接收明文
            plaintext = b''
            decrypted = tfm.recv(len(ciphertext))
            plaintext = pkcs7_unpad(decrypted)
            
            return plaintext
            
        finally:
            pass
    
    @staticmethod
    def hash_sha256(data: bytes) -> bytes:
        """
        计算 SHA256 哈希
        
        Args:
            data: 输入数据
            
        Returns:
            哈希值
        """
        
        try:
            # 连接
            client = socket.socket(SimpleAFAlgCrypto.AF_ALG, socket.SOCK_SEQPACKET, 0)
            client.bind(("hash", "sha256"))
            
            # 接受连接
            tfm, _ = client.accept()
            
            # 计算哈希
            tfm.sendall(data)
            
            # 接收哈希
            hash_value = tfm.recv(1024)
            
            return hash_value
            
        finally:
            pass

# 使用示例
if __name__ == "__main__":
    # 测试数据
    key = b"SixteenByteKey!!"  # 16字节 AES-128
    iv = b"SixteenByteIV123"  # 16字节 IV
    plaintext = b"Hello, World! This is a secret message."
    
    print("简化的 AF_ALG 加解密演示")
    print("=" * 50)
    
    # 加密
    ciphertext = SimpleAFAlgCrypto.encrypt_aes_cbc(plaintext, key, iv)
    print(f"明文: {plaintext}")
    print(f"密文: {ciphertext.hex()}")
    
    # 解密
    decrypted = SimpleAFAlgCrypto.decrypt_aes_cbc(ciphertext, key, iv)
    print(f"解密: {decrypted}")
    print(f"加解密成功: {decrypted == plaintext}")
    
    # 哈希
    hash_val = SimpleAFAlgCrypto.hash_sha256(plaintext)
    print(f"\nSHA256 哈希: {hash_val.hex()}")
    
    # 验证
    import hashlib
    expected = hashlib.sha256(plaintext).digest()
    print(f"验证哈希: {hash_val == expected}")
