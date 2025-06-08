import http.server
import socketserver
import requests
import argparse
from urllib.parse import urljoin, urlparse

# 配置参数解析
parser = argparse.ArgumentParser(description='HTTP to HTTPS Request Forwarder with Redirect Following')
parser.add_argument('--port', type=int, default=8080, help='HTTP listening port')
parser.add_argument('--target', type=str, required=True,
                    help='Target HTTPS base URL (e.g., https://example.com)')
parser.add_argument('--disable-verify', action='store_true',
                    help='Disable SSL certificate verification')
parser.add_argument('--max-redirects', type=int, default=5,
                    help='Maximum redirects to follow (default: 5)')
args = parser.parse_args()

class HTTPS_Forwarder(http.server.BaseHTTPRequestHandler):
    def handle_request(self, method):
        # 处理原始请求URL
        query_string = ""
        if '?' in self.path:
            path, query_string = self.path.split('?', 1)
            query_string = f"?{query_string}"
        else:
            path = self.path

        target_url = urljoin(args.target.rstrip('/') + "/", path.lstrip('/')) + query_string

        # 复制原始请求头
        headers = {k: v for k, v in self.headers.items()}
        # 移除可能需要修改的头
        headers.pop('Host', None)

        # 处理请求体
        content_length = int(self.headers.get('Content-Length', 0))
        request_body = self.rfile.read(content_length) if content_length > 0 else None

        # 发起请求并处理重定向
        try:
            redirect_count = 0
            while redirect_count <= args.max_redirects:
                # 发起请求
                verify = not args.disable_verify
                response = requests.request(
                    method,
                    target_url,
                    headers=headers,
                    data=request_body,
                    allow_redirects=False,  # 手动处理重定向
                    verify=verify,
                    stream=True,  # 流式处理大响应
                    timeout=30
                )

                # 如果响应不是重定向，返回结果
                if not 300 <= response.status_code < 400:
                    break

                # 处理重定向
                redirect_url = response.headers.get('Location', '')
                if not redirect_url:
                    self.send_error(500, "Redirect response missing Location header")
                    return

                # 处理相对URL
                if not urlparse(redirect_url).scheme:
                    redirect_url = urljoin(target_url, redirect_url)

                # 更新为新的请求URL
                target_url = redirect_url
                redirect_count += 1

                # 处理307/308重定向时保留方法
                if response.status_code in (307, 308) and method in ('POST', 'PUT', 'PATCH', 'DELETE'):
                    # 方法不变，保留请求体
                    continue

                # 对于301/302重定向，POST改为GET，清空请求体
                if method in ('POST', 'PUT', 'PATCH', 'DELETE'):
                    method = 'GET'
                    request_body = None
                    # 移除可能不再相关的请求头
                    headers.pop('Content-Type', None)
                    headers.pop('Content-Length', None)

            # 检查重定向次数是否超出限制
            if redirect_count > args.max_redirects:
                self.send_error(500, f"Exceeded maximum redirect limit of {args.max_redirects}")
                return

            # 发送最终响应
            self.send_response(response.status_code)

            # 复制响应头，过滤掉不安全的头
            excluded_headers = ['connection', 'content-length', 'transfer-encoding']
            for key, value in response.headers.items():
                if key.lower() not in excluded_headers:
                    self.send_header(key, value)

            # 特别处理内容长度
            self.send_header('Content-Length', str(len(response.content)))
            self.end_headers()

            # 发送内容
            self.wfile.write(response.content)

        except requests.exceptions.Timeout:
            self.send_error(504, "Gateway Timeout")
        except requests.exceptions.RequestException as e:
            self.send_error(502, f"Bad Gateway: {str(e)}")
        except Exception as e:
            self.send_error(500, f"Server Error: {str(e)}")

    # 支持所有 HTTP 方法
    def do_GET(self):
        self.handle_request('GET')

    def do_POST(self):
        self.handle_request('POST')

    def do_PUT(self):
        self.handle_request('PUT')

    def do_DELETE(self):
        self.handle_request('DELETE')

    def do_HEAD(self):
        self.handle_request('HEAD')

    def do_OPTIONS(self):
        self.handle_request('OPTIONS')

    def do_PATCH(self):
        self.handle_request('PATCH')

    def log_message(self, format, *args):
        # 禁用内置日志
        return

# 启动服务器
with socketserver.TCPServer(("", args.port), HTTPS_Forwarder) as httpd:
    print(f"Forwarding HTTP requests from http://localhost:{args.port} to {args.target}")
    print(f"Following up to {args.max_redirects} redirects")
    if args.disable_verify:
        print("WARNING: SSL certificate verification is disabled")

    try:
        httpd.serve_forever()
    except KeyboardInterrupt:
        print("\nServer stopped")
