#!/usr/bin/env python
# -*- coding: utf8 -*-
 
import MFRC522
 
rc = MFRC522.MFRC522()
# save src card data, then write the saved date to des card
dataBlock0 = [0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15]
sectorkey = [0xff, 0xff, 0xff, 0xff, 0xff, 0xff]
 
def read_block0():
    (status, tagType) = rc.MFRC522_Request(rc.PICC_REQIDL)
    if status != rc.MI_OK:
        print("MFRC522_Request error " + str(status))
        return
    print("MFRC522_Request success, tagType: %#x %#x" %(tagType[0], tagType[1]))
 
    (status, Uid) = rc.MFRC522_Anticoll()
    if status != rc.MI_OK:
        print("MFRC522_Anticoll error")
        return
    print("MFRC522_Anticoll success, Uid: %#x %#x %#x %#x" %(Uid[0], Uid[1], Uid[2], Uid[3]))
 
    status = rc.MFRC522_Select(Uid)
    if status != rc.MI_OK:
        print("MFRC522_Select error")
        return
    print("MFRC522_Select success")
 
    status = rc.MFRC522_Auth(rc.PICC_AUTHENT1A, 1, sectorkey, Uid)
    if status != rc.MI_OK:
        print("MFRC522_Auth error")
        return
    print("MFRC522_Auth success")
 
    (status, dataBlock0) = rc.MFRC522_ReadBolock(0)
    if status != rc.MI_OK:
        print("MFRC522_ReadBolock error")
        return
    print("MFRC522_ReadBolock success, block0: %#x %#x %#x %#x %#x %#x %#x %#x %#x %#x %#x %#x %#x %#x %#x %#x" 
          %(dataBlock0[0], dataBlock0[1], dataBlock0[2], dataBlock0[3], dataBlock0[4], dataBlock0[5], dataBlock0[6], dataBlock0[7], 
            dataBlock0[8], dataBlock0[9], dataBlock0[10], dataBlock0[11], dataBlock0[12], dataBlock0[13], dataBlock0[14], dataBlock0[15]))
 
    rc.MFRC522_StopCrypto1()
 
 
def write_block0():
    (status, tagType) = rc.MFRC522_Request(rc.PICC_REQIDL)
    if status != rc.MI_OK:
        print("MFRC522_Request error")
        return
    print("MFRC522_Request success, tagType: %#x %#x" %(tagType[0], tagType[1]))
 
    (status, Uid) = rc.MFRC522_Anticoll()
    if status != rc.MI_OK:
        print("MFRC522_Anticoll error")
        return
    print("MFRC522_Anticoll success, Uid: %#x %#x %#x %#x" %(Uid[0], Uid[1], Uid[2], Uid[3]))
 
    status = rc.MFRC522_Select(Uid)
    if status != rc.MI_OK:
        print("MFRC522_Select error")
        return
    print("MFRC522_Select success")
 
    status = rc.MFRC522_Halt()
    print("MFRC522_Halt %d" %status)
 
    status = rc.MFRC522_WriteCmd40()
    if (status != rc.MI_OK):
        print("MFRC522_ToCard 0x40 error, status:%d" %status)
        return
    print("MFRC522_ToCard 0x40 success")
 
    status = rc.MFRC522_WriteCmd43()
    if (status != rc.MI_OK):
        print("MFRC522_ToCard 0x43 error, status:%d" %status)
        return
    print("MFRC522_ToCard 0x43 success")
 
    status = rc.MFRC522_WriteBlock(0, dataBlock0)
    if status != rc.MI_OK:
        print("MFRC522_WriteBlock error")
        return
    print("MFRC522_WriteBlock success")
 
    rc.MFRC522_StopCrypto1()  
 
def closeSPI():
    rc.MFRC522_CloseSPI()
 
 
if __name__ == '__main__':
    while True:
        action = input("Enter action r/w: ")
        if action == 'r':
            read_block0()
        elif action == 'w':
            write_block0()
        elif action == 'q':
            closeSPI()
            break
    print("exit procedure")
