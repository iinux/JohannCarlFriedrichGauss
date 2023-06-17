#!/usr/bin/env python

from fuse import FUSE, FuseOSError, Operations

class MyFileSystem(Operations):
    def __init__(self, root):
        self.root = root

    def getattr(self, path, fh=None):
        # 实现获取文件属性的回调函数
        # 根据路径返回文件的属性，例如大小、权限、时间等
        print('getattr')
        pass

    def readdir(self, path, fh):
        # 实现读取目录内容的回调函数
        # 根据路径返回目录下的文件和子目录列表
        print('readdir')
        pass

    def read(self, path, length, offset, fh):
        # 实现读取文件内容的回调函数
        # 根据路径、偏移量和长度返回文件内容
        print('read')
        pass

    # 其他回调函数，根据需要实现相应的逻辑

if __name__ == '__main__':
    # 设置要挂载的根目录
    root = '/home/qzhang/git/JohannCarlFriedrichGauss/py/mount_point'

    # 创建文件系统实例
    fs = MyFileSystem(root)

    # 挂载文件系统到指定的挂载点
    mount_point = '/home/qzhang/git/JohannCarlFriedrichGauss/py/mount_point'
    FUSE(fs, mount_point, foreground=True)

