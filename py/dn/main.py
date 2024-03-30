import cv2
import sys

from run import process
from PIL import Image


def resize_image(input_path, output_path, target_size=(512, 512)):
    try:
        with Image.open(input_path) as image:
            resized_image = image.resize(target_size)
            resized_image.save(output_path, "PNG")
            print("图片已成功调整大小并保存为 PNG 格式！")
    except Exception as e:
        print("出现错误:", e)


def main():
    input_image_path = sys.argv[1]
    resize_image_path = "resize_image.png"
    resize_image(input_image_path, resize_image_path)

    # Read input image
    dress = cv2.imread(resize_image_path)

    watermark = process(dress)

    # Write output image
    cv2.imwrite(sys.argv[2], watermark)

    # Exit
    sys.exit()


if __name__ == '__main__':
    main()
