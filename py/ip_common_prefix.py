#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
读取 IPv4 地址列表，计算所有地址的共同前缀（子网掩码/网段）
"""

import sys
import ipaddress
from typing import List, Tuple, Optional


def read_ips_from_file(filepath: str) -> List[ipaddress.IPv4Address]:
    """从文件读取 IPv4 地址，每行一个"""
    ips = []
    with open(filepath, 'r', encoding='utf-8') as f:
        for line_num, line in enumerate(f, start=1):
            line = line.strip()
            if not line or line.startswith('#'):
                continue
            try:
                ip = ipaddress.IPv4Address(line)
                ips.append(ip)
            except ValueError as e:
                print(f"[警告] 第 {line_num} 行地址无效，已跳过: {line} ({e})", file=sys.stderr)
    return ips


def find_common_prefix(ips: List[ipaddress.IPv4Address]) -> Optional[Tuple[ipaddress.IPv4Address, int]]:
    """
    查找所有 IPv4 地址的最长共同前缀。
    返回: (网络地址, 前缀长度)
    """
    if not ips:
        return None
    if len(ips) == 1:
        return (ipaddress.IPv4Address(int(ips[0])), 32)

    # 转换为 32 位整数
    int_ips = [int(ip) for ip in ips]

    # 找最长共同前缀长度
    prefix_length = 0
    for i in range(32):
        mask = 0x80000000 >> i
        bit = int_ips[0] & mask
        if all((ip & mask) == bit for ip in int_ips):
            prefix_length += 1
        else:
            break

    # 计算网络地址（所有 IP 按前缀长度掩码后的共同部分）
    mask = 0xFFFFFFFF << (32 - prefix_length)
    network_int = int_ips[0] & mask
    network = ipaddress.IPv4Address(network_int)

    return (network, prefix_length)


def print_result(network: ipaddress.IPv4Address, prefix: int, total: int):
    """输出计算结果"""
    net = ipaddress.IPv4Network(f"{network}/{prefix}", strict=False)
    mask = net.netmask

    print(f"\n{'='*50}")
    print(f"  共读取 {total} 个 IPv4 地址")
    print(f"{'='*50}")
    print(f"  共同前缀长度: /{prefix}")
    print(f"  子网掩码:     {mask}")
    print(f"  网络地址:     {network}/{prefix}")
    print(f"  可用主机数:   {net.num_addresses}")
    print(f"  地址范围:     {net[0]} - {net[-1]}")
    print(f"{'='*50}\n")


def main():
    if len(sys.argv) < 2:
        print("用法: python ip_common_prefix.py <ip文件路径>", file=sys.stderr)
        print("示例: python ip_common_prefix.py ips.txt", file=sys.stderr)
        print("\n文件格式: 每行一个 IPv4 地址", file=sys.stderr)
        sys.exit(1)

    filepath = sys.argv[1]

    try:
        ips = read_ips_from_file(filepath)
    except FileNotFoundError:
        print(f"[错误] 文件不存在: {filepath}", file=sys.stderr)
        sys.exit(1)
    except Exception as e:
        print(f"[错误] 读取文件失败: {e}", file=sys.stderr)
        sys.exit(1)

    if not ips:
        print("[错误] 文件中没有有效的 IPv4 地址", file=sys.stderr)
        sys.exit(1)

    result = find_common_prefix(ips)
    if result:
        network, prefix = result
        print_result(network, prefix, len(ips))


if __name__ == "__main__":
    main()
