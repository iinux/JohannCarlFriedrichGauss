import time
import random
import pyautogui

txt = "I love python"

for _ in range(10):
    pyautogui.typewrite(txt)
    pyautogui.press('enter')
    time.sleep(2)