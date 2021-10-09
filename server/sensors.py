# Copyright Schulich Racing FSAE
# Written By Justin Tijunelis, Camilla Abdrazakov, Jonathan Mulyk

class Sensors:
  def __init__(self):
    # Get the sensor list from the database
    self.version = 1.0
    self.sensors = []

  def fetch_from_server(self):
    # Get the sensor list
    pass

  def get_version(self):
    # Return the version
    return self.version

  def get_diff(old_version):
    # Get difference between old and new version and return
    pass

