import socket
import struct
import ctypes
from threading import Thread


class MyStruct:
    _fields_ = [
        ("x", ctypes.c_double),
        ("y", ctypes.c_double),
        ("z", ctypes.c_float)
    ]


class Network:
    listenAddress = "0.0.0.0"
    listenPort = 8080

    def startServer(self):
        soc = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        # this is for easy starting/killing the app
        soc.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
        print('Socket created')

        try:
            soc.bind((self.listenAddress, self.listenPort))
            print('Socket bind complete')
        except socket.error as msg:
            import sys
            print('Bind failed. Error : ' + str(sys.exc_info()))
            print(msg)
            sys.exit()

        # Start listening on socket
        soc.listen(10)
        print('Socket now listening')

        # this will make an infinite loop needed for
        # not reseting server for every client
        while True:
            connection, addr = soc.accept()
            ip, port = str(addr[0]), str(addr[1])
            print('Accepting connection from ' + ip + ':' + port)
            try:
                Thread(target=self.clientThread, args=(
                    connection, ip, port)).start()
            except:
                print("Terrible error!")
                import traceback
                traceback.print_exc()
        soc.close()

    def clientThread(self, connection, IP, PORT, MAX_BUFFER_SIZE=4096):
        while 1:
            data = connection.recv(MAX_BUFFER_SIZE)
            if not data:
                break
            print("received data:", data)
            x = MyStruct()
            fmt = "<ddf"
            fmt_size = struct.calcsize(fmt)
            x.x, x.y, x.z = struct.unpack(fmt, data[:fmt_size])
            print("Fields\n  x: {:f}\n  y: {:f}\n  z: {:f}".format(
                x.x, x.y, x.z))
        connection.close()
