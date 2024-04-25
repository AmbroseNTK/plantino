import paho.mqtt.client as mqtt
from paho.mqtt import publish
import random
import json

import time

SERIAL_NUMBER = "1938-3092-9478-1823"


def on_connect(client, userdata, flags, rc, properties):
    print("Connected with result code " + str(rc))
    client.subscribe("platino/" + SERIAL_NUMBER + "/command")


def on_message(client, userdata, msg):
    # process received message
    # validate message
    if data is None:
        return
    print("Received message: " + str(data))

    command = data["command"]
    payload = data["payload"]


client = mqtt.Client(mqtt.CallbackAPIVersion.VERSION2)
client.on_connect = on_connect
client.on_message = on_message

client.connect("broker.emqx.io", 1883, 60)

while True:
    client.loop()
    # read temperature random between 20 to 30 in decimal
    temp = random.uniform(20, 30)
    # read moisture
    moisture = random.uniform(40, 70)
    data = {"temperature": temp, "moisture": moisture}
    print("Publishing data: " + str(data))
    # publish data
    publish.single(
        "plantino/" + SERIAL_NUMBER + "/data",
        json.dumps(data),
        hostname="broker.emqx.io",
    )
    # sleep for 5 seconds
    time.sleep(1)
