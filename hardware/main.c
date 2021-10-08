// Copyright Schulich Racing FSAE
// Written by Justin Tijunelis

#include <stdio.h>
#include <unistd.h>
#include "public/reader.h"
#include "public/receiver.h"
#include "public/sender.h"
#include "public/sensors.h"

// Define these in their respective locations?
#define SERVER_TCP_ENDPOINT "https://127.0.0.1:5000" // Change port in the future
#define SERVER_UDP_ENDPOINT "https://127.0.0.1:4000" // Change port in the future
#define API_KEY "xxx"                                // Change in the future

int main(void) {
  // Configuration

  // Make the request to the server asking for the sensor version list
  // Update the sensor information in flash
  consolidate_sensors();
  // TODO: Other configuration?

  // Run Forever
  unsigned char max_frequency = 60; // Hz, get from sensor list in the future
  unsigned char timestamp = 0;
  while (1) {
    struct Sample *sample = read_CAN_data(timestamp);
    // Read sensor data over CAN (simulate for now)
    // Send the data as an array of bytes over UDP
    // Dynamically allocate memory for sensor data?
    // Need to keep last value of sensor data in memory to check if it has changed.
    sleep((float)(1 / (unsigned int)max_frequency));
    timestamp = ++timestamp % max_frequency;
  }
  return 0;
}