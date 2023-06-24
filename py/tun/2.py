import struct

# 将整数和浮点数打包成二进制数据
data = struct.pack('i f', 42, 3.14)
print(data)  # 输出 b'*\x00\x00\x00\x0c\x8f\xc2@'

