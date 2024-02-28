#!/usr/bin/python3
import time

import requests
import logging
import my_config
from datetime import datetime, timedelta
from logging.handlers import TimedRotatingFileHandler

# line 110100014013
# station 110100014013023
url = my_config.url

log_formatter = logging.Formatter('%(asctime)s - %(levelname)s - %(message)s')

log_handler = TimedRotatingFileHandler(filename='rb.log', when='midnight', interval=1, backupCount=365, encoding='utf-8')
log_handler.suffix = "%Y%m%d"  # 指定切割后文件名的时间格式

log_handler.setFormatter(log_formatter)
logger = logging.getLogger()
logger.addHandler(log_handler)
logger.setLevel(logging.INFO)
output = open('rb.data', 'a')


def print_info(text):
    print(datetime.now(), text)
    output.write(text + "\n")
    output.flush()


def request():
    res = requests.get(url)
    logging.info(res.text)
    res_json = res.json()
    bus = res_json['buses'][0]
    line = bus['line']
    station = bus['station']
    station_index = bus['station_index']
    start_time = bus['start_time']
    end_time = bus['end_time']
    status = bus['status']
    # selected_car = bus['selected_car']

    trips = bus['trip']
    if len(trips) > 0:
        trip = trips[0]
        arrival = datetime.now() + timedelta(seconds=int(trip['arrival']))
        delay_time = trip['delay_time']
        dis = trip['dis']
        gps_id = trip['gps_id']
        speed = trip['speed']
        uuid = trip['uuid']
        x = trip['x']
        y = trip['y']
        print_info('arriveTime=%s %s %s %s' % (arrival, delay_time, dis, gps_id))


if __name__ == '__main__':
    while True:
        try:
            request()
        except Exception as e:
            print(e)
        time.sleep(30)
