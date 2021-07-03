import ssd1306
import machine
import time
import uos
import machine

print(uos.uname())
print("Freq: "  + str(machine.freq()) + " Hz")
print("128x64 SSD1306 SPI OLED on Raspberry Pi Pico")

WIDTH = 128
HEIGHT = 64

spi = machine.SPI(0)


CS = machine.Pin(2)
DC = machine.Pin(3)
RES = machine.Pin(4)
oled = ssd1306.SSD1306_SPI(WIDTH, HEIGHT, spi,DC, RES, CS)
oled.fill(0)

oled.text("MicroPython", 0, 0)
oled.text("OLED(ssd1306)", 0, 10)
oled.text("RPi Pico", 0, 20)
oled.show()

while True:
    time.sleep(1)
    oled.invert(1)
    time.sleep(1)
    oled.invert(0)
