from bluepy.btle import Scanner, DefaultDelegate
from multiprocessing import Process, Event


class ScanDelegate(DefaultDelegate):
    def __init__(self):
        DefaultDelegate.__init__(self)

    def handleDiscovery(self, dev, isNewDev, isNewData):
        if isNewDev:
            print("Discovered device", dev.addr)
        elif isNewData:
            print("Received new data from", dev.addr)

class BLEScanner:
    def __init__(self):
        self.scanner = Scanner().withDelegate(ScanDelegate())

        self.stop_event = Event()

    def start(self):

        self.stop_event.clear()

        self.process = Process(target=self.scan, args = ())
        self.process.start()
        return self

    def scan(self):
        while True:
            if self.stop_event.is_set():
                return
            self.devices = self.scanner.scan(5, passive=True)

    def stop(self):
        self.stop_event.set()
BLEScanner().start()
