#!/usr/bin/python3

import requests
import dns.message
import sys


def doh_query(domain, qtype='A', server='https://223.5.5.5/dns-query'):
    # 创建 DNS 查询
    query = dns.message.make_query(domain, qtype)
    wire = query.to_wire()

    # 发送请求
    headers = {'Content-Type': 'application/dns-message'}
    response = requests.post(server, data=wire, headers=headers)

    # 解析响应
    dns_response = dns.message.from_wire(response.content)
    return dns_response


if __name__ == '__main__':
    domain = sys.argv[1]
    qtype = sys.argv[2] if len(sys.argv) > 2 else 'A'

    response = doh_query(domain, qtype)
    for rrset in response.answer:
        print(rrset)