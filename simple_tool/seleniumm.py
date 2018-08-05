# coding=utf8

from selenium import webdriver
from wxpy import *
import time


def send_keys_slowly(element, text):
    for t in text:
        time.sleep(1)
        element.send_keys(t)


def click_button(text):
    ret = False
    buttons = driver.find_elements_by_tag_name('button')
    for button in buttons:
        button_text = button.text
        if button_text == text:
            button.click()
            ret = True
            break

    return ret


# driver = webdriver.Firefox()
driver = webdriver.Chrome()
driver.get('')
driver.switch_to.frame('alibaba-login-box')
name = driver.find_element_by_id('fm-login-id')
send_keys_slowly(name, "iinux")
password = driver.find_element_by_id('fm-login-password')
send_keys_slowly(password, '')

click_button(u'登录')
time.sleep(10)

r = click_button(u'签到领积分')
if r:
    print 'ok'
else:
    r = click_button(u'今日已签到')
    if r:
        print 'ok'
    else:
        print 'not ok'
# exit()
# embed()
# print a.get_attribute('href')
