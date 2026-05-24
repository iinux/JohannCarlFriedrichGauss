#!/usr/bin/env python3
"""
内核模块加载脚本
需要 root 权限运行
"""

import subprocess
import sys
import os

def load_kernel_modules():
    """加载必要的内核模块"""
    
    modules = [
        "cryptd",          # 加密框架
        "aesni_intel",     # AES-NI 硬件加速
        "sha256_ssse3",    # SHA256 硬件加速
        "algif_skcipher",  # AF_ALG 对称加密
        "algif_hash",      # AF_ALG 哈希
        "algif_aead",      # AF_ALG 认证加密
    ]
    
    print("加载内核模块...")
    for module in modules:
        try:
            result = subprocess.run(
                ["modprobe", module],
                capture_output=True,
                text=True
            )
            if result.returncode == 0:
                print(f"  ✓ 加载 {module}")
            else:
                print(f"  ✗ 加载 {module} 失败: {result.stderr.strip()}")
        except Exception as e:
            print(f"  ✗ 加载 {module} 异常: {e}")

def check_root():
    """检查是否为 root 权限"""
    if os.geteuid() != 0:
        print("需要 root 权限运行此脚本")
        print("请使用: sudo python3", sys.argv[0])
        sys.exit(1)

if __name__ == "__main__":
    check_root()
    load_kernel_modules()
