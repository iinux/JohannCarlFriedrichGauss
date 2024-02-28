import hyf_data
import requests
import smtplib
import sys
import time
from email.mime.text import MIMEText
from datetime import datetime, timedelta

user_agent = 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/536.36 (KHTML, like Gecko) Chrome/105.0.0.0 ' \
             'Safari/536.36 '


def print_x(text):
    print(datetime.now(), text)


def req(sn, port_num):
    res = requests.post(hyf_data.url, data={
        'sn': sn
    }, headers={
        'Content-Type': 'application/x-www-form-urlencoded',
        'User-Agent': user_agent,
    })

    res_data = res.json()['data']
    data = res_data['data']
    port_data = data['port_data']
    for port in port_data:
        if port['port'] == port_num:
            return port['status']
    print_x('error ' + res.text)
    return 0


def send(content, body='body'):
    my_email = smtplib.SMTP(hyf_data.smtp_server, hyf_data.smtp_port)
    my_email.login(hyf_data.from_email, hyf_data.from_email_password)

    msg = MIMEText(body, 'plain', 'utf-8')
    # msg = email.mime.text.MIMEText(content,_subtype='plain')
    msg['to'] = hyf_data.to_email
    msg['from'] = hyf_data.from_email
    msg['subject'] = content

    my_email.sendmail(hyf_data.from_email, [hyf_data.to_email], msg.as_string())
    my_email.quit()


if __name__ == '__main__':
    port = 5
    sn = '1090527'
    if len(sys.argv) > 1:
        port = sys.argv[1]
    if len(sys.argv) > 2:
        sn = sys.argv[2]
    print_x('sn=%s, port=%d' % (sn, port))
    status = req(sn, port)
    if status == 1:
        print_x('已经完成')
        send('已经完成')
    elif status == 2:
        while True:
            time.sleep(300)
            status = req(sn, port)
            if status == 1:
                print_x('完成')
                send('完成')
                break

    pass
