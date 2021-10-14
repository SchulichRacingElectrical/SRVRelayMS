from main import app

class Sender:
  def __init__(self, sensors):
    self.sensors = sensors
    # Make request to get sensors
    pass

  @app.route("/sensors/get_version")
  def get_version(self):
    # Return a version
    pass

  @app.route("/sensors/get_sensors")
  def get_sensors(self):
    # Return all sensors and send to the hardware
    pass

  @app.route("/sensors/get_diff")
  def get_diff(self, version):
    # Send the difference between version
    pass

  def send_message(self, message):
    # Send message to the car that will be on the display
    pass