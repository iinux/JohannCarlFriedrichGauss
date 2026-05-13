import random

from flask import Blueprint, render_template

bp = Blueprint('words', __name__)


@bp.route("/words")
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
