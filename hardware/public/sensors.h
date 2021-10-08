// Copyright Schulich Racing FSAE
// Written by Justin Tijunelis

#pragma once

struct Sensor {
  unsigned char *name;
  unsigned char id;
  unsigned short can_id; // Is it a short or int? Will be interpreted as hex anyways.
  unsigned char frequency; // Frequency will never be over 255, probably never over 60. 
};

// Keep sensors in a linked list on flash?
// Could alternatively use an array but that seems shitty. 
// Need to think of what data type to use to store all of these values.
// Should probably dynamically allocate memory for sensor information.
// Do we want to organize sensors categorically?
struct Sensors {
  unsigned char *version;
  struct Sensor *sensor;
};

// Could use static?
struct Sensors *sensor_data = 0;

/*
*
*/
void add_sensor(struct Sensor sensor);

/*
*
*/
void remove_sensor(unsigned char id);

/*
*
*/
void update_sensor(struct Sensor sensor);
