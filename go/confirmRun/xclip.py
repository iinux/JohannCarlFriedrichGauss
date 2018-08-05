import pyperclip
f=open('/home/abui/www/xclip3')
c=f.read()
c=c.decode('gbk')
pyperclip.copy(c)
#print pyperclip.paste()
