import socket
import subprocess
import time
import threading
import sys
from concurrent.futures import ThreadPoolExecutor, as_completed


def resolve_domain(domain):
    """
    解析域名获取所有IP地址
    """
    try:
        addr_info = socket.getaddrinfo(domain, None)
        ips = list(set([ip[4][0] for ip in addr_info]))  # 去重
        return ips
    except Exception as e:
        print(f"域名解析失败: {e}")
        return []


def ping_test(ip, count=4):
    """
    对指定IP进行ping测试
    返回平均延迟(ms)和丢包率
    """
    try:
        # 根据操作系统选择ping命令
        import platform
        if platform.system().lower() == "windows":
            cmd = ["ping", "-n", str(count), "-w", "timeout=1000", ip]
        else:
            cmd = ["ping", "-c", str(count), "-W", "1", ip]

        start_time = time.time()
        result = subprocess.run(cmd, capture_output=True, text=True, timeout=10)
        end_time = time.time()

        if result.returncode == 0:
            output = result.stdout
            # 解析ping结果
            if "ttl" in output.lower() or "time=" in output.lower():
                # 提取平均延迟
                import re
                times = re.findall(r"time[=<](\d+\.?\d*)", output)
                if times:
                    avg_time = sum(float(t) for t in times) / len(times)
                    return round(avg_time, 2)
                else:
                    # 粗略计算平均时间
                    return round((end_time - start_time) * 1000 / count, 2)
        return None
    except Exception as e:
        print(f"Ping {ip} 失败: {e}")
        return None


def tcp_port_test(ip, port, timeout=3):
    """
    测试TCP端口连通性和连接时间
    """
    try:
        start_time = time.time()
        sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        sock.settimeout(timeout)
        result = sock.connect_ex((ip, port))
        end_time = time.time()
        sock.close()

        if result == 0:
            connect_time = round((end_time - start_time) * 1000, 2)
            return connect_time
        else:
            return None
    except Exception as e:
        return None


def test_single_ip(ip, port):
    """
    对单个IP进行全面测试
    """
    print(f"正在测试 {ip} ...")

    # Ping测试
    ping_result = ping_test(ip)

    # TCP端口测试
    if port is not None:
        tcp_result = tcp_port_test(ip, port)
    else:
        tcp_result = None

    return {
        'ip': ip,
        'ping': ping_result,
        'tcp_connect': tcp_result,
        'status': 'reachable' if tcp_result is not None else 'unreachable'
    }


def main():
    args_len = len(sys.argv)
    domain = sys.argv[1] if args_len > 1 else "tcp.ap-northeast-1.clawcloudrun.com"
    port = int(sys.argv[2]) if args_len > 2 else None

    print(f"正在解析域名: {domain}")
    ips = resolve_domain(domain)

    if not ips:
        print("未能解析到任何IP地址")
        return

    print(f"解析到 {len(ips)} 个IP地址: {ips}")
    print("=" * 60)
    print(f"{'IP地址':<18} {'Ping延迟(ms)':<15} {'TCP连接(ms)':<15} {'状态':<10}")
    print("-" * 60)

    # 使用线程池并发测试所有IP
    results = []
    with ThreadPoolExecutor(max_workers=len(ips)) as executor:
        # 提交所有测试任务
        future_to_ip = {executor.submit(test_single_ip, ip, port): ip for ip in ips}

        # 收集结果
        for future in as_completed(future_to_ip):
            try:
                result = future.result()
                results.append(result)
            except Exception as e:
                ip = future_to_ip[future]
                print(f"测试 {ip} 时发生错误: {e}")

    # 按照TCP连接时间排序
    results.sort(key=lambda x: (x['tcp_connect'] is None, x['tcp_connect'] or float('inf')))

    # 输出结果
    for result in results:
        ping_str = f"{result['ping']}ms" if result['ping'] is not None else "N/A"
        tcp_str = f"{result['tcp_connect']}ms" if result['tcp_connect'] is not None else "N/A"
        status_str = result['status']

        print(f"{result['ip']:<18} {ping_str:<15} {tcp_str:<15} {status_str:<10}")

    print("=" * 60)
    reachable_count = sum(1 for r in results if r['status'] == 'reachable')
    print(f"总计: {len(ips)} 个IP, {reachable_count} 个可达, {len(ips) - reachable_count} 个不可达")


if __name__ == "__main__":
    main()
