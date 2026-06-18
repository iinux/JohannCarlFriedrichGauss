#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
OpenAI API 代理网关
功能：代理 OpenAI 请求，并修改 api_key 和 model
"""

import json
import re
import time
from http.server import HTTPServer, BaseHTTPRequestHandler
from urllib.request import Request, urlopen
from urllib.error import HTTPError
import ssl
import os
try:
    import my_config
    _DEFAULT_API_KEY = getattr(my_config, 'ali_key', None)
except ImportError:
    my_config = None
    _DEFAULT_API_KEY = None

# ============ 配置 ============
# 目标 OpenAI API 地址
TARGET_HOST = "dashscope.aliyuncs.com"
TARGET_URL = f"https://{TARGET_HOST}/compatible-mode/v1"

# 替换配置
REPLACEMENTS = {
    # 替换 API Key: {原始key或正则: 新key}
    "api_keys": {
        # 示例：替换特定 key
        "sk-original123": "sk-newApiKey456",
        # 替换所有 sk- 开头的 key，运行时由命令行参数填充
        r"sk-[a-zA-Z0-9]+": _DEFAULT_API_KEY,
    },

    # 替换 Model: {原始model: 新model}，运行时由命令行参数填充
    "models": {}
}

# 网关监听配置
GATEWAY_HOST = "0.0.0.0"
GATEWAY_PORT = 8081
PRINT_REQUEST = True
PRINT_RESPONSE = False


class OpenAIGatewayHandler(BaseHTTPRequestHandler):
    """OpenAI API 网关处理器"""

    def log_request(self, code='-', size='-'):
        pass  # 禁用 send_response 触发的自动日志，改为读完 body 后手动打印

    def log_message(self, format, *args):
        elapsed = getattr(self, '_elapsed', None)
        suffix = f" {elapsed:.2f}s" if elapsed is not None else ""
        print(f"[{self.log_date_time_string()}] {args[0]}{suffix}")

    def do_GET(self):
        """处理 GET 请求"""
        self._proxy_request("GET")

    def do_POST(self):
        """处理 POST 请求"""
        self._proxy_request("POST")

    def do_PUT(self):
        """处理 PUT 请求"""
        self._proxy_request("PUT")

    def do_DELETE(self):
        """处理 DELETE 请求"""
        self._proxy_request("DELETE")

    def do_OPTIONS(self):
        """处理 OPTIONS 请求（CORS 预检）"""
        self.send_response(200)
        self._send_cors_headers()
        self.end_headers()

    def _send_cors_headers(self):
        """发送 CORS 响应头"""
        self.send_header('Access-Control-Allow-Origin', '*')
        self.send_header('Access-Control-Allow-Methods', 'GET, POST, PUT, DELETE, OPTIONS')
        self.send_header('Access-Control-Allow-Headers', 'Content-Type, Authorization')

    def _modify_request_body(self, body: bytes) -> bytes:
        """
        修改请求体中的 api_key 和 model
        """
        try:
            data = json.loads(body.decode('utf-8'))
            if PRINT_REQUEST:
                print(f'[Request]{data}')

            # 替换 model
            if 'model' in data:
                original_model = data['model']
                data['model'] = self._replace_value(
                    original_model,
                    REPLACEMENTS['models']
                )
                if original_model != data['model']:
                    print(f"[Model] {original_model} -> {data['model']}")

            return json.dumps(data).encode('utf-8')
        except:
            return body

    def _modify_request_headers(self, headers: dict) -> dict:
        """
        修改请求头中的 Authorization (api_key)
        """
        key = 'authorization'
        auth_header = headers.get(key, '')
        if auth_header == '':
            key = 'Authorization'
            auth_header = headers.get(key, '')

        if auth_header.startswith('Bearer '):
            original_key = auth_header[7:]  # 去掉 "Bearer "
            new_key = self._replace_value(
                original_key,
                REPLACEMENTS['api_keys']
            )
            if original_key != new_key:
                # print(f"[API Key] {original_key[:20]}... -> {new_key[:20]}...")
                print(f"[API Key] {original_key} -> {new_key}")
                headers[key] = f'Bearer {new_key}'

        return headers

    def _replace_value(self, original: str, replacements: dict) -> str:
        """
        根据配置替换值（支持正则匹配）
        """
        for pattern, replacement in replacements.items():
            # 检查是否是正则表达式
            if pattern.startswith('r"') or pattern.startswith("r'"):
                # 去掉 r" 和 " 包装
                pattern = pattern[2:-1]

            try:
                # 尝试作为正则匹配
                if re.match(pattern, original):
                    return replacement
            except re.error:
                # 不是有效的正则，作为普通字符串匹配
                if pattern == original:
                    return replacement

        return original

    def _proxy_request(self, method: str):
        """
        代理请求到 OpenAI API
        """
        # 读取请求体
        key = 'content-length'
        content_length = int(self.headers.get(key, 0))
        if content_length == 0:
            key = 'Content-Length'
            content_length = int(self.headers.get(key, 0))
        body = self.rfile.read(content_length) if content_length > 0 else b''

        # 修改请求
        modified_body = self._modify_request_body(body)

        # 准备请求头
        headers = dict(self.headers)
        headers = self._modify_request_headers(headers)

        # 修改 Host 头
        headers['Host'] = TARGET_HOST

        if 'content-length' in headers:
            del headers['content-length']
        if 'Content-Length' in headers:
            del headers['Content-Length']

        # 移除 hop-by-hop 头
        hop_by_hop = ['connection', 'keep-alive', 'proxy-authenticate', 'proxy-authorization', 'te', 'trailers',
                      'transfer-encoding', 'upgrade']
        for header in hop_by_hop:
            headers.pop(header, None)
            headers.pop(header.title(), None)

        # 构建目标 URL
        target_url = f"{TARGET_URL}{self.path}"

        self._start = time.time()
        print(f"[{method}] {self.path} -> {target_url}")

        try:
            # 创建请求
            req = Request(
                url=target_url,
                data=modified_body if method in ['POST', 'PUT'] else None,
                headers=headers,
                method=method
            )

            # 发送请求（忽略 SSL 验证，生产环境请移除）
            ctx = ssl.create_default_context()
            ctx.check_hostname = False
            ctx.verify_mode = ssl.CERT_NONE

            with urlopen(req, context=ctx, timeout=120) as response:
                # 发送响应状态
                self.send_response(response.status)

                # 复制响应头
                for header, value in response.headers.items():
                    if header.lower() not in hop_by_hop:
                        self.send_header(header, value)

                self._send_cors_headers()
                self.end_headers()

                # 复制响应体
                rr = response.read()
                self._elapsed = time.time() - self._start
                if PRINT_RESPONSE:
                    print(f'[Response] {rr}')
                self.log_message('"%s" %s', self.requestline, str(response.status))
                print()
                self.wfile.write(rr)

        except HTTPError as e:
            self.send_response(e.code)
            self._send_cors_headers()
            self.end_headers()
            er = e.read()
            print(f'[Error] {e.code}: {e.reason} 错误响应：{er}')
            self.wfile.write(er)

        except Exception as e:
            print(f"[Error] {e}")
            self.send_response(500)
            self._send_cors_headers()
            self.end_headers()
            self.wfile.write(json.dumps({
                "error": {
                    "message": str(e),
                    "type": "gateway_error"
                }
            }).encode('utf-8'))


def select_model(target_url: str, api_key: str) -> str:
    ctx = ssl.create_default_context()
    ctx.check_hostname = False
    ctx.verify_mode = ssl.CERT_NONE
    req = Request(f"{target_url}/models", headers={"Authorization": f"Bearer {api_key}"})
    with urlopen(req, context=ctx, timeout=10) as resp:
        data = json.loads(resp.read())
    models = sorted(m["id"] for m in data.get("data", []))
    if not models:
        raise RuntimeError("未获取到可用模型")
    for i, m in enumerate(models, 1):
        print(f"  {i:>3}. {m}")
    while True:
        choice = input(f"请选择模型 [1-{len(models)}]: ").strip()
        if choice.isdigit() and 1 <= int(choice) <= len(models):
            return models[int(choice) - 1]
        print("输入无效，请重新选择")


def run_gateway(host_info):
    """
    启动网关服务器
    """
    server = HTTPServer(host_info, OpenAIGatewayHandler)
    print(f"=" * 60)
    print(f"OpenAI API 网关已启动")
    print(f"监听地址: {host_info}")
    print(f"目标地址: {TARGET_URL}")
    print(f"=" * 60)
    print()

    try:
        server.serve_forever()
    except KeyboardInterrupt:
        print("\n网关已停止")
        server.shutdown()


# ========== 使用示例 ==========

if __name__ == "__main__":
    import argparse

    parser = argparse.ArgumentParser(description="OpenAI API 代理网关")
    parser.add_argument("model", nargs="?", default=None, help="转发给上游的 model 名称，省略则从上游列表选择")
    parser.add_argument("--port", type=int, default=int(os.getenv("GATEWAY_PORT", GATEWAY_PORT)), help="监听端口（默认 8081）")
    parser.add_argument("--api-key", default=None, help="上游 API Key，默认取 my_config.ali_key")
    parser.add_argument("--target", default='ali', help="目标 API 主机")
    parser.add_argument("--print-response", action="store_true", help="打印响应体到终端（默认关闭）")
    parser.add_argument("--hide-request", action="store_false", help="打印请求到终端（默认打开）")
    parser.add_argument("--list-model", action="store_true", help="指定地址的支持模型列表")
    args = parser.parse_args()

    if args.list_model:
        select_model(args.target, args.api_key)
        sys.exit(-1)

    if args.target == 'glm':
        TARGET_HOST = "open.bigmodel.cn"
        TARGET_URL = f"https://{TARGET_HOST}/api/paas/v4"
    elif args.target == 'nvidia':
        TARGET_HOST = "integrate.api.nvidia.com"
        TARGET_URL = f"https://{TARGET_HOST}/v1"
    elif args.target == 'minimax':
        TARGET_HOST = 'api.minimaxi.com'
        TARGET_URL = f"https://{TARGET_HOST}/v1"
    elif args.target == 'deepseek':
        TARGET_HOST = 'api.deepseek.com'
        TARGET_URL = f'https://{TARGET_HOST}'


    api_key = args.api_key or _DEFAULT_API_KEY
    if not api_key:
        parser.error("未找到 API Key：请通过 --api-key 传入，或在 my_config.py 中设置 ali_key")
    REPLACEMENTS["api_keys"][r"sk-[a-zA-Z0-9]+"] = api_key

    model = args.model
    if model is None:
        model = select_model(TARGET_URL, api_key)
    print(f"使用模型: {model}")
    REPLACEMENTS["models"][r".*"] = model
    PRINT_REQUEST = args.hide_request
    PRINT_RESPONSE = args.print_response

    run_gateway((GATEWAY_HOST, args.port))
