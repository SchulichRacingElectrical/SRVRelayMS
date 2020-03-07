import socket
import struct
import ctypes
import sys
import time
import socketio
import threading

sio = socketio.Client()
sio.connect('http://127.0.0.1:4000')

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

    def __init__(self):
        self.lastPacketTime = 0
        self.lastPacketTime = -1

    def startServer(self):
        soc = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
        try:
            soc.bind(('', self.listenPort))
        except socket.error as msg:
            print('Bind failed. Error : ' + str(sys.exc_info()))
            print(msg)
            sys.exit()
        # Start listening on sock
        print('Socket now listening')
        packetResetter = threading.Thread(target=self.resetPacketTracker)
        packetResetter.start()
        self.readData(soc)

    def readData(self, sock):
        self.lastPacketID = -1
        while True:
            self.lastPacketTime = int(round(time.time() * 1000))
            message, address = sock.recvfrom(4096)
            # Padding needed to read doubles (xxxx)
            fmt = "<Iffffffffffffffffffffxxxxddff"
            fmt_size = struct.calcsize(fmt)
            y = struct.unpack(fmt, message[:fmt_size])
            # Check Packet is more recent than old
            if self.lastPacketID < y[0]:
                self.lastPacketID = y[0]
                i = 0
                for value in vehicleData:
                    # Keeping Decimal Values For Longi and Lati
                    if i != 22 and i != 21:
                        vehicleData[value] = round(y[i], 2)
                        i = i + 1
                    else:
                        vehicleData[value] = y[i]
                        i = i + 1
                print("{" + "\n".join("{!r}: {!r},".format(k, v)
                                      for k, v in vehicleData.items()) + "}")
                sio.emit('message', vehicleData)
            # Print PacketOutOfOrder Error and Ignore packet
            else:
                print('PacketOutOfOrder: Packet ' + str(y[0]) + ' dropped')

    def resetPacketTracker(self):
        print("Reset Started")
        while True:
            time.sleep(3)
            if int(round(time.time() * 1000)) - self.lastPacketTime > 3000 and self.lastPacketID != -1:
                print("Reset Packet Tracker")
                self.lastPacketID = -1
