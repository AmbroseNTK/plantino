# using paho mqtt client

import paho.mqtt.client as mqtt
from paho.mqtt import publish
from components.relay import Relay
from components.moisture_sensor import MoistureSensor
from components.temp_sensor import TemperatureSensor
import crypto as cr
import serial
import time
import json

PORT_NAME = "/dev/ttyUSB0"
SERIAL_NUMBER = "1938-3092-9478-1823"


def connect_rs485():
    print("Connecting to RS485 port")
    try:
        ser = serial.Serial(PORT_NAME, 9600, timeout=1)
        print("Connected to RS485 port")
        return ser
    except:
        print("Failed to connect to RS485 port")
        return None


def init_devices(serial_conn):
    relays = {
        "2": Relay(serial_conn, 2),
        "3": Relay(serial_conn, 3),
        "4": Relay(serial_conn, 4),
    }

    temp_sensor = TemperatureSensor(serial_conn, 1)
    moisture_sensor = MoistureSensor(serial_conn, 1)

    return relays, temp_sensor, moisture_sensor


relays, temp_sensor, moisture_sensor = init_devices(connect_rs485())


def on_connect(client, userdata, flags, rc):
    print("Connected with result code " + str(rc))
    client.subscribe("platino/" + SERIAL_NUMBER + "/command")


def on_message(client, userdata, msg):
    # process received message
    # validate message
    data = cr.verify_token(msg.payload)
    if data is None:
        return
    print("Received message: " + str(data))

    command = data["command"]
    payload = data["payload"]

    if command == "relay":
        relay_id = payload["relay_id"]
        action = payload["action"]
        if relay_id not in relays:
            print("Invalid relay id")
            return
        if action == "on":
            relays[relay_id].turn_on()
        elif action == "off":
            relays[relay_id].turn_off()
        else:
            print("Invalid action")


client = mqtt.Client()
client.on_connect = on_connect
client.on_message = on_message

client.connect("broker.emqx.io", 1883, 60)

while True:
    client.loop()
    # read temperature
    temp = temp_sensor.read_sensor()
    print("Temperature: " + str(temp))
    # read moisture
    moisture = moisture_sensor.read_sensor()
    print("Moisture: " + str(moisture))
    data = {"temperature": temp, "moisture": moisture}
    # publish data
    publish.single(
        "plantino/" + SERIAL_NUMBER + "/data",
        json.dumps(data),
        hostname="broker.emqx.io",
    )
    # sleep for 5 seconds
    time.sleep(1)
