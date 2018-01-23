# coding=utf8
import sqlite3
from sys import argv

conn = sqlite3.connect('a')
c = conn.cursor()

# Create table
c.execute('''CREATE TABLE stocks (date text, trans text, symbol text, qty real, price real)''')

# Insert a row of data
c.execute("INSERT INTO stocks VALUES ('2006-01-05','BUY','RHAT',100,35.14)")
# c.execute(u"insert into phrases(mlen,clen,m0,m1,m2,m3,category,phrase,freq) values(4,2,20,13,23,25,1,'微信',10700000)")

# Save (commit) the changes
conn.commit()

# We can also close the connection if we are done with it.
# Just be sure any changes have been committed or they will be lost.
conn.close()
