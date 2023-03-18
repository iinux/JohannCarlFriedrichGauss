# -*- coding: utf-8 -*-
import os
import time
import redis
import openai
import my_config

from flask import Flask, request, abort, render_template
from wechatpy import parse_message, create_reply
from wechatpy.utils import check_signature
from wechatpy.exceptions import (
    InvalidSignatureException,
    InvalidAppIdException,
)

# set token or get from environments
TOKEN = os.getenv("WECHAT_TOKEN", my_config.wechat_token)
AES_KEY = os.getenv("WECHAT_AES_KEY", my_config.wechat_aes_key)
APPID = os.getenv("WECHAT_APPID", my_config.wechat_app_id)
THINKING = 'THINKING'

openai.api_key = my_config.key_openai


def chat_with_gpt3(msg, user_model):
    if user_model in ['gpt-3.5-turbo']:
        response = openai.ChatCompletion.create(
            model="gpt-3.5-turbo",
            messages=[{'role': 'user', 'content': msg.content}],
            max_tokens=2048,
            n=1,
            stop=None,
            temperature=.9,
            user=msg.source,
        )
        message = response.choices[0].message.content
    else:
        response = openai.Completion.create(
            engine=user_model,
            prompt=msg.content,
            max_tokens=2048,
            n=1,
            stop=None,
            temperature=0.5,
            user=msg.source,
        )
        message = response.choices[0].text
    print(response)

    return message


def set_user_model(user, model):
    r = redis.Redis(host='localhost', port=6379, decode_responses=True)
    r.set('user_model_' + user, model)
    r.close()


def get_user_model(user):
    r = redis.Redis(host='localhost', port=6379, decode_responses=True)
    user_model = r.get('user_model_' + user)
    r.close()
    if user_model is None:
        user_model = 'text-davinci-003'
    return user_model


def ask(msg, allow_cache=True):
    r = redis.Redis(host='localhost', port=6379, decode_responses=True)

    user_model = get_user_model(msg.source)

    key = 'chatgpt_answer_' + msg.content
    answer_cache = r.get(key)
    if answer_cache == THINKING and allow_cache:
        print('get THINKING')
        reply_content = 'LOOP'
        read_times = 10
        while read_times > 0:
            read_times -= 1
            time.sleep(1)
            answer_cache = r.get(key)
            if answer_cache is None:
                reply_content = "loop result fail, try again"
                break
            elif answer_cache != THINKING:
                print('GOT')
                reply_content = answer_cache
                break
        reply_content += '[L]'
    elif answer_cache and allow_cache:
        print('answer from cache')
        reply_content = answer_cache
        reply_content += '[C]'
    else:
        print('call openapi')
        r.set(key, THINKING, 300)
        try:
            time_start = time.time()
            reply_content = chat_with_gpt3(msg, user_model).lstrip('?？').strip()
            time_end = time.time()
            reply_content += '\n(time cost %.3f s)' % (time_end - time_start)
            r.set(key, reply_content, 28800)
        except openai.error.APIConnectionError as e:
            print(e)
            reply_content = "call openapi fail, try again " + e.__str__()
            r.delete(key)
        reply_content += '[I]'
    return reply_content


app = Flask(__name__)


@app.route("/")
def index():
    host = request.url_root
    return render_template("index.html", host=host)


@app.route("/wechat", methods=["GET", "POST"])
def wechat():
    signature = request.args.get("signature", "")
    timestamp = request.args.get("timestamp", "")
    nonce = request.args.get("nonce", "")
    encrypt_type = request.args.get("encrypt_type", "raw")
    msg_signature = request.args.get("msg_signature", "")
    try:
        check_signature(TOKEN, signature, timestamp, nonce)
    except InvalidSignatureException:
        abort(403)
    if request.method == "GET":
        echo_str = request.args.get("echostr", "")
        return echo_str

    # POST request
    if encrypt_type == "raw":
        # plaintext mode
        msg = parse_message(request.data)
        print(msg)
        if msg.type == "text":
            reply_content = ""
            if msg.content == "list":
                reply_content = ','.join([x.id for x in openai.Model.list().get('data')])
            elif msg.content.startswith("set "):
                user_model = msg.content[4:]
                set_user_model(msg.source, user_model)
                reply_content = "set user model success"
            elif msg.content.startswith('call_api '):
                reply_content = ask(msg.content[9:], allow_cache=False)
            else:
                reply_content = ask(msg)
            print('ask {} response {}'.format(msg.content, reply_content))
            reply = create_reply(reply_content, msg)
        else:
            reply = create_reply("Sorry, can not handle this for now", msg)
        return reply.render()
    else:
        # encryption mode
        from wechatpy.crypto import WeChatCrypto

        crypto = WeChatCrypto(TOKEN, AES_KEY, APPID)
        try:
            msg = crypto.decrypt_message(request.data, msg_signature, timestamp, nonce)
        except (InvalidSignatureException, InvalidAppIdException):
            abort(403)
        else:
            msg = parse_message(msg)
            if msg.type == "text":
                reply = create_reply(msg.content, msg)
            else:
                reply = create_reply("Sorry, can not handle this for now", msg)
            return crypto.encrypt_message(reply.render(), nonce, timestamp)


if __name__ == "__main__":
    app.run("127.0.0.1", 8081, debug=False)
