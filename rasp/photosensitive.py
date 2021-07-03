import RPi.GPIO as GPIO
import time

# https://zhuanlan.zhihu.com/p/40594358
# https://shumeipai.nxez.com/raspberry-pi-pins-version-40
# https://www.jianshu.com/p/f31b90cc756f

no = 18

GPIO.setmode(GPIO.BCM)
GPIO.setup(no,GPIO.IN)

def m1():
    while True:
        time.sleep(1)
        if (GPIO.input(no)):
            print("####################################################################################################")
        else:
            print("##################################################")


def m2():
    while True:
        channel = GPIO.wait_for_edge(no, GPIO.RISING, timeout=5000)
        if channel is None:
            print('Timeout occurred')
        else:
            print('Edge detected on channel', channel)


def m3():
    GPIO.add_event_detect(no, GPIO.RISING)  # add rising edge detection on a channel
    while True:
        time.sleep(1)
        if GPIO.event_detected(no):
            print('Button pressed')
        else:
            print('Button not pressed')


def my_callback_one(channel):
    print('Callback one')

def my_callback_two(channel):
    print('Callback two')

def m4():
    GPIO.add_event_detect(no, GPIO.RISING)
    GPIO.add_event_callback(no, my_callback_one)
    GPIO.add_event_callback(no, my_callback_two)
    while True:
        time.sleep(1)

m1()


GPIO.cleanup()
