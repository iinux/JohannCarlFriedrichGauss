import my_config
import sensenova
import sys

# pip install sensenova


sensenova.access_key_id = my_config.sensenova_key_id
sensenova.secret_access_key = my_config.sensenova_secret

model_id = 'nova-ptc-xl-v1'
model_id = 'nova-ptc-xs-v1'
model_id = 'nova-embedding-stable'
model_id = 'nova-ptc-s-v2'
stream = True  # 流式输出或非流式输出


def print_info():
    # resp = sensenova.Model.list()
    resp = sensenova.Model.retrieve(id=model_id)
    print(resp)


def ask(content):
    resp = sensenova.ChatCompletion.create(
        messages=[{"role": "user", "content": content}],
        model=model_id,
        stream=stream,
    )

    print(resp)

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
