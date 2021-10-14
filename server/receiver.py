# Copyright Schulich Racing FSAE
# Written By Justin Tijunelis

import sys
import socket
import struct
import ctypes
import threading
import time

maptypes = {
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
      # TODO: Throw Exception
    packet_resetter = threading.Thread(target = self.reset_packet_tracker)
    packet_resetter.start()
    self.read_data(soc)

  def read_data(self, sock):
    self.last_packet_id = -1
    while True:
      self.last_packet_time = int(round(time.time() * 1000))
      # message, _ = sock.recvfrom(4096)
      message = str.encode(
        "6" +
        "abcdef" +
        "A\0\0\0" +          # 65
        "B\0\0\0" +          # Floating point (32-bit floating point)
        "\0\0\0\0\0\0\0\0" + # 0 (64-bit floating point)
        "\0\0\0\0\0\0\0\0" + # 0 (64-bit floating point)
        "E\0\0\0" +          # 69
        "\0\0\0\0\0\0\0\0"   # 0 (64-bit floating point)
      )
      data = self.parse_data(message)
  
  def parse_data(self, message):
    # Parse the message
    sensor_count = int(message.decode("utf-8")[0])
    sensor_ids = list(message.decode("utf-8")[1: sensor_count + 1])
    data_bytes = message[sensor_count + 1:]
    
    # Create the decode string based on sensor types
    # TODO: Update to read from sensor data
    # TODO: Figure out when padding is required
    # TODO: Add other types as needed
    data_format = "<"
    for i, sensor_id in enumerate(sensor_ids):
      sensor_type = maptypes[sensor_id]
      data_format += sensor_type

    # Decode the data
    data_format_size = struct.calcsize(data_format)
    data = struct.unpack(data_format, data_bytes)
    return data

  def reset_packet_tracker(self):
    while True:
      time.sleep(1)
      current_time = int(round(time.time() * 1000))
      if current_time - self.last_packet_time > 1000:
        if self.last_packet_id != -1:
          self.last_packet_id = -1
