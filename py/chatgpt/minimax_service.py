import requests
import my_config

group_id = my_config.minimax_group_id
api_key = my_config.minimax_api_key

url = f"https://api.minimax.chat/v1/text/chatcompletion_pro?GroupId={group_id}"
headers = {"Authorization": f"Bearer {api_key}", "Content-Type": "application/json"}
model_id_1 = 'abab5.5-chat'


def get_default_model():
    return model_id_1


def get_model_list():
    return [model_id_1]


def ask_from_wechat(msg, user_model=None):
    return ask(msg.content, user_model)


def ask(content, user_model=None):
    if not user_model:
        user_model = get_default_model()
    # tokens_to_generate/bot_setting/reply_constraints 可自行修改
    request_body = {
        "model": user_model,
        "tokens_to_generate": 1024,
        "reply_constraints": {"sender_type": "BOT", "sender_name": "MM智能助理"},
        "messages": [
            {"sender_type": "USER", "sender_name": "小明", "text": content}
        ],
        "bot_setting": [
            {
                "bot_name": "MM智能助理",
                "content": "MM智能助理是一款由MiniMax自研的，没有调用其他产品的接口的大型语言模型。MiniMax是一家中国科技公司，一直致力于进行大模型相关的研究。",
            }
        ],
    }
    response = requests.post(url, headers=headers, json=request_body)
    # print(response.text)
    reply = response.json()["reply"]
    return reply
    # request_body["messages"].extend(response.json()["choices"][0]["messages"])


if __name__ == '__main__':
    t = ask('你好')
    print(t)
    pass
