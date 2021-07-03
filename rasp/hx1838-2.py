#!/usr/bin/python
# -*- coding:utf-8 -*-  HX1838
import RPi.GPIO as GPIO
import time

ir_pin = 14;

GPIO.setmode(GPIO.BCM)
#GPIO.setup(ir_pin,GPIO.IN,GPIO.PUD_UP)
GPIO.setup(ir_pin, GPIO.IN,  pull_up_down = GPIO.PUD_UP)
GPIO.add_event_detect(ir_pin, GPIO.FALLING,  bouncetime = 200)

def ir_1838(ir_pin):
    ir_key=""
    if GPIO.input(ir_pin) == 0:
        
        count = 0
        while GPIO.input(ir_pin) == 0 and count < 200:
            count += 1
            time.sleep(0.00006)
        count = 0

        while GPIO.input(ir_pin) == 1 and count < 80:
            count += 1
            time.sleep(0.00006)
        
        idx = 0
        cnt = 0
        data = [0,0,0,0]
        for i in range(0,32):
            count = 0
            while GPIO.input(ir_pin) == 0 and count < 15:
                count += 1
                time.sleep(0.00006)

            count = 0
            while GPIO.input(ir_pin) == 1 and count < 40:
                count += 1
                time.sleep(0.00006)

            if count > 8:
                data[idx] |= 1<<cnt
            if cnt == 7:
                cnt = 0
                idx += 1
            else:
                cnt += 1
        if data[0]+data[1] == 0xFF and data[2]+data[3] == 0xFF:
            #print("Get the key: 0x%02x" %data[2])
            if  (data[2]==0x45):
                ir_key="1"
            elif(data[2]==0x46):
                ir_key="2"
            elif(data[2]==0x47):
                ir_key="3"
            elif(data[2]==0x44):
                ir_key="4"
            elif(data[2]==0x40):
                ir_key="5"
            elif(data[2]==0x43):
                ir_key="6"
            elif(data[2]==0x07):
                ir_key="7"
            elif(data[2]==0x15):
                ir_key="8"
            elif(data[2]==0x09):
                ir_key="9"
            elif(data[2]==0x16):
                ir_key="*"
            elif(data[2]==0x19):
                ir_key="0"
            elif(data[2]==0x0d):
                ir_key="#"
            elif(data[2]==0x18):
                ir_key="上"
            elif(data[2]==0x52):
                ir_key="下"
            elif(data[2]==0x08):
                ir_key="左"
            elif(data[2]==0x5a):
                ir_key="右"
            elif(data[2]==0x1c):
                ir_key="OK"
            print("ir_1838()检测到按键：  "+ir_key)
    #return ir_key

GPIO.add_event_callback(ir_pin, ir_1838)  #开始运行ir_1838()， 以便随时响应遥控器按键

############################
# 主程序运行中
############################
print("等待中，请按下遥控器按钮......")
i=1
try:
    while True:
        print("主程序运行中，第"+str(i)+"次")
        time.sleep(5)
        i=i+1        
except KeyboardInterrupt:
    GPIO.cleanup();
