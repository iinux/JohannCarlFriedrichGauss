import sys
import time
from datetime import datetime

import requests

import hyf_data

user_agent = 'Mozilla/5.0 (Linux; Android 12; Redmi 5 Plus Build/SQ3A.220705.004; wv) ' \
             'AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/86.0.4240.99 XWEB/4433 MMWEBSDK/20230504 ' \
             'Mobile Safari/537.36 MMWEBID/1289 MicroMessenger/8.0.37.2368(0x28002548) WeChat/arm64 Weixin ' \
             'GPVersion/1 NetType/WIFI Language/en ABI/arm64 MiniProgramEnv/android'


def print_x(text):
    print(datetime.now(), text)


def req(sn, port_num):
    res = requests.get(hyf_data.url2, params={
        'gtel': sn
    }, headers={
        'Content-Type': 'text/html;charset=utf-8',
        'User-Agent': user_agent,
    })

    if res.json()[0]['glzt%d' % port_num] == 1:
        return 2
    else:
        return 1


if __name__ == '__main__':
    port = 8
    sn = '68000026956'
    if len(sys.argv) > 1:
        port = int(sys.argv[1])
    if len(sys.argv) > 2:
        sn = sys.argv[2]
    print_x('sn=%s, port=%d' % (sn, port))
    status = req(sn, port)
    if status == 1:
        print_x('已经完成')
        hyf_data.notify('已经完成')
    elif status == 2:
        while True:
            time.sleep(300)
            status = req(sn, port)
            if status == 1:
                print_x('完成')
                hyf_data.notify('完成')
                break

    pass
