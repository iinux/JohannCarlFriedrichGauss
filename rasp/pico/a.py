from machine import Pin,PWM
import utime

led = Pin(25, Pin.OUT)
sound = Pin(0, Pin.OUT)
o = Pin(1, Pin.OUT)
i = Pin(2, Pin.IN)


def led_and_sound():
    while (1):
        led.value(1)
        sound.value(1)
        utime.sleep(1)
        led.value(0)
        sound.value(0)
        utime.sleep(1)


        
def keyboard():
    count = 0
    o.value(1)
    while(1):
        utime.sleep_ms(1000)
        count+=1
        print(str(count) + ":" + str(i.value()))
        if (i.value() == 1):
            sound.value(0)
            led.value(1)
        else:
            sound.value(1)    
            led.value(0)


PWM_PulseWidth=0
#使用树莓派Pico板上LED，构建PWM对象pwm_LED
pwm_LED=PWM(Pin(25))
#设置pwm_LED频率
pwm_LED.freq(500)
while True:
    while PWM_PulseWidth<65535:
        PWM_PulseWidth=PWM_PulseWidth+50
        utime.sleep_ms(1)   #延时1ms
        pwm_LED.duty_u16(PWM_PulseWidth)
    while PWM_PulseWidth>0:
        PWM_PulseWidth=PWM_PulseWidth-50
        utime.sleep_ms(1)
        pwm_LED.duty_u16(PWM_PulseWidth)
