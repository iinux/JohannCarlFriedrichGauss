import base64
import io
import os
import re
import time

from flask import Blueprint, flash, redirect, render_template, request, send_from_directory, url_for
from pdf2image import convert_from_path
from PIL import Image
from werkzeug.utils import secure_filename

from paths import ALLOWED_EXTENSIONS, upload_dir

bp = Blueprint('upload', __name__)

# rembg is optional. If installed it gives ML-based background removal
# (best quality). Otherwise we fall back to a simple PIL corner-color
# threshold method that only works on plain-color backgrounds.
try:
    from rembg import remove as _rembg_remove
    HAS_REMBG = True
except Exception:
    HAS_REMBG = False


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
            timestamp = time.strftime("%Y%m%d_%H%M%S", time.gmtime(time.time() + 8 * 3600))
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
            timestamp = time.strftime("%Y%m%d_%H%M%S", time.gmtime(time.time() + 8 * 3600))
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


# ---------------------------------------------------------------------------
# Background removal
# ---------------------------------------------------------------------------

BG_IMAGE_EXTS = {'jpg', 'jpeg', 'png', 'webp'}


def _remove_bg_pil(img, tolerance=30):
    """PIL fallback: treat the brightest corner color as background and
    make all pixels within `tolerance` transparent. Works only when the
    background is a single plain color (typical for ID photos etc.)."""
    rgba = img.convert('RGBA')
    w, h = rgba.size
    corners = [rgba.getpixel((0, 0)), rgba.getpixel((w - 1, 0)),
               rgba.getpixel((0, h - 1)), rgba.getpixel((w - 1, h - 1))]
    bg = max(corners, key=lambda p: p[0] + p[1] + p[2])
    br, bgc, bb, _ = bg

    pixels = rgba.load()
    tol = int(tolerance)
    for y in range(h):
        for x in range(w):
            r, g, b, _a = pixels[x, y]
            if abs(r - br) <= tol and abs(g - bgc) <= tol and abs(b - bb) <= tol:
                pixels[x, y] = (255, 255, 255, 0)
    return rgba


def _remove_bg(img_bytes):
    """Return PNG bytes with the background removed. Uses rembg when
    available, otherwise the PIL fallback."""
    if HAS_REMBG:
        return _rembg_remove(img_bytes)
    img = Image.open(io.BytesIO(img_bytes))
    out = _remove_bg_pil(img)
    buf = io.BytesIO()
    out.save(buf, format='PNG')
    return buf.getvalue()


@bp.route("/upload/bg", methods=['GET', 'POST'])
def upload_bg():
    """Upload an image, strip its background, save the result as PNG."""
    result_b64 = None
    original_b64 = None
    output_name = None
    error = None

    if request.method == 'POST':
        if 'file' not in request.files:
            error = 'No file selected'
        else:
            f = request.files['file']
            if f.filename == '':
                error = 'No file selected'
            elif not allowed_file(f.filename):
                error = 'Unsupported file type'
            else:
                ext = f.filename.rsplit('.', 1)[1].lower()
                if ext not in BG_IMAGE_EXTS:
                    error = 'Only image files are supported'
                else:
                    if not os.path.exists(upload_dir):
                        os.makedirs(upload_dir)
                    raw = f.read()
                    timestamp = time.strftime(
                        "%Y%m%d_%H%M%S",
                        time.gmtime(time.time() + 8 * 3600))
                    output_name = f"nobg_{timestamp}.png"
                    output_path = os.path.join(upload_dir, output_name)
                    try:
                        result_bytes = _remove_bg(raw)
                        with open(output_path, 'wb') as out:
                            out.write(result_bytes)
                        result_b64 = base64.b64encode(result_bytes).decode()
                        original_b64 = base64.b64encode(raw).decode()
                    except Exception as e:
                        error = f'Background removal failed: {e}'

    return render_template(
        "upload_bg.html",
        result_b64=result_b64,
        original_b64=original_b64,
        output_name=output_name,
        error=error,
        has_rembg=HAS_REMBG,
        allowed_image_exts=', '.join(sorted(BG_IMAGE_EXTS)),
    )


@bp.route("/upload/bg/download/")
def download_bg(filename):
    """Download the background-removed image."""
    safe_name = re.sub(r'[^a-zA-Z0-9_\-.]', '', filename)
    return send_from_directory(upload_dir, safe_name, as_attachment=True)
