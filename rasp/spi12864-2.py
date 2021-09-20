#!/usr/bin/python/
# coding: utf-8
import time
import Adafruit_GPIO.SPI as SPI
import Adafruit_SSD1306
from PIL import Image
from PIL import ImageDraw
from PIL import ImageFont

import socket
import fcntl
import struct
def get_ip_address(ifname):
    s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    return socket.inet_ntoa(fcntl.ioctl(
        s.fileno(),
        0x8915,  # SIOCGIFADDR
        struct.pack('256s', ifname[:15])
    )[20:24])
	

while True:
	# 打开温度传感器文件
	tfile = open("/sys/bus/w1/devices/28-0115a83f87ff/w1_slave")
	# 读取文件所有内容
	text = tfile.read()
	# 关闭文件
	tfile.close()
	# 用换行符分割字符串成数组，并取第二行
	secondline = text.split("\n")[1]
	# 用空格分割字符串成数组，并取最后一个，即 t=23000
	temperaturedata = secondline.split(" ")[9]
	# 取 t = 后面的数值，并转换为浮点型
	temperature = float(temperaturedata[2:])
	# 转换单位为摄氏度
	temperature = temperature / 1000

	# Raspberry Pi pin configuration:
	RST = 17
	# Note the following are only used with SPI:
	DC = 27
	SPI_PORT = 0
	SPI_DEVICE = 0

	# 128x64 display with hardware SPI:
	disp = Adafruit_SSD1306.SSD1306_128_64(rst=RST, dc=DC, spi=SPI.SpiDev(SPI_PORT, SPI_DEVICE, max_speed_hz=8000000))

	# Initialize library.
	disp.begin()

	# Clear display.
	disp.clear()
	disp.display()

	# Create blank image for drawing.
	# Make sure to create image with mode '1' for 1-bit color.
	width = disp.width
	height = disp.height
	image = Image.new('1', (width, height))

	# Get drawing object to draw on image.
	draw = ImageDraw.Draw(image)

	# Draw a black filled box to clear the image.
	draw.rectangle((0,0,width,height), outline=0, fill=0)

	# Draw some shapes.
	# First define some constants to allow easy resizing of shapes.
	padding = 1
	top = padding
	x = padding
	# Load default font.
	font = ImageFont.load_default()

	# Alternatively load a TTF font.
	# Some other nice fonts to try: http://www.dafont.com/bitmap.php
	#font = ImageFont.truetype('Minecraftia.ttf', 8)

	# Write two lines of text.
	draw.text((x, top), "Chenxu's Raspebrry Pi", font=font, fill=255)
	draw.text((x, top+10), time.strftime("%Y-%m-%d %H:%M:%S",time.localtime(time.time())), font=font, fill=255)
	draw.text((x, top+20), 'The temperature:' + str(temperature), font=font, fill=255)
	draw.text((x, top+30), "Pi's ip:", font=font, fill=255)
	draw.text((x, top+40), 'eth  ip:'+ get_ip_address('eth0'), font=font, fill=255)
	draw.text((x, top+50), 'wlan ip:'+ get_ip_address('lo'), font=font, fill=255)

	# Display image.
	disp.image(image)
	disp.display()
	time.sleep(5)
