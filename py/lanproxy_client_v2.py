import argparse
import queue
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
            zero = 0
            bs += zero.to_bytes(1, 'big')

        if data is not None:
            if isinstance(data, bytes):
                bs += data
            else:
                bs += data.encode()

        return bs


class ClientThread(threading.Thread):
    def __init__(self, proxy_socket, client_socket, uri):
        super().__init__()
        self.proxy_socket = proxy_socket
        self.client_socket = client_socket
        self.uri = uri

    def run(self):
        while True:
            if getattr(self.client_socket, '_closed'):
                print('socket is close')
                break
            try:
                data = self.client_socket.recv(1024)
                data_len = len(data)
                # print('uri %s socket %s recv data len %d' % (self.uri, self.client_socket.getsockname(), data_len))
                if data_len > 0:
                    self.proxy_socket.sendall(pm.encode(ProxyMessage.P_TYPE_TRANSFER, 0, self.uri, data))
                else:
                    break
            except socket.timeout:
                print('socket ins timeout')
                break

        self.client_socket.close()
        del client_socket_map[self.uri]
        psp.put_idle(self.proxy_socket)


class HeartBeat(threading.Thread):
    def __init__(self, heart_beat_socket):
        super().__init__()
        self.heart_beat_socket = heart_beat_socket

    def run(self):
        while True:
            time.sleep(10)
            self.heart_beat_socket.sendall(pm.encode(ProxyMessage.TYPE_HEARTBEAT, 0, None, None))


class ProxySocketPool:
    def __init__(self):
        self.idle_queue = queue.Queue()

    def get_socket(self):
        if self.idle_queue.empty():
            proxy_socket = socket.socket()
            proxy_socket.connect(proxy_server_ip_port)
            return proxy_socket, True
        else:
            return self.idle_queue.get(), False

    def put_idle(self, idle_proxy_socket):
        self.idle_queue.put(idle_proxy_socket)


class Handler(threading.Thread):
    def __init__(self, proxy_socket):
        super().__init__()
        self.proxy_socket = proxy_socket

        heat_beat_thread = HeartBeat(self.proxy_socket)
        heat_beat_thread.setDaemon(True)
        heat_beat_thread.start()

    def recv(self, recv_len):
        rbs = bytes()
        while len(rbs) < recv_len:
            rbs_one = self.proxy_socket.recv(recv_len - len(rbs))
            if len(rbs_one) < 1:
                print('recv empty break')
                self.proxy_socket.close()
                break

            rbs += rbs_one

        if len(rbs) < recv_len:
            raise Exception('want read %d bus read %d' % (recv_len, len(rbs)))
        return rbs

    def run(self):
        while True:
            msg_len_str = self.recv(4)
            msg_len = int.from_bytes(msg_len_str, 'big')
            # print('msg_len is ' + str(msg_len))

            msg = self.recv(msg_len)

            pack = struct.unpack('!bqb', msg[0:10])
            msg_type = pack[0]
            sn = pack[1]
            uri_len = pack[2]

            uri = msg[10: 10 + uri_len]
            data = msg[10 + uri_len:]

            # print('msg_type is ' + str(msg_type))

            if msg_type == ProxyMessage.TYPE_CONNECT:
                connect_ip_port = data.decode()
                print('type connect ' + connect_ip_port)

                ds = connect_ip_port.split(':')
                ip_port = (ds[0], int(ds[1]))
                client_socket = socket.socket()
                client_socket.connect(ip_port)

                proxy_socket = psp.get_socket()

                client_thread = ClientThread(proxy_socket[0], client_socket, uri)
                client_thread.setDaemon(True)
                client_thread.start()

                client_socket_map[uri] = client_socket

                proxy_socket[0].sendall(pm.encode(ProxyMessage.TYPE_CONNECT, 0, uri.decode() + '@' + key, None))

                if proxy_socket[1]:
                    handler_thread = Handler(proxy_socket[0])
                    handler_thread.setDaemon(True)
                    handler_thread.start()
            elif msg_type == ProxyMessage.P_TYPE_TRANSFER:
                # print('type transfer')

                if client_socket_map[uri] is None:
                    print('get client socket by uri is none')
                else:
                    client_socket = client_socket_map[uri]
                    client_socket.sendall(data)
                    # print('uri %s socket %s send data len %d' % (uri, client_socket.getsockname(), len(data)))
            elif msg_type == ProxyMessage.TYPE_HEARTBEAT:
                # print('type heart beat')
                pass

            time.sleep(0.1)


parser = argparse.ArgumentParser(description='lanproxy client')
parser.add_argument('-s', '--server', help='proxy server host', default='127.0.0.1')
parser.add_argument('-p', '--port', help='proxy server port', default=4900)
parser.add_argument('-k', '--key', help='client key', default='4ec72fba10664f748cf4bfda005cefc6')
args = parser.parse_args()

proxy_server_ip_port = (args.s, args.p)
key = args.k

core_socket = socket.socket()
core_socket.connect(proxy_server_ip_port)

pm = ProxyMessage()
core_socket.sendall(pm.encode(ProxyMessage.C_TYPE_AUTH, 0, key, None))

client_socket_map = {}
psp = ProxySocketPool()

try:
    print('started')
    Handler(core_socket).run()
except KeyboardInterrupt as e:
    print('KeyboardInterrupt , exit')
except Exception as e:
    print('Exception ' + str(e) + ' exit')

print('close core socket')
core_socket.close()
