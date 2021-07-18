-- https://www.cnblogs.com/0pen1/p/12592906.html
-- https://item.taobao.com/item.htm?spm=a1z09.2.0.0.1cfd2e8dmPHEaE&id=531755241333&_u=aoohcvued60
-- https://esp8266.ru/downloads/
-- https://nodemcu-build.com/index.php
-- https://github.com/nodemcu/nodemcu-flasher
-- https://www.jianshu.com/p/3aba3ce1ad12
-- http://www.lcdwiki.com/zh/1.44inch_SPI_Arduino_Module_Black_SKU:MAR1442
-- https://sankios.imediabank.com/nodemcu-oled-ssd1306
-- https://blog.csdn.net/weixin_42268054/article/details/104254955
-- https://lceda.cn/
-- https://www.ai-thinker.com/home
-- https://www.espressif.com/
-- https://nodemcu.readthedocs.io/en/release/upload/
-- https://nodemcu.readthedocs.io/en/release/lua-developer-faq/#how-do-i-avoid-a-panic-loop-in-initlua
-- http://httpbin.org/ip
-- http://nodemcu-dev.doit.am/2.html

print(wifi.sta.getip())
wifi.setmode(wifi.STATION)
cfg={}
cfg.ssid="xxx"
cfg.pwd="yyy"
wifi.sta.config(cfg)

tmr.delay(10*1000000)
print(wifi.sta.getip())
