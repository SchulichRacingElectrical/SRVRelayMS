#!/bin/bash

# Compile and run the hardware simulation

gcc main.c private/reader.c private/receiver.c private/sender.c private/sensors.c -o main -Wall

./main