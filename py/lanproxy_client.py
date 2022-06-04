import socket
import struct
import threading
import time


class ProxyMessage:
    # /** 心跳消息 */
    TYPE_HEARTBEAT = 0x07

    # /** 认证消息，检测clientKey是否正确 */
    C_TYPE_AUTH = 0x01

    # // /** 保活确认消息 */
    # // public static final byte TYPE_ACK = 0x02;

    # /** 代理后端服务器建立连接消息 */
    TYPE_CONNECT = 0x03

    # /** 代理后端服务器断开连接消息 */
    TYPE_DISCONNECT = 0x04

    # /** 代理数据传输 */
    P_TYPE_TRANSFER = 0x05

    # /** 用户与代理服务器以及代理客户端与真实服务器连接是否可写状态同步 */
    C_TYPE_WRITE_CONTROL = 0x06

    # /** 消息类型 */
    # byte type;

    # /** 消息流水号 */ long
    # long serialNumber;

    # /** 消息命令请求信息 */ String
    # String uri;

    # /** 消息传输数据 */ byte[]
    # byte[] data;

    TYPE_SIZE = 1
    SERIAL_NUMBER_SIZE = 8
    URI_LENGTH_SIZE = 1

    def encode(self, msg_type, serial_number, uri, data):
        body_len = self.TYPE_SIZE + self.SERIAL_NUMBER_SIZE + self.URI_LENGTH_SIZE

        if uri is not None:
            body_len += len(uri)

        if data is not None:
            body_len += len(data)

        bs = body_len.to_bytes(4, 'big')
        bs += msg_type.to_bytes(1, 'big')
        bs += serial_number.to_bytes(8, 'big')

        if uri is not None:
            bs += len(uri).to_bytes(1, 'big')
            if isinstance(uri, bytes):
                bs += uri
            else:
                bs += uri.encode()
        else:
            bs += 0

        if data is not None:
            if isinstance(data, bytes):
                bs += data
            else:
                bs += data.encode()

        return bs


class ClientThread(threading.Thread):
    def __init__(self, socket_ins, uri):
        super().__init__()
        self.socket_ins = socket_ins
        self.uri = uri

    def run(self):
        while True:
            if getattr(self.socket_ins, '_closed'):
                print('socket is close')
                break

            data = self.socket_ins.recv(1024)
            data_len = len(data)
            print("socket ins recv len " + str(data_len))
            if data_len > 0:
                core_socket.sendall(pm.encode(ProxyMessage.P_TYPE_TRANSFER, 0, self.uri, data))


ip_port = ('127.0.0.1', 4900)
key = '4ec72fba10664f748cf4bfda005cefc6'

core_socket = socket.socket()  # 创建套接字
core_socket.connect(ip_port)  # 连接服务器

pm = ProxyMessage()
core_socket.sendall(pm.encode(ProxyMessage.C_TYPE_AUTH, 0, key, None))

socket_map = {}

while True:
    msg_len_str = core_socket.recv(4)
    msg_len = int.from_bytes(msg_len_str, 'big')
    print("msg_len is " + str(msg_len))

    msg = core_socket.recv(msg_len)

    pack = struct.unpack('!bqb', msg[0:10])
    msg_type = pack[0]
    sn = pack[1]
    uri_len = pack[2]

    uri = msg[10: 10 + uri_len]
    data = msg[10 + uri_len:]

    print(msg_type)

    if msg_type == ProxyMessage.TYPE_CONNECT:
        print("type connect")

        ds = data.decode().split(':')
        ip_port = (ds[0], int(ds[1]))
        client_socket = socket.socket()
        client_socket.connect(ip_port)
        ClientThread(client_socket, uri).start()
        socket_map[uri] = client_socket
        core_socket.sendall(pm.encode(ProxyMessage.TYPE_CONNECT, 0, uri.decode() + '@' + key, None))
    elif msg_type == ProxyMessage.P_TYPE_TRANSFER:
        print("type transfer")

        if socket_map[uri] is None:
            print('uri is none')
        else:
            socket_map[uri].sendall(data)

    time.sleep(0.1)

# s.close()
