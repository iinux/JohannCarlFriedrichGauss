import requests
import traffic_card_data

url = 'http://www.8989iot.com/api/Card/loginCard'
user_agent = 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/536.36 (KHTML, like Gecko) Chrome/105.0.0.0 ' \
             'Safari/536.36 '


def req(number):
    res = requests.post(url, json={
        'number': number
    }, headers={
        'Content-Type': 'application/json',
        'User-Agent': user_agent,
    })

    res_data = res.json()['data']
    # print(res_data)
    print('-' * 40)
    print('number: ' + number)
    print('usage: %s G / %s G' % (res_data['remainAmount'] / 1024, res_data['totalAmount'] / 1024))
    print('packageName: ' + res_data['packageName'])
    print('expiretime: ' + res_data['expiretime'])
    print('status: ' + res_data['status'])
    print('-' * 40)


if __name__ == '__main__':
    req(traffic_card_data.number1)
    req(traffic_card_data.number2)
    pass
