import os
import base64
import requests
import random
import subprocess
import json
import threading
import time
import my_config
from flask import Flask, render_template, request, redirect, url_for, flash, Response
from werkzeug.utils import secure_filename
from pdf2image import convert_from_path

work_dir = "."
mp4_dir = work_dir + '/mp4'
img_dir = work_dir + '/img'
upload_dir = work_dir + '/upload'
app = Flask(__name__, static_folder=work_dir + '/static', static_url_path='')
app.config['UPLOAD_FOLDER'] = upload_dir
app.config['MAX_CONTENT_LENGTH'] = 500 * 1024 * 1024  # 500MB max file size
app.secret_key = 'upload-secret-key-for-flash-messages'

ALLOWED_EXTENSIONS = {'mp4', 'jpg', 'jpeg', 'png', 'webp', 'gif', 'pdf'}

def allowed_file(filename):
    return '.' in filename and filename.rsplit('.', 1)[1].lower() in ALLOWED_EXTENSIONS


def get_video_duration(video_path):
    """使用 ffprobe 获取视频时长（秒）"""
    try:
        cmd = [
            'ffprobe',
            '-v', 'error',
            '-show_entries', 'format=duration',
            '-of', 'json',
            video_path
        ]
        result = subprocess.run(cmd, capture_output=True, text=True, check=True)
        data = json.loads(result.stdout)
        return float(data['format']['duration'])
    except Exception as e:
        print(f"获取视频时长失败: {e}")
        return None

def generate_random_time(duration):
    """生成随机时间点（格式 HH:MM:SS）"""
    if duration is None or duration <= 0:
        return "00:00:05"  # 默认值
    
    # 确保时间点在视频范围内（最后5秒前）
    max_time = max(0, duration - 5)
    random_seconds = random.uniform(0, max_time)
    
    # 转换为 HH:MM:SS 格式
    hours = int(random_seconds // 3600)
    minutes = int((random_seconds % 3600) // 60)
    seconds = random_seconds % 60
    return f"{hours:02d}:{minutes:02d}:{seconds:06.3f}"

@app.route("/index/random")
def random_page():
    # 1. 列出当前目录下所有 mp4 文件
    mp4_files = [f for f in os.listdir(mp4_dir) if f.endswith('.mp4')]
    
    if not mp4_files:
        print("当前目录下没有找到任何 mp4 文件")
        return
    
    # 2. 随机选择一个 mp4 文件
    selected_file = random.choice(mp4_files)
    selected_file_url = mp4_dir + '/' + selected_file
    print(f"随机选择的文件: {selected_file}")
    
    # 3. 构建对应的截图路径
    screenshot_dir = "screenshot"
    screenshot_path = os.path.join(screenshot_dir, os.path.splitext(selected_file)[0] + ".jpg")
    
    # 4. 检查截图是否存在
    if not os.path.exists(screenshot_path):
        print(f"截图不存在: {screenshot_path}")
        
        # 确保截图目录存在
        if not os.path.exists(screenshot_dir):
            os.makedirs(screenshot_dir)
            print(f"已创建目录: {screenshot_dir}")
        
        # 5. 获取视频时长并生成随机时间点
        duration = get_video_duration(selected_file_url)
        random_time = generate_random_time(duration)
        print(f"视频时长: {duration:.2f}秒, 随机时间点: {random_time}")
        
        # 6. 调用 ffmpeg 生成截图
        try:
            command = [
                'ffmpeg',
                '-ss', random_time,         # 随机时间点
                '-i', selected_file_url,        # 输入文件
                '-vframes', '1',            # 只取一帧
                '-q:v', '2',                # 高质量截图
                '-y',                       # 覆盖已存在文件
                screenshot_path             # 输出文件
            ]
            
            print(f"正在生成截图: {' '.join(command)}")
            result = subprocess.run(command, capture_output=True, text=True)
            
            if os.path.exists(screenshot_path):
                print(f"截图已成功生成: {screenshot_path}")
            else:
                print(f"截图生成失败，请检查 ffmpeg 是否安装正确")
                if result.stderr:
                    print(f"ffmpeg 错误输出:\n{result.stderr[:500]}...")  # 只显示前500字符
                    
        except Exception as e:
            print(f"执行 ffmpeg 命令出错: {e}")
    else:
        print(f"截图已存在: {screenshot_path}")

    return render_template("random.html", screenshot=screenshot_path, mp4=selected_file_url)

def get_images():
    images = []
    for filename in os.listdir(img_dir):
        if filename.endswith(".jpg") or filename.endswith(".png") or filename.endswith(".webp"):
            images.append(filename)
    return images


def get_mp4s():
    mp4s = []
    path = mp4_dir
    files_and_dirs = os.listdir(path)

    # 过滤掉目录，只对文件进行排序
    files = [f for f in files_and_dirs if os.path.isfile(os.path.join(path, f))]

    # 按文件大小排序
    sorted_files = sorted(files, key=lambda x: os.path.getsize(os.path.join(path, x)))

    # 按文件时间排序
    #sorted_files = sorted(files_and_dirs, key=lambda x: os.path.getmtime(os.path.join(path, x)))

    for filename in sorted_files:
        if filename.endswith(".mp4"):
            fn = {}
            if filename.startswith("upload"):
                fn['show_name'] = filename[39:]
            else:
                fn['show_name'] = filename
            fn['real_name'] = filename
            mp4s.append(fn)
    return mp4s

def get_flvs():
    flvs = []
    path = './flv'
    files_and_dirs = os.listdir(path)

    sorted_files = files_and_dirs

    for filename in sorted_files:
        if filename.endswith(".flv"):
            fn = {}
            fn['show_name'] = filename
            fn['real_name'] = filename
            flvs.append(fn)
    return flvs

@app.route("/index")
def index():
    images = get_images()
    mp4s = get_mp4s()
    flvs = get_flvs()
    return render_template("index.html", images=images, mp4s=mp4s, flvs=flvs)

@app.route("/index/image/<filename>")
def image(filename):
    """图片详情页"""
    image_data = open(img_dir + '/' + filename, "rb").read()
    base64_data = base64.b64encode(image_data).decode()
    return render_template("image.html", filename=filename, base64_data=base64_data)


@app.route("/words")
def words():
    """Random vocabulary page"""
    dict_path = "dict2.txt"
    lines = []
    try:
        with open(dict_path, 'r') as f:
            lines = [line.strip() for line in f if line.strip()]
    except Exception as e:
        print(f"读取词典文件失败: {e}")
    
    if len(lines) > 10:
        random_lines = random.sample(lines, 10)
    else:
        random_lines = lines
    
    return render_template("words.html", words=random_lines)


@app.route("/rtsp/ch00")
def rtsp_ch00():
    """RTSP 视频流 - 通道 00"""
    return render_template("rtsp_ch00.html")

@app.route("/rtsp/ch01")
def rtsp_ch01():
    """RTSP 视频流 - 通道 01"""
    return render_template("rtsp_ch01.html")


# RTSP 流配置 - 请根据实际摄像头修改
RTSP_CONFIGS = {
    'ch00': {
        'rtsp_url': my_config.rtsp_url1,
        'width': 1920,
        'height': 1080
    },
    'ch01': {
        'rtsp_url': my_config.rtsp_url2,
        'width': 1920,
        'height': 1080
    }
}

# 每个频道只跑一个 ffmpeg，最新一帧写入 stream_cache，所有 HTTP 请求从这里读
stream_cache = {}
stream_lock = threading.Lock()
stream_threads = {}


def _drain_stderr(process):
    """持续读取 ffmpeg stderr，避免 64KB 管道写满后阻塞主进程"""
    try:
        for _ in iter(process.stderr.readline, b''):
            pass
    except Exception:
        pass


def _stream_reader(channel):
    """后台线程：跑 ffmpeg、解析 JPEG 帧、写入 stream_cache，异常时自动重连"""
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
        process = None
        try:
            process = subprocess.Popen(cmd, stdout=subprocess.PIPE, stderr=subprocess.PIPE, bufsize=0)
            threading.Thread(target=_drain_stderr, args=(process,), daemon=True).start()

            buffer = b''
            while True:
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

        print(f"RTSP 流 {channel} ffmpeg 退出，3 秒后重连")
        time.sleep(3)


def _ensure_stream(channel):
    """首次请求时启动该频道的后台读取线程"""
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
    """从 stream_cache 取最新帧，输出 MJPEG multipart"""
    _ensure_stream(channel)
    last_frame = None
    while True:
        with stream_lock:
            frame = stream_cache.get(channel)
        if frame is not None and frame is not last_frame:
            last_frame = frame
            yield (b'--frame\r\n'
                   b'Content-Type: image/jpeg\r\n\r\n' + frame + b'\r\n')
        time.sleep(0.04)


@app.route("/rtsp/stream/<channel>")
def rtsp_stream(channel):
    """RTSP 流 MJPEG 视频流"""
    if channel not in RTSP_CONFIGS:
        return "Channel not found", 404

    return Response(
        generate_frames(channel),
        mimetype='multipart/x-mixed-replace; boundary=frame'
    )


@app.route("/rtsp/feed/<channel>")
def rtsp_feed(channel):
    """RTSP 流单帧图片 (供前端轮询)"""
    if channel not in RTSP_CONFIGS:
        return "Channel not found", 404

    _ensure_stream(channel)
    with stream_lock:
        frame = stream_cache.get(channel)

    if frame:
        return Response(frame, mimetype='image/jpeg')
    else:
        return "No frame", 404

@app.route("/upload", methods=['GET', 'POST'])
def upload():
    """上传页面及处理"""
    if request.method == 'POST':
        # 检查是否有文件
        if 'file' not in request.files:
            flash('没有选择文件')
            return redirect(request.url)
        
        file = request.files['file']
        if file.filename == '':
            flash('没有选择文件')
            return redirect(request.url)
        
        if file and allowed_file(file.filename):
            # 确保上传目录存在
            if not os.path.exists(upload_dir):
                os.makedirs(upload_dir)
            
            filename = secure_filename(file.filename)
            # 添加时间戳避免文件名冲突
            import time
            timestamp = str(int(time.time()))
            name, ext = os.path.splitext(filename)
            unique_filename = f"upload_{timestamp}_{filename}"
            
            filepath = os.path.join(upload_dir, unique_filename)
            file.save(filepath)
            flash(f'文件 {filename} 上传成功！')
            return redirect(url_for('upload'))
        else:
            flash(f'不支持的文件类型，允许的格式: {", ".join(ALLOWED_EXTENSIONS)}')
            return redirect(request.url)
    
    # 列出已上传的文件
    uploaded_files = []
    if os.path.exists(upload_dir):
        for f in sorted(os.listdir(upload_dir), key=lambda x: os.path.getmtime(os.path.join(upload_dir, x)), reverse=True):
            full_path = os.path.join(upload_dir, f)
            if os.path.isfile(full_path):
                size = os.path.getsize(full_path)
                # 格式化文件大小
                if size < 1024 * 1024:
                    size_str = f"{size / 1024:.1f} KB"
                else:
                    size_str = f"{size / (1024 * 1024):.1f} MB"
                uploaded_files.append({'name': f, 'size': size_str})
    
    return render_template("upload.html", uploaded_files=uploaded_files, allowed_extensions=ALLOWED_EXTENSIONS)


@app.route("/upload/pdf", methods=['GET', 'POST'])
def upload_pdf():
    """上传PDF并转换为图片"""
    if request.method == 'POST':
        if 'file' not in request.files:
            flash('没有选择文件')
            return redirect(request.url)

        file = request.files['file']
        if file.filename == '':
            flash('没有选择文件')
            return redirect(request.url)

        if file and file.filename.lower().endswith('.pdf'):
            if not os.path.exists(upload_dir):
                os.makedirs(upload_dir)

            filename = secure_filename(file.filename)
            import time
            timestamp = str(int(time.time()))
            unique_filename = f"upload_{timestamp}_{filename}"
            filepath = os.path.join(upload_dir, unique_filename)
            file.save(filepath)

            try:
                images = convert_from_path(filepath, dpi=200)
                converted_names = []
                for i, image in enumerate(images):
                    img_filename = f"pdf_{timestamp}_page_{i + 1}.jpg"
                    img_path = os.path.join(upload_dir, img_filename)
                    image.save(img_path, "JPEG")
                    converted_names.append(img_filename)

                os.remove(filepath)
                flash(f'PDF转换成功! 共转换 {len(converted_names)} 页')
            except Exception as e:
                flash(f'PDF转换失败: {str(e)}')
                if os.path.exists(filepath):
                    os.remove(filepath)

            return redirect(url_for('upload_pdf'))
        else:
            flash('请上传PDF文件')
            return redirect(request.url)

    return render_template("upload_pdf.html")


@app.route("/upload/delete/<filename>", methods=['POST'])
def delete_upload(filename):
    """删除已上传的文件"""
    import re
    # 安全检查：防止路径遍历攻击
    safe_name = re.sub(r'[^a-zA-Z0-9_\-.]', '', filename)
    filepath = os.path.join(upload_dir, safe_name)

    if os.path.exists(filepath) and os.path.isfile(filepath):
        os.remove(filepath)
        flash(f'文件 {filename} 已删除')

    return redirect(url_for('upload'))


@app.route("/upload/download/<filename>")
def download_upload(filename):
    """下载已上传的文件"""
    import re
    from flask import send_from_directory
    safe_name = re.sub(r'[^a-zA-Z0-9_\-.]', '', filename)
    return send_from_directory(upload_dir, safe_name, as_attachment=True)


if __name__ == "__main__":
    #app.run(host='192.168.1.4')
    app.run()
