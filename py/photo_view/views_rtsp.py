import subprocess
import threading
import time
import re

import my_config
from flask import Blueprint, Response, render_template

bp = Blueprint('rtsp', __name__)


def find_ip_by_mac(target_mac):
    """通过 MAC 地址查找局域网设备的 IP 地址"""
    try:
        # 先尝试从 arp 缓存中查找
        result = subprocess.run(['arp', '-a'], capture_output=True, text=True, timeout=10)
        for line in result.stdout.split('\n'):
            # 格式: host (192.168.1.1) at 98:03:cf:a4:64:1c on enp0s3 ...
            match = re.search(r'\((\d+\.\d+\.\d+\.\d+)\)\s+at\s+([0-9a-fA-F:]+)', line)
            if match:
                ip = match.group(1)
                mac = match.group(2).lower().replace('-', ':')
                if mac == target_mac.lower().replace('-', ':'):
                    return ip

        # 如果 arp 缓存没有，尝试 ping 整个网段再查
        # 获取本机 IP 和子网
        local_ip = subprocess.run(
            ["hostname", "-I"], capture_output=True, text=True, timeout=5
        ).stdout.strip().split()[0]
        subnet = '.'.join(local_ip.split('.')[:3])

        # ping 扫描子网
        subprocess.run(
            f'for i in $(seq 1 254); do ping -c 1 -W 1 {subnet}.$i & done',
            shell=True, capture_output=True, timeout=30
        )
        time.sleep(2)

        # 再次尝试从 arp 缓存中查找
        result = subprocess.run(['arp', '-a'], capture_output=True, text=True, timeout=10)
        for line in result.stdout.split('\n'):
            match = re.search(r'\((\d+\.\d+\.\d+\.\d+)\)\s+at\s+([0-9a-fA-F:]+)', line)
            if match:
                ip = match.group(1)
                mac = match.group(2).lower().replace('-', ':')
                if mac == target_mac.lower().replace('-', ':'):
                    return ip
    except Exception as e:
        print(f"查找 IP 失败: {e}")
    return None


# 启动前解析 RTSP URL 中的 <ip> 占位符
def _resolve_rtsp_url(template_url):
    """解析 RTSP URL 中的 <ip> 占位符为实际 IP"""
    if '<ip>' in template_url:
        ip = find_ip_by_mac(my_config.mac)
        if ip:
            return template_url.replace('<ip>', ip)
        print(f"未找到 MAC {my_config.mac} 对应的 IP，使用原 URL")
    return template_url


# RTSP 流配置 - 请根据实际摄像头修改
RTSP_CONFIGS = {
    'ch00': {
        'rtsp_url': _resolve_rtsp_url(my_config.rtsp_url1),
        'width': 1920,
        'height': 1080
    },
    'ch01': {
        'rtsp_url': _resolve_rtsp_url(my_config.rtsp_url2),
        'width': 1920,
        'height': 1080
    }
}

# 每个频道按需启动单个 ffmpeg：有人看才跑，全部空闲后关掉
stream_cache = {}
stream_lock = threading.Lock()
stream_threads = {}
stream_viewers = {}      # channel -> 当前活跃 MJPEG 客户端数
stream_last_feed = {}    # channel -> 最近一次 /rtsp/feed 访问时间戳
STREAM_FEED_GRACE = 30   # 单帧轮询的宽限期（秒）；过后视为空闲


def _stream_active(channel):
    """是否还有活跃使用者：长连接观众或最近的轮询请求"""
    if stream_viewers.get(channel, 0) > 0:
        return True
    return (time.time() - stream_last_feed.get(channel, 0)) < STREAM_FEED_GRACE


def _drain_stderr(process):
    """持续读取 ffmpeg stderr，避免 64KB 管道写满后阻塞主进程"""
    try:
        for _ in iter(process.stderr.readline, b''):
            pass
    except Exception:
        pass


def _stream_reader(channel):
    """后台线程：跑 ffmpeg、解析 JPEG 帧、写入 stream_cache；空闲时自我注销并退出"""
    config = RTSP_CONFIGS[channel]
    cmd = [
        'ffmpeg',
        '-rtsp_transport', 'tcp',
        '-i', config['rtsp_url'],
        '-f', 'mjpeg',
        '-q:v', '5',
        '-vf', f'scale={config["width"]}:{config["height"]}',
        '-'
    ]

    while True:
        # 启动 ffmpeg 之前先确认仍有人需要
        with stream_lock:
            if not _stream_active(channel):
                stream_threads.pop(channel, None)
                stream_cache.pop(channel, None)
                print(f"RTSP 流 {channel} 空闲，停止 ffmpeg")
                return

        process = None
        try:
            process = subprocess.Popen(cmd, stdout=subprocess.PIPE, stderr=subprocess.PIPE, bufsize=0)
            threading.Thread(target=_drain_stderr, args=(process,), daemon=True).start()

            buffer = b''
            idle_check_at = time.time()
            while True:
                # 每秒检查一次是否还有人在看，没人就退出循环把 ffmpeg 杀掉
                now = time.time()
                if now - idle_check_at >= 1.0:
                    idle_check_at = now
                    with stream_lock:
                        if not _stream_active(channel):
                            break

                chunk = process.stdout.read(8192)
                if not chunk:
                    break
                buffer += chunk
                while True:
                    soi = buffer.find(b'\xff\xd8')
                    if soi < 0:
                        buffer = b''
                        break
                    eoi = buffer.find(b'\xff\xd9', soi + 2)
                    if eoi < 0:
                        buffer = buffer[soi:]
                        break
                    frame = buffer[soi:eoi + 2]
                    buffer = buffer[eoi + 2:]
                    with stream_lock:
                        stream_cache[channel] = frame
        except Exception as e:
            print(f"RTSP 流 {channel} 读取错误: {e}")
        finally:
            if process is not None:
                try:
                    process.kill()
                    process.wait(timeout=5)
                except Exception:
                    pass

        # ffmpeg 退出后再决定：还有人看就重连，没人看就退出
        with stream_lock:
            if not _stream_active(channel):
                stream_threads.pop(channel, None)
                stream_cache.pop(channel, None)
                print(f"RTSP 流 {channel} 空闲，停止 ffmpeg")
                return

        print(f"RTSP 流 {channel} ffmpeg 退出，3 秒后重连")
        time.sleep(3)


def _ensure_stream(channel):
    """请求到来时确保该频道的后台读取线程在跑"""
    if channel not in RTSP_CONFIGS:
        return
    with stream_lock:
        t = stream_threads.get(channel)
        if t and t.is_alive():
            return
        t = threading.Thread(target=_stream_reader, args=(channel,), daemon=True)
        stream_threads[channel] = t
        t.start()


def generate_frames(channel):
    """从 stream_cache 取最新帧输出 MJPEG multipart；以引用计数维持 ffmpeg 存活"""
    with stream_lock:
        stream_viewers[channel] = stream_viewers.get(channel, 0) + 1
    _ensure_stream(channel)
    try:
        last_frame = None
        while True:
            with stream_lock:
                frame = stream_cache.get(channel)
            if frame is not None and frame is not last_frame:
                last_frame = frame
                yield (b'--frame\r\n'
                       b'Content-Type: image/jpeg\r\n\r\n' + frame + b'\r\n')
            time.sleep(0.04)
    finally:
        with stream_lock:
            stream_viewers[channel] = max(0, stream_viewers.get(channel, 0) - 1)


@bp.route("/rtsp/ch00")
def rtsp_ch00():
    """RTSP 视频流 - 通道 00"""
    return render_template("rtsp_ch00.html")


@bp.route("/rtsp/ch01")
def rtsp_ch01():
    """RTSP 视频流 - 通道 01"""
    return render_template("rtsp_ch01.html")


@bp.route("/rtsp/stream/<channel>")
def rtsp_stream(channel):
    """RTSP 流 MJPEG 视频流"""
    if channel not in RTSP_CONFIGS:
        return "Channel not found", 404

    return Response(
        generate_frames(channel),
        mimetype='multipart/x-mixed-replace; boundary=frame'
    )


@bp.route("/rtsp/feed/<channel>")
def rtsp_feed(channel):
    """RTSP 流单帧图片 (供前端轮询)"""
    if channel not in RTSP_CONFIGS:
        return "Channel not found", 404

    with stream_lock:
        stream_last_feed[channel] = time.time()
    _ensure_stream(channel)
    with stream_lock:
        frame = stream_cache.get(channel)

    if frame:
        return Response(frame, mimetype='image/jpeg')
    else:
        return "No frame", 404
