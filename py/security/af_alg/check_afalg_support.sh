#!/bin/bash
# check_afalg_support.sh
# 检查系统是否支持 AF_ALG

echo "=== 检查 AF_ALG 支持 ==="

# 1. 检查内核版本
echo -n "1. 内核版本: "
uname -r
KERNEL_MAJOR=$(uname -r | cut -d. -f1)
KERNEL_MINOR=$(uname -r | cut -d. -f2)
if [ $KERNEL_MAJOR -gt 2 ] || ([ $KERNEL_MAJOR -eq 2 ] && [ $KERNEL_MINOR -ge 6 ] && [ $(uname -r | cut -d. -f3 | cut -d- -f1) -ge 38 ]); then
    echo "   ✓ 内核版本 >= 2.6.38"
else
    echo "   ✗ 内核版本过低，需要 >= 2.6.38"
fi

# 2. 检查 AF_ALG 套接字
echo -n "2. 检查 AF_ALG 套接字: "
python3 -c "
import socket
try:
    socket.socket(38, socket.SOCK_SEQPACKET, 0)
    print('✓ 支持 AF_ALG')
except:
    print('✗ 不支持 AF_ALG')
"

# 3. 检查已加载的内核模块
echo "3. 已加载的加密模块:"
lsmod | grep -E "(aes|crypt|sha)" | head -10 || echo "   未找到相关模块"

# 4. 检查可用的算法
echo "4. 可用加密算法:"
if [ -f /proc/crypto ]; then
    grep -E "name\s*:\s*(cbc|ecb|ctr|gcm|sha|md)" /proc/crypto | sort | uniq | head -20
else
    echo "   /proc/crypto 不存在"
fi

# 5. 加载必要的模块
echo "5. 建议加载的模块:"
echo "   sudo modprobe cryptd"
echo "   sudo modprobe aesni_intel"
echo "   sudo modprobe sha256_ssse3"
