# Copyright Schulich Racing FSAE
# Written By Justin Tijunelis

from flask import Flask
from receiver import Receiver
from sensors import Sensors
# from sender import Sender
from relay import Relay

app = Flask(__name__)

if __name__ == "__main__":
  # Initialize sensor information
  sensors = Sensors()
  sensors.fetch_from_server()

  # Create the relay server
  relay = Relay(sensors = sensors)

  # Start the receiver
  recv = Receiver(sensors, relay)
  recv.start_receiver(4500)

  # Start the sender
  # sender = Sender(sensors)

  # Start the TCP Server (Sender)
  app.run()