import os
import base64
import requests
import random
import subprocess
import json
from flask import Flask, render_template, request

img_dir = "."
app = Flask(__name__, static_folder=img_dir, static_url_path='')


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
    mp4_files = [f for f in os.listdir(img_dir) if f.endswith('.mp4')]
    
    if not mp4_files:
        print("当前目录下没有找到任何 mp4 文件")
        return
    
    # 2. 随机选择一个 mp4 文件
    selected_file = random.choice(mp4_files)
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
        duration = get_video_duration(selected_file)
        random_time = generate_random_time(duration)
        print(f"视频时长: {duration:.2f}秒, 随机时间点: {random_time}")
        
        # 6. 调用 ffmpeg 生成截图
        try:
            command = [
                'ffmpeg',
                '-ss', random_time,         # 随机时间点
                '-i', selected_file,        # 输入文件
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

    return render_template("random.html", screenshot=screenshot_path, mp4=selected_file)

def get_images():
    images = []
    for filename in os.listdir(img_dir):
        if filename.endswith(".jpg") or filename.endswith(".png"):
            images.append(filename)
    return images


def get_mp4s():
    mp4s = []
    path = img_dir
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

@app.route("/image/<filename>")
def image(filename):
    """图片详情页"""
    image_data = open(img_dir + '/' + filename, "rb").read()
    base64_data = base64.b64encode(image_data).decode()
    return render_template("image.html", filename=filename, base64_data=base64_data)


if __name__ == "__main__":
    #app.run(host='192.168.1.4')
    app.run()
