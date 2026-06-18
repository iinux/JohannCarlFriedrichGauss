import browser_cookie3
import json

# 导出 Chrome Cookie
cookies = browser_cookie3.chrome(domain_name='baidu.com')
cookie_dict = [c.__dict__ for c in cookies]

with open('cookies.json', 'w') as f:
    json.dump(cookie_dict, f)

# 导入 Cookie
#import requests
#session = requests.Session()
#for cookie in cookie_dict:
#    session.cookies.set(**cookie)

if __name__ == '__main__':
    pass