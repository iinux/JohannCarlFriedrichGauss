import hyf_data
import requests
import sys
import time
from datetime import datetime, timedelta

user_agent = 'Mozilla/5.0 (Linux; Android 12; Redmi 5 Plus Build/SQ3A.220705.004; wv) ' \
             'AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/86.0.4240.99 XWEB/4433 MMWEBSDK/20230504 ' \
             'Mobile Safari/537.36 MMWEBID/1289 MicroMessenger/8.0.37.2368(0x28002548) WeChat/arm64 Weixin ' \
             'GPVersion/1 NetType/WIFI Language/en ABI/arm64 MiniProgramEnv/android'


def print_x(text, notify=False):
    print(datetime.now(), text)
    if notify:
        hyf_data.notify(text)


def req(date, hall_id=2):
    res = requests.get(hyf_data.cstm_url, headers={
        'Content-Type': 'application/json',
        'User-Agent': user_agent,
    })

    res_data = res.json()['data']
    for item in res_data:
        if item['currentDate'] == date:
            for i2 in item['hallTicketPoolVOS']:
                if i2['hallId'] == hall_id:
                    num = i2['ticketPool']
                    print_x('num is %s' % num)
                    return num
    return -1


if __name__ == '__main__':
    date = '2024-09-01'
    hall_id = 2
    if len(sys.argv) > 1:
        date = sys.argv[1]
    if len(sys.argv) > 2:
        hall_id = int(sys.argv[2])
    print_x('date=%s, hall_id=%d' % (date, hall_id))
    while True:
        num = req(date, hall_id)
        if num >= 3:
            print_x('完成', True)
            break
        elif num == -1:
            print_x('error', True)
            break
        time.sleep(30)
