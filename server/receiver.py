# Copyright Schulich Racing FSAE
# Written By Justin Tijunelis, Camilla Abdrazakov, Jonathan Mulyk

import sys
import socket
import struct
import ctypes
import threading
import time

maptypes = {
  'I': 'Q', # Always included, the interval or UTC of the data (defines ID)
  'a': 'I',
  'b': 'f',
  'c': 'd',
  'd': 'd',
  'e': 'I',
  'f': 'd'
}

class Receiver:
  def __init__(self, sensors, relay):
    self.sensors = sensors
    self.relay = relay
    self.last_packet_time = -1

  def start_receiver(self, port):
    soc = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    try:
      soc.bind(('', port))
    except socket.error as msg:
      print("Bind failed. Error: " + str(sys.exc_info()))
      # TODO: Handle exception
    packet_resetter = threading.Thread(target = self.__handle_disconnect)
    packet_resetter.start()
    self.read_data(soc)

  def read_data(self, sock):
    self.last_packet_id = -1 # Equivalent to the UTC (should we use a count instead?)
    while True:
      # Wait for a new message
      message, _ = sock.recvfrom(4096)
      self.last_packet_time = int(round(time.time() * 1000))

      # Parse the message
      sensor_count = message[0]
      sensor_ids = list(message.decode()[1:sensor_count + 1])
      
      # Create the decode string based on sensor types
      # TODO: Update to read from sensor data rather than maptypes
      # TODO: Add other types as needed (Padding required for variables < 4 bytes)
      data_format = "<"
      for i, sensor_id in enumerate(sensor_ids):
        data_format += maptypes[sensor_id]

      # Decode the data
      data = struct.unpack(data_format, message[sensor_count + 1 :])
      
      # Handle the data
      if self.last_packet_id < data[0]:
        self.last_packet_id = data[0]
        self.sensors.update_values(data)
        self.relay.send_data()
      else:
        # TODO: The data came out of order, but we still want to use it
        pass

  def __handle_disconnect(self):
    while True:
      time.sleep(1)
      current_time = int(round(time.time() * 1000))
      if current_time - self.last_packet_time > 1000:
        if self.last_packet_id != -1:
          self.last_packet_id = -1
          # This means there was a disconnection, we should create and store the data
