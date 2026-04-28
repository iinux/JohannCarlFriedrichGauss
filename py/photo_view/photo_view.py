import os
import base64
import requests
import random
import subprocess
import json
from flask import Flask, render_template, request, redirect, url_for, flash
from werkzeug.utils import secure_filename

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
    dict_path = "dict.txt"
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
    
    return render_template("upload.html", uploaded_files=uploaded_files)


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


if __name__ == "__main__":
    #app.run(host='192.168.1.4')
    app.run()
