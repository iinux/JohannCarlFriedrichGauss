import os
from pdf2image import convert_from_path
from PIL import Image


# pip install pdf2image Pillow
# sudo apt-get install poppler-utils
# sudo yum install poppler-utils
# brew install poppler


def pdf_to_jpg_page_by_page(pdf_path, output_dir):
    # 检查输出目录是否存在，如果不存在则创建
    if not os.path.exists(output_dir):
        os.makedirs(output_dir)

    # 将PDF转换为图像
    images = convert_from_path(pdf_path)

    # 保存每个图像为JPG文件
    for i, image in enumerate(images):
        jpg_path = os.path.join(output_dir, f"page_{i + 1}.jpg")
        image.save(jpg_path, "JPEG")

    print("转换完成！")


def pdf_to_one_jpg(pdf_path):
    # 读取PDF并生成图片列表
    images = convert_from_path(pdf_path)

    # 将图片合并为一张大图
    width = sum(image.width for image in images)
    height = max(image.height for image in images)
    result_image = Image.new('RGB', (width, height))

    x_offset = 0
    for image in images:
        result_image.paste(image, (x_offset, 0))
        x_offset += image.width

    if __name__ == '__main__':
        # 保存合并后的图片
        result_image.save('output.jpg')


def pdf_to_one_jpg2(pdf_path):
    # 读取PDF并生成图片列表
    images = convert_from_path(pdf_path)

    # 计算结果图片的尺寸
    width = max(image.width for image in images)
    height = sum(image.height for image in images)

    # 创建结果图片对象
    result_image = Image.new('RGB', (width, height))

    y_offset = 0
    for image in images:
        result_image.paste(image, (0, y_offset))
        y_offset += image.height

    # 保存合并后的图片
    result_image.save('output2.jpg')


pdf_path = "kid.pdf"
output_dir = "."

if __name__ == '__main__':
    # 调用函数进行转换
    pdf_to_jpg_page_by_page(pdf_path, output_dir)
    pdf_to_one_jpg(pdf_path)
    pdf_to_one_jpg2(pdf_path)
