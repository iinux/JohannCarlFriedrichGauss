import socket
import threading

import socks

debug = False


def p2p(socket1, socket2):
    while True:
        response = socket1.recv(1024)
        if len(response) == 0:
            print("socket1 %s is disconnected" % socket1)
            return
        socket2.sendall(response)
        print("socket1 %s receive %d %s and send to socket2 %s" % (socket1, len(response), response if debug else '',
                                                                   socket2))


def handle_connection(client_socket):
    """处理客户端连接"""

    request = client_socket.recv(1024)
    host = '2600:9000:23d1:2400:17:b174:6d00:93a1'
    port = 443

    # 创建 socks 代理
    proxy = socks.socksocket()
    proxy.set_proxy(socks.SOCKS5, "127.0.0.1", 7890)

    # 通过 socks 代理请求目标网站
    proxy.connect((host, int(port)))
    print("connect %s %s success" % (host, port))
    proxy.sendall(request)
    print("sendall %s %s success" % (host, port))

    t1 = threading.Thread(target=p2p, args=(proxy, client_socket))
    t1.setDaemon(True)
    t1.start()

    t2 = threading.Thread(target=p2p, args=(client_socket, proxy))
    t2.setDaemon(True)
    t2.start()


def main():
    """主函数"""

    # 创建监听 socket
    listen_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    listen_socket.bind(("0.0.0.0", 443))
    listen_socket.listen(1)

    try:
        while True:
            # 等待客户端连接
            client_socket, _ = listen_socket.accept()

            # 处理客户端连接
            t = threading.Thread(target=handle_connection, args=(client_socket,))
            t.setDaemon(True)
            t.start()
    except:
        print("close listen socket")
        listen_socket.close()


if __name__ == "__main__":
    main()
