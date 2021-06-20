from machine import Pin
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


led_and_sound()
