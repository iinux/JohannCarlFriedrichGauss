# https://mp.weixin.qq.com/s/wx9Tao9HqW3p4jc_kEowIw

import sys

total = {}
print('target file: %s' % sys.argv[1])
with open(sys.argv[1]) as fp:
    while True:
        line = fp.readline()
        if line:
            # print(line)
            fields = line.split('|')
            if len(fields) != 7:
                # print(line)
                continue
            # apnic|CN|ipv6|240f:c000::|24|20190917|allocated
            country = fields[1]
            type = fields[2]
            if type != 'ipv4':
                continue
            ip = fields[3]
            num = int(fields[4])
            # print('ip: %s, num: %d' % (ip, num))
            if country not in total:
                total[country] = 0
            total[country] += num
        else:
            break


print(dict(sorted(total.items(), key=lambda item: item[1], reverse=True)))
