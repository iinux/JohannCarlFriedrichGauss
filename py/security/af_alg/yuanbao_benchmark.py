#!/usr/bin/env python3
"""
AF_ALG 与 cryptography 库性能对比
"""

import time
import statistics
import os
from cryptography.hazmat.primitives.ciphers import Cipher, algorithms, modes
from cryptography.hazmat.backends import default_backend

def benchmark_afalg():
    """AF_ALG 性能测试"""
    from yuanbao_human2 import SimpleAFAlgCrypto
    
    # 准备测试数据
    key = os.urandom(32)  # AES-256
    iv = os.urandom(16)
    data_sizes = [1024, 4096, 16384, 65536, 262144]  # 1KB 到 256KB
    
    print("AF_ALG 性能测试")
    print("=" * 60)
    
    for size in data_sizes:
        data = os.urandom(size)
        times = []
        
        # 预热
        SimpleAFAlgCrypto.encrypt_aes_cbc(data[:1024], key, iv)
        
        # 运行 10 次取平均
        for _ in range(10):
            start = time.perf_counter()
            ciphertext = SimpleAFAlgCrypto.encrypt_aes_cbc(data, key, iv)
            end = time.perf_counter()
            times.append(end - start)
        
        avg_time = statistics.mean(times)
        throughput = size / avg_time / 1024 / 1024  # MB/s
        
        print(f"数据大小: {size:8d} bytes | "
              f"平均时间: {avg_time*1000:6.2f} ms | "
              f"吞吐量: {throughput:6.2f} MB/s")

def benchmark_cryptography():
    """cryptography 库性能测试"""
    
    # 准备测试数据
    key = os.urandom(32)
    iv = os.urandom(16)
    data_sizes = [1024, 4096, 16384, 65536, 262144]
    
    print("\ncryptography 库性能测试")
    print("=" * 60)
    
    for size in data_sizes:
        data = os.urandom(size)
        times = []
        
        # 预热
        cipher = Cipher(algorithms.AES(key), modes.CBC(iv), backend=default_backend())
        encryptor = cipher.encryptor()
        encryptor.update(data[:1024]) + encryptor.finalize()
        
        # 运行 10 次取平均
        for _ in range(10):
            start = time.perf_counter()
            
            cipher = Cipher(algorithms.AES(key), modes.CBC(iv), backend=default_backend())
            encryptor = cipher.encryptor()
            ciphertext = encryptor.update(data) + encryptor.finalize()
            
            end = time.perf_counter()
            times.append(end - start)
        
        avg_time = statistics.mean(times)
        throughput = size / avg_time / 1024 / 1024  # MB/s
        
        print(f"数据大小: {size:8d} bytes | "
              f"平均时间: {avg_time*1000:6.2f} ms | "
              f"吞吐量: {throughput:6.2f} MB/s")

if __name__ == "__main__":
    benchmark_afalg()
    benchmark_cryptography()
