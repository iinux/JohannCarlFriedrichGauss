import requests
import shutil
import sys
import os
from urllib.parse import urlparse
from tqdm import tqdm
import argparse

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
        total_size = int(response.headers.get('content-length', 0))
        filename = get_filename_from_url(url)
        
        with open(filename, 'wb') as out_file, tqdm(
            total=total_size, unit='B', unit_scale=True, desc=filename
        ) as progress_bar:
            for chunk in response.iter_content(chunk_size=1024):
                if chunk:
                    out_file.write(chunk)
                    progress_bar.update(len(chunk))

        print(f"\n文件下载完成: {filename}")
    else:
        print(f"下载失败，状态码: {response.status_code}")

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="使用socks代理下载文件")
    parser.add_argument('url', type=str, help='文件的URL')
    parser.add_argument('--proxy_host', type=str, default='127.0.0.1', help='socks5代理地址 (默认: 127.0.0.1)')
    parser.add_argument('--proxy_port', type=int, default=1080, help='socks5代理端口 (默认: 1080)')

    args = parser.parse_args()

    proxy = f'socks5h://{args.proxy_host}:{args.proxy_port}'
    download_file(args.url, proxy)

