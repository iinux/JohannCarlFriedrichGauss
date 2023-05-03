#coding=utf8
import threading
import time
import sys
import signal

def showperson(name):
    while True:
       time.sleep(1)
       print('show person :%s'%name)

print('%s thread start!'%(time.ctime()))

def quit(signal_num,frame):
    print("you stop the threading")
    sys.exit()

try:
   signal.signal(signal.SIGINT, quit)
   signal.signal(signal.SIGTERM, quit)

   list=[]
   for i in range(3):
       t =threading.Thread(target=showperson,args=("person-%d"%i,))
       list.append(t)
       t.setDaemon(True)
       t.start()

   while True:
       pass
except Exception as e:
    print(e)

print('%s thread end!'%(time.ctime()))
