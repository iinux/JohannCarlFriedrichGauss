#!/usr/bin/env python3
"""PDF转JPG工具"""

# pip install pdf2image Pillow
# 并安装 poppler-utils
# sudo apt install poppler-utils

import os
import sys
from pdf2image import convert_from_path
from PIL import Image

def pdf_to_jpg(pdf_path, output_dir=".", dpi=200):
    """
    将PDF转换为JPG图片

    Args:
        pdf_path: PDF文件路径
        output_dir: 输出目录，默认当前目录
        dpi: 图片分辨率，默认200
    """
    if not os.path.exists(pdf_path):
        print(f"错误: 文件不存在 - {pdf_path}")
        return

    if not os.path.exists(output_dir):
        os.makedirs(output_dir)

    images = convert_from_path(pdf_path, dpi=dpi)

    for i, image in enumerate(images):
        jpg_path = os.path.join(output_dir, f"page_{i + 1}.jpg")
        image.save(jpg_path, "JPEG")
        print(f"已保存: {jpg_path}")

    print(f"转换完成! 共 {len(images)} 页")

if __name__ == '__main__':
    if len(sys.argv) < 2:
        print("用法: python pdf_to_jpg.py <pdf文件> [输出目录] [dpi]")
        print("示例: python pdf_to_jpg.py input.pdf ./output 300")
        sys.exit(1)

    pdf_path = sys.argv[1]
    output_dir = sys.argv[2] if len(sys.argv) > 2 else "."
    dpi = int(sys.argv[3]) if len(sys.argv) > 3 else 200

    pdf_to_jpg(pdf_path, output_dir, dpi)
