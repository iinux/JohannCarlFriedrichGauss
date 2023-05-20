# 导入easyocr
import easyocr

# 创建reader对象
reader = easyocr.Reader(['ch_sim', 'en'])
# 读取图像
result = reader.readtext('test1.jpg')
# 结果
print(result)

for i in result:
    word = i[1]
    print(word)

if __name__ == '__main__':
    pass
