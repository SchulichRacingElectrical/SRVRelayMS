import socketio

class Relay:
  sio = socketio.Client()

  def __init__(self, ip = "127.0.0.1", port = 5000, sensors = None):
    self.ip = ip
    self.port = port
    self.sensors = sensors
    address = "http://" + self.ip + ":" + str(self.port)
    try:
      self.sio.connect(address)
    except:
      pass

  def send_data(self, data):
    # TODO: Set other sensor data to be the most recent?
    # When passing data, send the value of every existing sensor
    pass