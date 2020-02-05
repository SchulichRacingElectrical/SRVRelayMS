import socket
import struct
import ctypes
import time
import socketio
from threading import Thread

# sio = socketio.Client()
# sio.connect('http://127.0.0.1:4000')

vehicleData = {
    "count": 0,
    "rearLeft": 0,
    "rearRight": 0,
    "frontLeft": 0,
    "frontRight": 0,
    "TPS": 0,
    "IPW": 0,
    "baro": 0,
    "MAP": 0,
    "AFR": 0,
    "IAT": 0,
    "engineTemp": 0,
    "oilPressure": 0,
    "oilTemp": 0,
    "fuelTemp": 0,
    "xAccel": 0,
    "yAccel": 0,
    "zAccel": 0,
    "roll": 0,
    "pitch": 0,
    "yaw": 0,
    "longitude": 0,
    "latitude": 0,
    "speed": 0,
    "distance": 0
}

class Network:
    listenAddress = "0.0.0.0"
    listenPort = 8000

    def startServer(self):
        soc = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
        try:
            soc.bind(('', self.listenPort))
        except socket.error as msg:
            import sys
            print('Bind failed. Error : ' + str(sys.exc_info()))
            print(msg)
            sys.exit()
        # Start listening on sock
        print('Socket now listening')
        while True:
            message, address = soc.recvfrom(4096)
            fmt = "<fffffffffffffffffffffffff"
            fmt_size = struct.calcsize(fmt)
            y = struct.unpack(fmt, message[:fmt_size])
            i = 0
            for value in vehicleData:
                vehicleData[value] = round(y[i],2)
                i = i + 1
            print(y)

    def clientThread(self, connection, IP, PORT, MAX_BUFFER_SIZE=4096):
        b = 0
        while 1:
            data = connection.recv(MAX_BUFFER_SIZE)
            import time
            millis = int(round(time.time() * 1000))
            print (millis)
            if not data:
                break
            fmt = "<fffffffffffffffffffffffff"
            fmt_size = struct.calcsize(fmt)
            y = struct.unpack(fmt, data[:fmt_size])
            i = 0
            for value in vehicleData:
                vehicleData[value] = round(y[i],2)
                i = i + 1
            print(y)
            sio.emit('message', vehicleData)
        connection.close()
