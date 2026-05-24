#!/usr/bin/env python3
"""
Linux AF_ALG 套接字加解密程序
需要：Linux 内核 2.6.38+，支持 AF_ALG
"""

import socket
import os
import struct
import hashlib
from typing import Optional, Tuple, Union
import array
import traceback
from cryptography.hazmat.primitives import padding

def pkcs7_pad(data: bytes, block_size: int = 16) -> bytes:
    pad_len = block_size - (len(data) % block_size)
    return data + bytes([pad_len] * pad_len)

def pkcs7_unpad(data: bytes) -> bytes:
    pad_len = data[-1]
    return data[:-pad_len]

class AFAlgCrypto:
    """
    使用 AF_ALG 套接字进行加解密的类
    支持算法：AES-CBC, AES-ECB, AES-CTR, SHA256, HMAC-SHA256 等
    """
    
    # AF_ALG 常量
    AF_ALG = 38
    SOL_ALG = 279
    
    def __init__(self, algorithm: str = "skcipher", cipher: str = "cbc(aes)"):
        """
        初始化 AF_ALG 加密上下文
        
        Args:
            algorithm: 算法类型，如 "skcipher"(对称加密), "hash"(哈希), "aead"(认证加密)
            cipher: 具体算法，如 "cbc(aes)", "ecb(aes)", "sha256", "hmac(sha256)"
        """
        self.algorithm = algorithm
        self.cipher = cipher
        self.sock = None
        self.tfm_sock = None
        self.connected = False
        
    def __enter__(self):
        """上下文管理器入口"""
        return self
    
    def __exit__(self, exc_type, exc_val, exc_tb):
        """上下文管理器退出，确保资源释放"""
        self.close()
    
    def _create_socket(self) -> socket.socket:
        """创建 AF_ALG 套接字"""
        try:
            return socket.socket(self.AF_ALG, socket.SOCK_SEQPACKET, 0)
        except AttributeError:
            raise RuntimeError("系统不支持 AF_ALG 套接字")
        except OSError as e:
            raise RuntimeError(f"创建 AF_ALG 套接字失败: {e}")
    
    def setup_cipher(self, key: bytes, iv: Optional[bytes] = None) -> None:
        """
        设置加密算法和密钥
        
        Args:
            key: 加密密钥
            iv: 初始化向量（CBC/CTR 模式需要）
        """
        if self.sock or self.tfm_sock:
            self.close()
        
        # 1. 创建主套接字
        self.sock = self._create_socket()
        
        # 2. 绑定算法
        alg_name = f"{self.algorithm}({self.cipher})"
        try:
            self.sock.bind((self.algorithm, self.cipher))
        except OSError as e:
            self.close()
            raise RuntimeError(f"绑定算法 {alg_name} 失败: {e}")
        
        try:
            
            # 设置密钥
            if key:
                self.sock.setsockopt(self.SOL_ALG, 1, key)
        except OSError as e:
            client_sock.close()
            self.close()
            raise RuntimeError(f"设置密钥失败: {e}")
        
        # 5. 接受连接
        self.tfm_sock, _ = self.sock.accept()
        
        # 6. 设置 IV（如果提供）
        if iv and self.tfm_sock:
            self.iv = iv
        self.connected = True
    
    def encrypt(self, plaintext: bytes) -> bytes:
        """
        加密数据
        
        Args:
            plaintext: 明文数据
            
        Returns:
            密文数据
        """
        if not self.connected or not self.tfm_sock:
            raise RuntimeError("未设置加密算法和密钥")
        
        try:
            # 发送明文
            plaintext = pkcs7_pad(plaintext)
            self.tfm_sock.sendmsg_afalg([plaintext], op=socket.ALG_OP_ENCRYPT, iv=self.iv)
            
            # 接收密文
            ciphertext = self._recv_all(self.tfm_sock, len(plaintext))
            return ciphertext
            
        except OSError as e:
            raise RuntimeError(f"加密失败: {e}")
    
    def decrypt(self, ciphertext: bytes) -> bytes:
        """
        解密数据
        
        Args:
            ciphertext: 密文数据
            
        Returns:
            明文数据
        """
        if not self.connected or not self.tfm_sock:
            raise RuntimeError("未设置加密算法和密钥")
        
        try:
            # 对于解密，需要重新设置套接字
            # 关闭旧的变换套接字
            if self.tfm_sock:
                self.tfm_sock.close()
            
            # 接受新的连接
            self.tfm_sock, _ = self.sock.accept()
            
            # 发送密文
            self.tfm_sock.sendmsg_afalg([ciphertext], op=socket.ALG_OP_DECRYPT, iv=self.iv)
            
            # 接收明文
            plaintext = self._recv_all(self.tfm_sock, len(ciphertext))
            plaintext = pkcs7_unpad(plaintext)
            return plaintext
            
        except OSError as e:
            raise RuntimeError(f"解密失败: {e}")
    
    def _recv_all(self, sock: socket.socket, size: int) -> bytes:
        """接收指定大小的数据"""
        data = b''
        while len(data) < size:
            chunk = sock.recv(size - len(data))
            if not chunk:
                break
            data += chunk
        return data
    
    def close(self) -> None:
        """关闭套接字释放资源"""
        if self.tfm_sock:
            self.tfm_sock.close()
            self.tfm_sock = None
        if self.sock:
            self.sock.close()
            self.sock = None
        self.connected = False

class AFAlgHash:
    """使用 AF_ALG 计算哈希"""
    
    def __init__(self, algorithm: str = "sha256"):
        self.algorithm = algorithm
        self.sock = None
        self.tfm_sock = None
        
    def __enter__(self):
        return self
    
    def __exit__(self, exc_type, exc_val, exc_tb):
        self.close()
    
    def setup(self) -> None:
        """设置哈希算法"""
        try:
            # 创建套接字
            self.sock = socket.socket(38, socket.SOCK_SEQPACKET, 0)
            
            # 绑定算法
            alg_name = f"hash({self.algorithm})"
            self.sock.bind(("hash", self.algorithm))
            
            # 接受连接
            self.tfm_sock, _ = self.sock.accept()
        except Exception as e:
            self.close()
            raise RuntimeError(f"设置哈希算法失败: {e}")
    
    def update(self, data: bytes) -> None:
        """更新哈希计算"""
        if not self.tfm_sock:
            raise RuntimeError("哈希算法未设置")
        
        self.tfm_sock.sendall(data)
    
    def digest(self) -> bytes:
        """获取哈希值"""
        if not self.tfm_sock:
            raise RuntimeError("哈希算法未设置")
        
        # 接收哈希值
        hash_value = self.tfm_sock.recv(1024)
        return hash_value
    
    def hexdigest(self) -> str:
        """获取十六进制哈希值"""
        return self.digest().hex()
    
    def close(self) -> None:
        """关闭套接字"""
        if self.tfm_sock:
            self.tfm_sock.close()
            self.tfm_sock = None
        if self.sock:
            self.sock.close()
            self.sock = None

def test_afalg_aes_cbc():
    """测试 AES-CBC 加解密"""
    print("测试 AES-CBC 加解密")
    print("=" * 50)
    
    # 测试数据
    key = os.urandom(32)  # AES-256
    iv = os.urandom(16)   # AES块大小
    plaintext = b"Hello, AF_ALG! This is a test message for AES-CBC encryption." * 4
    
    print(f"密钥: {key.hex()}")
    print(f"IV: {iv.hex()}")
    print(f"明文长度: {len(plaintext)} 字节")
    
    # 使用 AF_ALG 加密
    with AFAlgCrypto("skcipher", "cbc(aes)") as cipher:
        cipher.setup_cipher(key, iv)
        ciphertext = cipher.encrypt(plaintext)
        print(f"密文长度: {len(ciphertext)} 字节")
        
        # 重置 IV 进行解密
        cipher.setup_cipher(key, iv)
        decrypted = cipher.decrypt(ciphertext)
        
        print(f"解密成功: {decrypted == plaintext}")
        
        if decrypted != plaintext:
            print("解密失败！")
            print(f"原始明文: {plaintext[:50]}...")
            print(f"解密结果: {decrypted[:50]}...")

def test_afalg_hash():
    """测试哈希计算"""
    print("\n测试 SHA256 哈希计算")
    print("=" * 50)
    
    data = b"Hello, AF_ALG! This is a test message for hash."
    
    # 使用 AF_ALG 计算哈希
    with AFAlgHash("sha256") as hasher:
        hasher.setup()
        hasher.update(data)
        afalg_hash = hasher.hexdigest()
        print(f"AF_ALG SHA256: {afalg_hash}")
    
    # 使用 hashlib 验证
    lib_hash = hashlib.sha256(data).hexdigest()
    print(f"hashlib SHA256: {lib_hash}")
    print(f"哈希匹配: {afalg_hash == lib_hash}")

def benchmark_afalg_vs_openssl():
    """性能对比：AF_ALG vs 纯 Python"""
    print("\n性能对比测试")
    print("=" * 50)
    
    import time
    from cryptography.hazmat.primitives.ciphers import Cipher, algorithms, modes
    from cryptography.hazmat.backends import default_backend
    
    # 准备测试数据
    key = os.urandom(32)
    iv = os.urandom(16)
    data = os.urandom(1024 * 1024)
    
    # 测试 AF_ALG
    print("测试 AF_ALG 加密速度...")
    start_time = time.time()
    
    with AFAlgCrypto("skcipher", "cbc(aes)") as cipher:
        cipher.setup_cipher(key, iv)
        ciphertext = cipher.encrypt(data)
    
    afalg_time = time.time() - start_time
    print(f"AF_ALG 加密数据耗时: {afalg_time:.4f} 秒")
    print(f"AF_ALG 吞吐量: {len(data) / afalg_time / 1024 / 1024:.2f} MB/s")
    
    # 测试 cryptography 库
    print("\n测试 cryptography 库加密速度...")
    start_time = time.time()
    
    padder = padding.PKCS7(algorithms.AES.block_size).padder()
    padded_data = padder.update(data) + padder.finalize()

    cipher = Cipher(algorithms.AES(key), modes.CBC(iv), backend=default_backend())
    encryptor = cipher.encryptor()
    ciphertext2 = encryptor.update(padded_data) + encryptor.finalize()
    
    crypto_time = time.time() - start_time
    print(f"cryptography 加密数据耗时: {crypto_time:.4f} 秒")
    print(f"cryptography 吞吐量: {len(data) / crypto_time / 1024 / 1024:.2f} MB/s")
    
    # 验证结果一致
    print(f"\n加密结果一致: {ciphertext == ciphertext2}")

def list_available_algorithms():
    """列出系统支持的加密算法"""
    print("系统支持的 AF_ALG 算法")
    print("=" * 50)
    
    try:
        # 创建套接字查询可用算法
        sock = socket.socket(38, socket.SOCK_SEQPACKET, 0)
        
        # 尝试绑定不同类型的算法
        algorithms_to_try = [
            "skcipher",
            "aead", 
            "hash",
            "rng"
        ]
        
        for alg_type in algorithms_to_try:
            print(f"\n{alg_type.upper()} 算法:")
            print("-" * 30)
            
            # 读取算法列表文件
            alg_path = f"/proc/crypto"
            if os.path.exists(alg_path):
                with open(alg_path, 'r') as f:
                    content = f.read()
                    # 简单提取算法名
                    import re
                    patterns = {
                        'skcipher': r'name\s*:\s*(cbc|ecb|ctr)\([a-z0-9]+\)',
                        'hash': r'name\s*:\s*(sha|md)[a-z0-9]+',
                        'aead': r'name\s*:\s*(gcm|ccm)\([a-z0-9]+\)',
                    }
                    
                    if alg_type in patterns:
                        matches = re.findall(patterns[alg_type], content, re.IGNORECASE)
                        for match in set(matches):
                            print(f"  - {match}")
            
    except Exception as e:
        print(f"查询算法失败: {e}")
    finally:
        if 'sock' in locals():
            sock.close()

def main():
    """主函数"""
    print("Linux AF_ALG 套接字加解密演示程序")
    print("=" * 60)
    
    # 检查是否支持 AF_ALG
    try:
        socket.socket(38, socket.SOCK_SEQPACKET, 0)
        print("✓ 系统支持 AF_ALG")
    except:
        print("✗ 系统不支持 AF_ALG")
        print("请确保：")
        print("1. Linux 内核版本 >= 2.6.38")
        print("2. 已加载 cryptd 模块: sudo modprobe cryptd")
        print("3. 已加载 aesni_intel 模块: sudo modprobe aesni_intel")
        return
    
    # 运行测试
    try:
        # 列出可用算法
        list_available_algorithms()
        
        # 测试加解密
        test_afalg_aes_cbc()
        
        # 测试哈希
        test_afalg_hash()
        
        # 性能测试
        benchmark_afalg_vs_openssl()
        
    except Exception as e:
        traceback.print_exc()
        print(f"\n错误: {e}")
        print("\n可能的原因：")
        print("1. 内核模块未加载，尝试运行:")
        print("   sudo modprobe cryptd")
        print("   sudo modprobe aesni_intel")
        print("2. 权限不足，尝试以 root 运行")
        print("3. 内核不支持某些算法")

if __name__ == "__main__":
    main()
