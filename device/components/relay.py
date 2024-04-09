import time
from components.rs485_device import Rs485Device


class Relay(Rs485Device):
    def __init__(self, serial, address):
        super().__init__(serial, address)

    def turn_on(self):
        super.send([6, 0, 0, 0, 255, 200, 91])

    def turn_off(self):
        super.send([6, 0, 0, 0, 0, 136, 27])

    def read_status(self):
        self.send([])
        time.sleep(0.5)
        return self.read()
