import requests
import shutil
import sys
import os
from urllib.parse import urlparse

def get_filename_from_url(url):
    """
    从URL中截取文件名
    """
    parsed_url = urlparse(url)
    return os.path.basename(parsed_url.path)

def download_file(url, proxy):
    """
    使用socks代理下载文件
    """
    proxies = {
        'http': proxy,
        'https': proxy,
    }

    # 发送请求并下载文件
    response = requests.get(url, proxies=proxies, stream=True)

    # 检查请求是否成功
    if response.status_code == 200:
        filename = get_filename_from_url(url)
        with open(filename, 'wb') as out_file:
            shutil.copyfileobj(response.raw, out_file)
        print(f"文件下载完成: {filename}")
    else:
        print(f"下载失败，状态码: {response.status_code}")

if __name__ == "__main__":
    if len(sys.argv) != 4:
        print("用法: python script.py <URL> <socks5代理地址> <socks5代理端口>")
        sys.exit(1)

    url = sys.argv[1]
    proxy_host = sys.argv[2]
    proxy_port = sys.argv[3]

    proxy = f'socks5h://{proxy_host}:{proxy_port}'
    download_file(url, proxy)

