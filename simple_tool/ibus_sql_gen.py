# coding=utf8

from __future__ import print_function

m='iiis'
c='沙河'
clen=2;

mlen=len(m)
mArr=[]
for item in m:
    mArr.append(ord(item)-96)

sql="insert into phrases(mlen,clen,m0,m1,m2,m3,category,phrase,freq,user_freq) values(%d,%d,%d,%d,%d,%d,1,'%s',10700000, 0)" % (mlen, clen, mArr[0], mArr[1], mArr[2], mArr[3], c)
print(sql)
cmd="sudo sqlite3 /usr/share/ibus-table/tables/wubi-jidian86.db \"%s\"" % sql
print(cmd)
