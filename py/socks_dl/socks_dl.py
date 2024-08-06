import requests
import shutil

# pip install 'requests[socks]'
# pip install tqdm


# 配置socks代理
proxies = {
    #'http': 'socks5h://username:password@proxy_host:proxy_port',
    #'https': 'socks5h://username:password@proxy_host:proxy_port',
    'http': 'socks5h://127.0.0.1:1080',
    'https': 'socks5h://127.0.0.1:1080',
}

# 文件下载的URL
url = 'https://github.com/alist-org/alist/releases/download/v3.36.0/alist-linux-amd64.tar.gz'

# 发送请求并下载文件
response = requests.get(url, proxies=proxies, stream=True)

# 保存文件到本地
with open('downloaded_file.zip', 'wb') as out_file:
    shutil.copyfileobj(response.raw, out_file)

print("文件下载完成")

