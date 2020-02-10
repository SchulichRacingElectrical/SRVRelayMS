import socket
import struct
import ctypes
import sys
import time
import socketio
from threading import Thread

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
        self.readData(soc)

    def readData(self, sock):
        lastPacketID = -1
        while True:
            message, address = sock.recvfrom(4096)
            # Padding needed to read doubles (xxxx)
            fmt = "<Iffffffffffffffffffffxxxxddff"
            fmt_size = struct.calcsize(fmt)
            y = struct.unpack(fmt, message[:fmt_size])
            # Check Packet is more recent than old
            if lastPacketID < y[0]:
                lastPacketID = y[0]
                i = 0
                for value in vehicleData:
                    # Keeping Decimal Values For Longi and Lati
                    if i != 22 and i != 21:
                        vehicleData[value] = round(y[i], 2)
                        i = i + 1
                    else:
                        print(sys.getsizeof(y[i]))
                        vehicleData[value] = y[i]
                        i = i + 1
                print("{" + "\n".join("{!r}: {!r},".format(k, v)
                                      for k, v in vehicleData.items()) + "}")
                sio.emit('message', vehicleData)
            # Print PacketOutOfOrder Error and Ignore packet
            else:
                print('PacketOutOfOrder: Packet ' + str(y[0]) + ' dropped')

    def clientThread(self, connection, IP, PORT, MAX_BUFFER_SIZE=4096):
        b = 0
        while 1:
            data = connection.recv(MAX_BUFFER_SIZE)
            import time
            millis = int(round(time.time() * 1000))
            print(millis)
            if not data:
                break
            fmt = "<fffffffffffffffffffffffff"
            fmt_size = struct.calcsize(fmt)
            y = struct.unpack(fmt, data[:fmt_size])
            i = 0
            for value in vehicleData:
                vehicleData[value] = round(y[i], 2)
                i = i + 1
            print(y)
        connection.close()
