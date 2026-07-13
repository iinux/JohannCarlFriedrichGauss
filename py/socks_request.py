import requests

port = '1081'

proxies = {
    'http': 'socks5h://127.0.0.1:'+port,
    'https': 'socks5h://127.0.0.1:'+port,
}
#proxies = None

url = 'http://httpbin.org/ip'
url = 'http://ip.sb'
response = requests.get(url, proxies=proxies)

print(response.text)

