# Copyright Schulich Racing FSAE
# Written By Justin Tijunelis

import sys
import socket
import struct
import ctypes
import threading
import time

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
      sys.exit() # Throw exception or inform somehow
    packet_resetter = threading.Thread(target = self.reset_packet_tracker)
    packet_resetter.start()
    self.read_data(soc)

  def read_data(self, sock):
    self.last_packet_id = -1
    while True:
      self.last_packet_time = int(round(time.time() * 1000))
      # Read data from socket and create sender
      # TODO: Have time expiry that closes the socket and resets? 10 seconds
    return

  def decode_data(self, data):
    
    pass

  def reset_packet_tracker(self):
    while True:
      time.sleep(1)
      current_time = int(round(time.time() * 1000))
      if current_time - self.last_packet_time > 1000:
        if self.last_packet_id != -1:
          self.last_packet_id = -1
