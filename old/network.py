import socket
import struct
import ctypes
import sys
import time
import socketio
import threading

sio = socketio.Client()
sio.connect('http://127.0.0.1:5500')

vehicleData = {
    "count": 0,
    "rlSuspension": 0,
    "rrSuspension": 0,
    "flSuspension": 0,
    "frSuspension": 0,
    "tp": 0,
    "ipw": 0,
    "baro": 0,
    "map": 0,
    "atf": 0,  # AFR
    "iat": 0,
    "engineTemp": 0,
    "oilPres": 0,
    "oilTemp": 0,
    "fuelTemp": 0,
    "x": 0,
    "y": 0,
    "z": 0,
    "roll": 0,
    "pitch": 0,
    "yaw": 0,
    "longitude": 0,
    "latitude": 0,
    "speed": 0,
    "distance": 0,
    "egt1": 0,
    "egt2": 0,
    "egt3": 0,
    "egt4": 0,
    "o2": 0,
    "cam": 0,
    "crank": 0,
    "neutral": 0,
    "flSpeed": 0,
    "frSpeed": 0,
    "rlSpeed": 0,
    "rrSpeed": 0,
    "fbrakes": 0,
    "rbrakes": 0,
    "rotPot": 0,
    "voltage": 0,
    "rpm": 0
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
            message, _ = sock.recvfrom(4096)
            # Padding needed to read doubles (xxxx)
            fmt = "<Iffffffffffffffffffffxxxxddfffffffffffffffffff"
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
                sio.emit('message', vehicleData)

    def resetPacketTracker(self):
        while True:
            time.sleep(1)
            if int(round(time.time() * 1000)) - self.lastPacketTime > 1000 and self.lastPacketID != -1:
                print("Reset Packet Tracker ")
                self.lastPacketID = -1
