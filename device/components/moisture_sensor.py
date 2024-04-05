from components.rs485_device import Rs485Device

import time


class MoistureSensor(Rs485Device):
    def __init__(self, pin):
        self.pin = pin

    def read(self):
        self.send([3, 0, 7, 0, 1, 53, 203])
        time.sleep(0.5)
        return self.read()
