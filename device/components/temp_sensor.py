from components.rs485_device import Rs485Device
import time


class TemperatureSensor(Rs485Device):
    def __init__(self, pin):
        self.pin = pin

    def read(self):
        self.send([3, 0, 6, 0, 1, 100, 11])
        time.sleep(0.5)
        return self.read()
