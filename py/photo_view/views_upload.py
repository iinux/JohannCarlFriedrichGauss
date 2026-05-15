import base64
import io
import os
import re
import time

from flask import Blueprint, flash, redirect, render_template, request, send_from_directory, url_for
from pdf2image import convert_from_path
from werkzeug.utils import secure_filename

from paths import ALLOWED_EXTENSIONS, upload_dir

bp = Blueprint('upload', __name__)


def allowed_file(filename):
    return '.' in filename and filename.rsplit('.', 1)[1].lower() in ALLOWED_EXTENSIONS


@bp.route("/upload", methods=['GET', 'POST'])
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
            timestamp = str(int(time.time()))
            unique_filename = f"upload_{timestamp}_{filename}"

            filepath = os.path.join(upload_dir, unique_filename)
            file.save(filepath)
            flash(f'文件 {filename} 上传成功！')
            return redirect(url_for('upload.upload'))
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


@bp.route("/upload/pdf", methods=['GET', 'POST'])
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

            return redirect(url_for('upload.upload_pdf'))
        else:
            flash('请上传PDF文件')
            return redirect(request.url)

    return render_template("upload_pdf.html")


@bp.route("/upload/delete/<filename>", methods=['POST'])
def delete_upload(filename):
    """删除已上传的文件"""
    # 安全检查：防止路径遍历攻击
    safe_name = re.sub(r'[^a-zA-Z0-9_\-.]', '', filename)
    filepath = os.path.join(upload_dir, safe_name)

    if os.path.exists(filepath) and os.path.isfile(filepath):
        os.remove(filepath)
        flash(f'文件 {filename} 已删除')

    return redirect(url_for('upload.upload'))


@bp.route("/upload/preview/<filename>")
def preview_upload(filename):
    """将 PDF 每页转为图片并在页面中展示"""
    safe_name = re.sub(r'[^a-zA-Z0-9_\-.]', '', filename)
    filepath = os.path.join(upload_dir, safe_name)
    images_b64 = []
    error = None
    try:
        pages = convert_from_path(filepath, dpi=150)
        for page in pages:
            buf = io.BytesIO()
            page.save(buf, 'JPEG')
            images_b64.append(base64.b64encode(buf.getvalue()).decode())
    except Exception as e:
        error = str(e)
    return render_template("upload_preview.html", filename=safe_name, images=images_b64, error=error)


@bp.route("/upload/download/<filename>")
def download_upload(filename):
    """下载已上传的文件"""
    safe_name = re.sub(r'[^a-zA-Z0-9_\-.]', '', filename)
    return send_from_directory(upload_dir, safe_name, as_attachment=True)
