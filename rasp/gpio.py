import RPi.GPIO as GPIO
import time

# https://zhuanlan.zhihu.com/p/40594358
# https://shumeipai.nxez.com/raspberry-pi-pins-version-40

no = 26

GPIO.setmode(GPIO.BCM)
GPIO.setup(no,GPIO.OUT)

while True:
    GPIO.output(no,GPIO.HIGH)
    time.sleep(1)
    GPIO.output(no,GPIO.LOW)
    time.sleep(1)

GPIO.cleanup()
