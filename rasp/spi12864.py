#!/usr/bin/python/
# coding: utf-8
import time
import Adafruit_GPIO.SPI as SPI
import Adafruit_SSD1306
from PIL import Image
from PIL import ImageDraw
from PIL import ImageFont

# https://www.chenxublog.com/2016/08/20/raspi-12864-spi-oled.html
# https://mc.dfrobot.com.cn/thread-13396-1-1.html

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
#font = ImageFont.truetype('/home/pi/MingImperial.TTF', 8)

# Write two lines of text.
draw.text((x, top), 'This is first line', font=font, fill=255)
draw.text((x, top+13), 'This is second line', font=font, fill=255)
draw.text((x, top+23), 'This is third line', font=font, fill=255)
draw.text((x, top+33), 'This is fourth line', font=font, fill=255)
draw.text((x, top+43), 'This is fifth line', font=font, fill=255)
draw.text((x, top+53), 'This is last line', font=font, fill=255)

# Display image.
disp.image(image)
disp.display()
