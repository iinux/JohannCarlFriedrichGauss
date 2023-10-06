import os
import base64
import requests
from flask import Flask, render_template, request

app = Flask(__name__)

img_dir = "."

def get_images():
    """获取当前目录下的所有图片"""
    images = []
    for filename in os.listdir(img_dir):
        if filename.endswith(".jpg") or filename.endswith(".png"):
            images.append(filename)
    return images


@app.route("/")
def index():
    """主页"""
    images = get_images()
    return render_template("index.html", images=images)


@app.route("/image/<filename>")
def image(filename):
    """图片详情页"""
    image_data = open(img_dir + '/' + filename, "rb").read()
    base64_data = base64.b64encode(image_data).decode()
    return render_template("image.html", filename=filename, base64_data=base64_data)


if __name__ == "__main__":
    app.run(debug=True)

