# Copyright Schulich Racing FSAE
# Written By Justin Tijunelis

from flask import Flask
from receiver import Receiver
app = Flask(__name__)

if __name__ == "__main__":
  app.run()
  recv = Receiver()
  recv.start_receiver()