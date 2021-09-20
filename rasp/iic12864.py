# !/usr/bin/env python

import RPi.GPIO as GPIO
import time
import datetime as datetime

from luma.core.interface.serial import i2c, spi
from luma.core.render import canvas
from luma.oled.device import ssd1306, ssd1325, ssd1331, sh1106
from PIL import ImageDraw, Image, ImageFont

GPIO.setmode(GPIO.BCM)
GPIO.setup(18,GPIO.OUT)
GPIO.output(18, GPIO.LOW)

device = sh1106(port=1, address=0x3C)
font = ImageFont.load_default()
fontYear = ImageFont.truetype('/usr/share/fonts/truetype/freefont/FreeMonoBold.ttf', 18)
font2 = ImageFont.truetype('/usr/share/fonts/truetype/freefont/FreeMonoBold.ttf', 16)

def blink():
    GPIO.output(18, GPIO.LOW)
    time.sleep(1)
    GPIO.output(18, GPIO.HIGH)
    time.sleep(1)

def Show(d, fullDt):
    y = fullDt.strftime('%Y-')
    dt = fullDt.strftime('%m-%d')
    tm = fullDt.strftime('%H:%M:%S')
    with canvas(d) as draw:
        draw.text((40, 0), "TIME", font=fontYear, fill=255)
        draw.text((10, 22), y, font=font2, fill=255)
        draw.text((60, 22), dt, font=font2, fill=255)
        draw.text((20, 44), tm, font=font2, fill=255)

def main():
    nowDt = datetime.datetime.now() + datetime.timedelta(hours=8)
    blink()

    while True:
        Show(device, nowDt)
        nowDt = nowDt + datetime.timedelta(seconds=1)
        time.sleep(1)

if __name__ == "__main__":
    main()
