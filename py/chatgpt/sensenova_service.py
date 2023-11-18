import my_config
import sensenova

# pip install sensenova


sensenova.access_key_id = my_config.sensenova_key_id
sensenova.secret_access_key = my_config.sensenova_secret

model_id_1 = 'nova-ptc-xl-v1'
model_id_2 = 'nova-ptc-xs-v1'
model_id_3 = 'nova-embedding-stable'
model_id_4 = 'nova-ptc-s-v2'
stream = True  # 流式输出或非流式输出


def get_default_model():
    return model_id_4


def get_model_list():
    return [model_id_1, model_id_2, model_id_3, model_id_4]


def print_info():
    # resp = sensenova.Model.list()
    resp = sensenova.Model.retrieve(id=get_default_model())
    print(resp)


def ask_from_wechat(msg, user_model=None):
    return ask(msg.content, user_model)


def ask(content, user_model=None):
    if not user_model:
        user_model = get_default_model()
    resp = sensenova.ChatCompletion.create(
        messages=[{"role": "user", "content": content}],
        model=user_model,
        stream=stream,
    )

    # print(resp)

    return_text = ''

    if not stream:
        resp = [resp]

    for part in resp:
        choices = part['data']["choices"]
        for c_idx, c in enumerate(choices):
            if len(choices) > 1:
                return_text += ("===== Chat Completion {} =====\n".format(c_idx))
            if stream:
                delta = c.get("delta")
                if delta:
                    return_text += delta
            else:
                return_text += (c["message"])
                if len(choices) > 1:
                    return_text += "\n"

    return return_text


if __name__ == '__main__':
    t = ask('你好')
    print(t)
    pass
