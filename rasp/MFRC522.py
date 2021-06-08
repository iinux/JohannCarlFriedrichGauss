# RC522-pin     Pi-pin
# 3.3V          1(3.3V)
# RST           22(GPIO25)
# GND           6(GND)
# IRQ           none
# MISO          21(GPIO9)
# MOSI          19(GPIO10)
# SCK           23(GPIO11)
# SDA           24(GPIO8)
 
import RPi.GPIO as GPIO
import spi
import signal
import time
 
class MFRC522:
 
    PCD_IDLE        = 0x00     # cancel current command
    PCD_AUTHENT     = 0x0E     # authent
    PCD_RECEIVE     = 0x08     # receive data
    PCD_TRANSMIT    = 0x04     # send data
    PCD_TRANSCEIVE  = 0x0C     # send & receive data
    PCD_RESETPHASE  = 0x0F     # reset
    PCD_CALCCRC     = 0x03     # CRC calculate
 
    PICC_REQIDL     = 0x26     # detect cards, not slepp
    PICC_REQALL     = 0x52     # detect cards, all 
    PICC_ANTICOLL   = 0x93     # anti collision
    PICC_SElECTTAG  = 0x93     # select card
    PICC_AUTHENT1A  = 0x60     # authent A
    PICC_AUTHENT1B  = 0x61     # authent B
    PICC_READ       = 0x30     # read block
    PICC_WRITE      = 0xA0     # write block
    PICC_DECREMENT  = 0xC0     # dec money
    PICC_INCREMENT  = 0xC1     # inc money
    PICC_RESTORE    = 0xC2     # send data to buf
    PICC_TRANSFER   = 0xB0     # save date from buf
    PICC_HALT       = 0x50     # idle status
 
    Reserved00      = 0x00     # register
    CommandReg      = 0x01
    CommIEnReg      = 0x02
    DivlEnReg       = 0x03
    CommIrqReg      = 0x04
    DivIrqReg       = 0x05
    ErrorReg        = 0x06
    Status1Reg      = 0x07
    Status2Reg      = 0x08
    FIFODataReg     = 0x09
    FIFOLevelReg    = 0x0A
    WaterLevelReg   = 0x0B
    ControlReg      = 0x0C
    BitFramingReg   = 0x0D
    CollReg         = 0x0E
    Reserved01      = 0x0F
 
    Reserved10      = 0x10
    ModeReg         = 0x11
    TxModeReg       = 0x12
    RxModeReg       = 0x13
    TxControlReg    = 0x14
    TxAutoReg       = 0x15
    TxSelReg        = 0x16
    RxSelReg        = 0x17
    RxThresholdReg  = 0x18
    DemodReg        = 0x19
    Reserved11      = 0x1A
    Reserved12      = 0x1B
    MifareReg       = 0x1C
    Reserved13      = 0x1D
    Reserved14      = 0x1E
    SerialSpeedReg  = 0x1F
 
    Reserved20        = 0x20
    CRCResultRegM     = 0x21
    CRCResultRegL     = 0x22
    Reserved21        = 0x23
    ModWidthReg       = 0x24
    Reserved22        = 0x25
    RFCfgReg          = 0x26
    GsNReg            = 0x27
    CWGsPReg          = 0x28
    ModGsPReg         = 0x29
    TModeReg          = 0x2A
    TPrescalerReg     = 0x2B
    TReloadRegH       = 0x2C
    TReloadRegL       = 0x2D
    TCounterValueRegH = 0x2E
    TCounterValueRegL = 0x2F
 
    Reserved30      = 0x30
    TestSel1Reg     = 0x31
    TestSel2Reg     = 0x32
    TestPinEnReg    = 0x33
    TestPinValueReg = 0x34
    TestBusReg      = 0x35
    AutoTestReg     = 0x36
    VersionReg      = 0x37
    AnalogTestReg   = 0x38
    TestDAC1Reg     = 0x39
    TestDAC2Reg     = 0x3A
    TestADCReg      = 0x3B
    Reserved31      = 0x3C
    Reserved32      = 0x3D
    Reserved33      = 0x3E
    Reserved34      = 0x3F
 
    NRSTPD          = 22
    MAX_LEN         = 18
 
    MI_OK           = 0
    MI_NOTAGERR     = 1
    MI_ERR          = 2
    MI_TIMEOUT      = 3
 
 
    def __init__(self, dev='/dev/spidev0.0', spd=1000000):
        self.dev0 = spi.openSPI(device=dev, speed=spd)
        GPIO.setmode(GPIO.BOARD)
        GPIO.setup(self.NRSTPD, GPIO.OUT)
        GPIO.output(self.NRSTPD, 1)
        self.MFRC522_Init()
 
    def __WriteReg(self, addr, val):
        spi.transfer(self.dev0, ((addr<<1)&0x7E, val))
 
    def __ReadReg(self, addr):
        print(addr, self.dev0)
        val = spi.transfer(self.dev0, (((addr<<1)&0x7E) | 0x80, 0))
        print(val)
        return val[1]
 
    def __SetRegBitMask(self, reg, mask):
        tmp = self.__ReadReg(reg)
        self.__WriteReg(reg, tmp | mask)
 
    def __ClearRegBitMask(self, reg, mask):
        tmp = self.__ReadReg(reg)
        self.__WriteReg(reg, tmp & (~mask))
 
    def __ToCard(self, command, sendData):
        retStatus = self.MI_OK
        backData = []
        backLen = 0
        irqEn = 0x00
        waitIRq = 0x00
        lastBits = None
        n = 0
        i = 0
 
        if command == self.PCD_AUTHENT:
            irqEn   = 0x12
            waitIRq = 0x10
        if command == self.PCD_TRANSCEIVE:
            irqEn   = 0x77
            waitIRq = 0x30
 
        self.__WriteReg(self.CommIEnReg, irqEn|0x80)      # enable interupt
        self.__ClearRegBitMask(self.CommIrqReg, 0x80)     # clear interupt flags
        self.__SetRegBitMask(self.FIFOLevelReg, 0x80)     # init FIFO
        self.__WriteReg(self.CommandReg, self.PCD_IDLE)   # cancel current command
 
        for i in range(len(sendData)):
            self.__WriteReg(self.FIFODataReg, sendData[i])
        self.__WriteReg(self.CommandReg, command)
 
        if command == self.PCD_TRANSCEIVE:
            self.__SetRegBitMask(self.BitFramingReg, 0x80)
 
        i = 2000
        while True:
            n = self.__ReadReg(self.CommIrqReg)
            i -= 1
            if ~((i!=0) and ~(n&0x01) and ~(n&waitIRq)):
                break
 
        self.__ClearRegBitMask(self.BitFramingReg, 0x80)
 
        if i != 0:
            if (self.__ReadReg(self.ErrorReg) & 0x1B) == 0x00:
                if n & irqEn & 0x01:
                    print("retStatus = self.MI_NOTAGERR")
                    retStatus = self.MI_NOTAGERR
 
                if command == self.PCD_TRANSCEIVE:
                    n = self.__ReadReg(self.FIFOLevelReg)
                    lastBits = self.__ReadReg(self.ControlReg) & 0x07
                    if lastBits != 0:
                        backLen = (n-1)*8 + lastBits
                    else:
                        backLen = n*8
 
                    if n == 0:
                        n = 1
                    if n > self.MAX_LEN:
                        n = self.MAX_LEN
 
                    for i in range(n):
                        backData.append(self.__ReadReg(self.FIFODataReg))
            else:
                print("retStatus = self.MI_ERR")
                retStatus = self.MI_ERR
        else:
            print("retStatus = self.MI_TIMEOUT")
            retStatus = self.MI_TIMEOUT
 
        return (retStatus, backData, backLen)
 
    def __CalulateCRC(self, indata):
        self.__ClearRegBitMask(self.DivIrqReg, 0x04)
        self.__SetRegBitMask(self.FIFOLevelReg, 0x80)
        for i in range(len(indata)):
            self.__WriteReg(self.FIFODataReg, indata[i])
        self.__WriteReg(self.CommandReg, self.PCD_CALCCRC)
 
        i = 255
        while True:
            n = self.__ReadReg(self.DivIrqReg)
            i -= 1
            if ~((i != 0) and ~(n&0x04)):
                break
 
        crc = []
        crc.append(self.__ReadReg(self.CRCResultRegL))
        crc.append(self.__ReadReg(self.CRCResultRegM))
        return crc
 
    def MFRC522_Init(self):
        self.MFRC522_Reset()
        self.__WriteReg(self.TModeReg, 0x8D)
        self.__WriteReg(self.TPrescalerReg, 0x3E)
        self.__WriteReg(self.TReloadRegL, 30)
        self.__WriteReg(self.TReloadRegH, 0)
        self.__WriteReg(self.TxAutoReg, 0x40)
        self.__WriteReg(self.ModeReg, 0x3D)
        self.MFRC522_AntennaOn()
 
    def MFRC522_Reset(self):
        # reg: 0x01
        # buf: 0x0F
        self.__WriteReg(self.CommandReg, self.PCD_RESETPHASE)
 
    def MFRC522_AntennaOn(self):
        # reg: 0x14
        # buf: 0bxxxxxx11
        temp = self.__ReadReg(self.TxControlReg)
        if not(temp & 0x03):
            self.__SetRegBitMask(self.TxControlReg, 0x03)
 
    def MFRC522_AntennaOff(self):
        # reg: 0x14
        # buf: 0bxxxxxx00
        self.__ClearRegBitMask(self.TxControlReg, 0x03)
 
    # function  : read block
    # parameter : blockAddr(0~63)
    # return    : retStatus
    #             backData[16]
    def MFRC522_ReadBolock(self, blockAddr):
        retStatus = self.MI_OK
 
        # cmd: 0x0c
        # buf: 0x30 blockAddr crc[2]
        buf = []
        buf.append(self.PICC_READ)
        buf.append(blockAddr)
        crc = self.__CalulateCRC(buf)
        buf.append(crc[0])
        buf.append(crc[1])
        (status, backData, backLen) = self.__ToCard(self.PCD_TRANSCEIVE, buf)
        if (status != self.MI_OK) or (backLen != 8*self.MAX_LEN):
            retStatus = self.MI_ERR
 
        return (retStatus, backData)
 
    # function  : write block
    # parameter : blockAddr(0~63)
    #             writeData[16]
    # return    : retStatus
    def MFRC522_WriteBlock(self, blockAddr, writeData):
        retStatus = self.MI_OK
 
        # cmd: 0x0c
        # buf: 0xA0 blockAddr crc[2]
        buf = []
        buf.append(self.PICC_WRITE)
        buf.append(blockAddr)
        crc = self.__CalulateCRC(buf)
        buf.append(crc[0])
        buf.append(crc[1])
        (status, backData, backLen) = self.__ToCard(self.PCD_TRANSCEIVE, buf)
        if (status != self.MI_OK) or ((backData[0]&0x0F) != 0x0A) or (backLen != 4):
            retStatus = self.MI_ERR
 
        if status == self.MI_OK:
            # cmd: 0x0c
            # buf: writeData[16] crc[2]
            buf2 = []
            for i in range(16):
                buf2.append(writeData[i])
            crc = self.__CalulateCRC(buf2)
            buf2.append(crc[0])
            buf2.append(crc[1])
            (status, backData, backLen) = self.__ToCard(self.PCD_TRANSCEIVE, buf2)
            if (status != self.MI_OK) or ((backData[0]&0x0F) != 0x0A) or (backLen != 4):
                retStatus = self.MI_ERR
 
        return retStatus
 
    # function  : detect card
    # parameter : reqMode: detect mode
    #               0x52 = detect all cards 
    #               0x26 = detect not sleep cards
    # return    : retStatus
    #             backData: card type(2 bytes)
    #               0x4400 = Mifare_UltraLight
    #               0x0400 = Mifare_One(S50)
    #               0x0200 = Mifare_One(S70)
    #               0x0800 = Mifare_Pro(X)
    #               0x4403 = Mifare_DESFire
    def MFRC522_Request(self, reqMode): 
        retStatus = self.MI_OK
 
        self.__WriteReg(self.BitFramingReg, 0x07)
        # cmd: 0x0c
        # buf: 0x26/0x52
        buf = []
        buf.append(reqMode)
        (status, backData, backLen) = self.__ToCard(self.PCD_TRANSCEIVE, buf)
        if (status != self.MI_OK) or (backLen != 0x10):
            retStatus = self.MI_ERR
 
        return (retStatus, backData)
 
    # function  : anticoll
    # parameter :
    # return    : retStatus
    #             backData(Uid, 4 bytes)
    def MFRC522_Anticoll(self):
        retStatus = self.MI_OK
        serNumCheck = 0
 
        self.__WriteReg(self.BitFramingReg, 0x00)
        # cmd: 0x0c
        # buf: 0x93 0x20
        buf = []
        buf.append(self.PICC_ANTICOLL)
        buf.append(0x20)
        (status, backData, backLen) = self.__ToCard(self.PCD_TRANSCEIVE, buf)
        if (status==self.MI_OK) and (len(backData)==5):
            for i in range(4):
                serNumCheck ^= backData[i]
            if serNumCheck != backData[4]:
                retStatus = self.MI_ERR
        else:
            retStatus = self.MI_ERR
 
        return (retStatus, backData)
 
    # function  : select card
    # parameter : Uid
    # return    : retStatus
    def MFRC522_Select(self, Uid):
        retStatus = self.MI_OK
        serNumCheck = 0
 
        # cmd: 0x0c
        # buf: 0x93 0x70 Uid[4] check crc[2]
        buf = []
        buf.append(self.PICC_SElECTTAG)
        buf.append(0x70)
        for i in range(4):
            buf.append(Uid[i])
            serNumCheck ^= Uid[i]
        buf.append(serNumCheck)
        crc = self.__CalulateCRC(buf)
        buf.append(crc[0])
        buf.append(crc[1])
        (status, backData, backLen) = self.__ToCard(self.PCD_TRANSCEIVE, buf)
        if (status != self.MI_OK) or (backLen != 0x18):
            retStatus = self.MI_ERR
 
        return retStatus
 
    # function  : auth sectorkey
    # parameter : authMode(PICC_AUTHENT1A/PICC_AUTHENT1B)
    #             blockAddr
    #             sectorkey(6 bytes)
    #             Uid(4 bytes)
    # return    : retStatus
    def MFRC522_Auth(self, authMode, blockAddr, sectorkey, Uid):
        retStatus = self.MI_OK
 
        # cmd: 0x0e
        # buf: authMode blockAddr sectorkey[6] Uid[4]
        buf = []
        buf.append(authMode)
        buf.append(blockAddr)
        for i in range(6):
            buf.append(sectorkey[i])
        for i in range(4):
            buf.append(Uid[i])
        (status, backData, backLen) = self.__ToCard(self.PCD_AUTHENT, buf)
        if (status != self.MI_OK) or not(self.__ReadReg(self.Status2Reg)&0x08):
            retStatus = self.MI_ERR
 
        return retStatus
 
    # function  : idle
    # parameter :
    # return    : retStatus
    def MFRC522_Halt(self):
        retStatus = self.MI_OK
 
        # cmd: 0x0c
        # buf: 0x50 0x00 crc[2]
        buf = []
        buf.append(self.PICC_HALT)
        buf.append(0)
        crc = self.__CalulateCRC(buf)
        buf.append(crc[0])
        buf.append(crc[1])
        (status, backData, backLen) = self.__ToCard(self.PCD_TRANSCEIVE, buf)
        retStatus = status
 
        return retStatus
 
    def MFRC522_WriteCmd40(self):
        retStatus = self.MI_OK
 
        self.__WriteReg(self.BitFramingReg, 0x07)
        # cmd: 0x0c
        # buf: 0x40
        buf = []
        buf.append(0x40)
        (status, backData, backLen) = self.__ToCard(self.PCD_TRANSCEIVE, buf)
        if (status != self.MI_OK) or (backData[0] != 0x0a):
            retStatus = status
 
        return retStatus
 
    def MFRC522_WriteCmd43(self):
        retStatus = self.MI_OK
 
        self.__WriteReg(self.BitFramingReg, 0x00)
        # cmd: 0x0c
        # buf: 0x43
        buf = []
        buf.append(0x43)
        (status, backData, backLen) = self.__ToCard(self.PCD_TRANSCEIVE, buf)
        if (status != self.MI_OK) or (backData[0] != 0x0a):
            retStatus = status
 
        return retStatus
 
    def MFRC522_StopCrypto1(self):
        self.__ClearRegBitMask(self.Status2Reg, 0x08)
 
    def MFRC522_CloseSPI(self):
        spi.closeSPI(self.dev0)
