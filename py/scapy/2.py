from scapy.all import *

# 创建IPv6扩展头
ipv6_extension_header = IPv6ExtHdr()
ipv6_extension_header['version'] = 0  # 版本号
ipv6_extension_header['length'] = len(payload)  # 扩展头长度
ipv6_extension_header['type'] = XCAP_TYPE_DATA  # 扩展头类型

# 构建IPv6报文
ipv6_packet = IPv6() / ipv6_extension_header / payload  # 注意替换payload为实际的数据部分

# 发送IPv6报文
send(ipv6_packet)

if __name__ == '__main__':
    pass
