#!/usr/bin/env python
# -*- coding: utf-8 -*-

import os
import sys
import zipfile

# print "Processing File " + sys.argv[1]

file = zipfile.ZipFile(sys.argv[1], "r")
for name in file.namelist():
    # /usr/lib/python3.9/zipfile.py: 1359
    # /usr/lib/python3.9/zipfile.py: 1362
    # 可能用cp437可能用utf-8,一个不行则试另一个
    utf8name = name.encode('cp437')
    utf8name = utf8name.decode('gbk')
    #    print("Extracting " + utf8name)
    pathname = os.path.dirname(utf8name)
    if not os.path.exists(pathname) and pathname != "":
        os.makedirs(pathname)
    data = file.read(name)
    if not os.path.exists(utf8name):
        fo = open(utf8name, "wb")
        fo.write(data)
        fo.close()
file.close()
