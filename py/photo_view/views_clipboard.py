import os
import threading

from flask import Blueprint, jsonify, render_template, request

from paths import clipboard_file

bp = Blueprint('clipboard', __name__)

MAX_CLIPBOARD_LENGTH = 1024 * 1024  # 1MB

clipboard_lock = threading.Lock()


def _read_clipboard():
    if not os.path.exists(clipboard_file):
        return ''
    try:
        with open(clipboard_file, 'r', encoding='utf-8') as f:
            return f.read()
    except Exception as e:
        print(f"读取剪切板文件失败: {e}")
        return ''


def _write_clipboard(text):
    tmp_path = clipboard_file + '.tmp'
    with open(tmp_path, 'w', encoding='utf-8') as f:
        f.write(text)
    os.replace(tmp_path, clipboard_file)


@bp.route("/clipboard")
def clipboard():
    """Shared clipboard page"""
    return render_template("clipboard.html")


@bp.route("/clipboard/content", methods=['GET'])
def clipboard_get():
    with clipboard_lock:
        return jsonify({'text': _read_clipboard()})


@bp.route("/clipboard/content", methods=['POST'])
def clipboard_set():
    data = request.get_json(silent=True) or {}
    text = data.get('text', '')
    if not isinstance(text, str):
        return jsonify({'ok': False, 'error': '内容格式错误'}), 400
    if len(text) > MAX_CLIPBOARD_LENGTH:
        return jsonify({'ok': False, 'error': '内容过长（上限 1MB）'}), 413

    with clipboard_lock:
        try:
            _write_clipboard(text)
        except Exception as e:
            print(f"写入剪切板文件失败: {e}")
            return jsonify({'ok': False, 'error': '保存失败'}), 500

    return jsonify({'ok': True})
