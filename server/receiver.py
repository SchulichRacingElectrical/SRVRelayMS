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
  'c': 'd'
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
      # Read data from socket and create sender
      message, _ = sock.recvfrom(4096)
      message = "3abcAAAABC"
      print(message)
      num_sensor = int(message[0])
      id_sensor = list(message[1:num_sensor+1])
      fmt = list(message[num_sensor+1:])
      msg = tuple(zip(id_sensor,fmt))
      print(msg)
      
      for i in range(id_sensor):
        if id_sensor[i - 1] is not 'c' and i is not 0:
          pass
          # insert 4 x
          # check a
          # check b
        # else replace c with d
  
  def parse_data(self, data):
    pass

  def reset_packet_tracker(self):
    while True:
      time.sleep(1)
      current_time = int(round(time.time() * 1000))
      if current_time - self.last_packet_time > 1000:
        if self.last_packet_id != -1:
          self.last_packet_id = -1
