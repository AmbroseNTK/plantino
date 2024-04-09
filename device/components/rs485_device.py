import serial as se


class Rs485Device:
    def __init__(self, serial, address):
        self.serial = serial
        self.address = address

    def send(self, data=[]):
        if len(data) != 7:
            raise ValueError("Data must be 7 bytes long")
        data = [self.address] + data
        print(data)
        self.serial.write(data)

    def read(self):
        bytesToRead = self.serial.inWaiting()
        if bytesToRead == 0:
            return 0
        data = self.serial.read(bytesToRead)
        print(data)
        dataArray = list(data)
        print(dataArray)
        if len(dataArray) >= 7:
            array_size = len(dataArray)
            value = dataArray[array_size - 4] * 256 + dataArray[array_size - 3]
            return value
        else:
            return -1
