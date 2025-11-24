import random
import socket
import sys
import threading
import select


class TCPForwarder:
    def __init__(self, local_port, remote_host, remote_port):
        self.local_port = local_port
        self.remote_host = remote_host
        self.remote_port = remote_port
        self.remote_ips = []
        self.active_connections = {}
        self.lock = threading.Lock()

    def resolve_hostname(self):
        """解析域名获取所有IP地址"""
        try:
            import socket
            addr_info = socket.getaddrinfo(self.remote_host, None)
            ips = list(set([ip[4][0] for ip in addr_info]))  # 去重
            print(f"解析到的IP地址: {ips}")
            return ips
        except Exception as e:
            print(f"域名解析失败: {e}")
            return []

    def check_connectivity(self, ip, port, timeout=3):
        """检查TCP连接是否可达"""
        try:
            sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
            sock.settimeout(timeout)
            result = sock.connect_ex((ip, port))
            sock.close()
            return result == 0
        except Exception:
            return False

    def filter_reachable_ips(self, ips):
        """过滤出可连接的IP地址"""
        reachable_ips = []
        for ip in ips:
            if self.check_connectivity(ip, self.remote_port):
                reachable_ips.append(ip)
                print(f"IP {ip} 可达")
            else:
                print(f"IP {ip} 不可达，已排除")
        return reachable_ips

    def handle_local_connection(self, local_socket):
        """处理本地连接"""
        remote_sock = None
        try:
            ip = random.choice(self.remote_ips)
            try:
                remote_sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
                remote_sock.settimeout(5)
                remote_sock.connect((ip, self.remote_port))
            except Exception as e:
                print(f"连接到 {ip}:{self.remote_port} 失败: {e}")
                return

                # 使用select进行多路复用
            while True:
                try:
                    ready, _, _ = select.select([local_socket, remote_sock], [], [], 1)
                    if not ready:
                        continue

                    for sock in ready:
                        try:
                            data = sock.recv(4096)
                            if not data:
                                raise ConnectionResetError("连接已断开")

                            # 将数据转发到其他socket
                            if sock is local_socket:
                                try:
                                    remote_sock.send(data)
                                except Exception as e:
                                    print(f"发送数据到远程服务器失败: {e}")
                            else:
                                # 从远程发往本地
                                try:
                                    local_socket.send(data)
                                except Exception as e:
                                    print(f"发送数据到本地客户端失败: {e}")
                        except ConnectionResetError:
                            raise
                        except Exception as e:
                            print(f"接收数据时发生错误: {e}")
                            raise
                except Exception as e:
                    print(f"转发过程中发生错误: {e}")
                    break

        except Exception as e:
            print(f"连接处理过程中发生错误: {e}")
        finally:
            local_socket.close()
            try:
                if remote_sock is not None:
                    remote_sock.close()
            except:
                pass

    def start_server(self):
        """启动本地服务器"""
        # 解析域名并筛选可连接的IP
        print("正在解析域名...")
        all_ips = self.resolve_hostname()
        if not all_ips:
            print("未能解析到任何IP地址")
            return

        print("正在检查IP可达性...")
        self.remote_ips = self.filter_reachable_ips(all_ips)

        if not self.remote_ips:
            print("没有可连接的远程IP地址")
            return

        print(f"将使用以下IP进行转发: {self.remote_ips}")

        # 创建本地监听socket
        server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        server_socket.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
        server_socket.bind(('0.0.0.0', self.local_port))
        server_socket.listen(5)
        server_socket.settimeout(1)

        print(f"本地服务器已在端口 {self.local_port} 上启动")

        try:
            while True:
                try:
                    client_socket, addr = server_socket.accept()
                    print(f"接受来自 {addr} 的连接")

                    # 为每个连接创建新线程
                    client_thread = threading.Thread(
                        target=self.handle_local_connection,
                        args=(client_socket,)
                    )
                    client_thread.daemon = True
                    client_thread.start()

                except socket.timeout:
                    continue
                except Exception as e:
                    print(f"接受连接时发生错误: {e}")

        except KeyboardInterrupt:
            print("\n正在关闭服务器...")
        finally:
            server_socket.close()


def main():
    # 配置参数
    LOCAL_PORT = int(sys.argv[1])
    REMOTE_HOST = sys.argv[2]
    REMOTE_PORT = int(sys.argv[3])

    # 创建并启动转发器
    forwarder = TCPForwarder(LOCAL_PORT, REMOTE_HOST, REMOTE_PORT)
    forwarder.start_server()


if __name__ == "__main__":
    main()
