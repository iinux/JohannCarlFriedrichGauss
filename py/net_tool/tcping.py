import socket
import time
import sys


def tcping(host, port, timeout=5):
    try:
        start_time = time.time()
        sock = socket.create_connection((host, port), timeout)
        sock.sendall(b'A')
        received_data = sock.recv(1)
        # print("receive: {}".format(received_data))
        sock.close()
        end_time = time.time()

        rtt = (end_time - start_time) * 1000  # to ms
        print('TCPing to {}: {}, RTT: {:.2f} ms'.format(host, port, rtt))
        # print(f'TCPing to {host}:{port}, RTT: {rtt:.2f} ms')
        return [rtt]

    except (socket.timeout, socket.error) as e:
        print('TCPing to {}: {}, Error: {}'.format(host, port, e))
        # print(f'TCPing to {host}:{port}, Error: {e}')


def print_detail(detail):
    count = len(detail)
    max_value = max(detail)
    min_value = min(detail)
    avg_value = sum(detail) / count

    print('\n{} times, max {:.2f} ms, min {:.2f} ms, avg {:.2f} ms' .format(count, max_value, min_value, avg_value))


if __name__ == "__main__":
    target_host = sys.argv[1]
    target_port = sys.argv[2]

    max_time = 0
    if len(sys.argv) > 3:
        max_time = int(sys.argv[3])

    count = 0
    detail = []

    try:
        while True:
            detail += tcping(target_host, target_port)
            count += 1
            if max_time != 0 and count >= max_time:
                break
            time.sleep(1)
        print_detail(detail)
    except:
        print_detail(detail)
