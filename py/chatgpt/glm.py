from openai import OpenAI
import my_config

client = OpenAI(
    base_url="https://open.bigmodel.cn/api/paas/v4",
    api_key=my_config.glm_key
)

completion = client.chat.completions.create(
    model="glm-4.7-flash",
    messages=[{"content": "who are you", "role": "user"}],
    temperature=1,
    top_p=0.95,
    max_tokens=8192,
    extra_body={"chat_template_kwargs": {"thinking": True}},
    stream=True
)

for chunk in completion:
    if not getattr(chunk, "choices", None):
        continue
    reasoning = getattr(chunk.choices[0].delta, "reasoning_content", None)
    if reasoning:
        print(reasoning, end="")
    if chunk.choices and chunk.choices[0].delta.content is not None:
        print(chunk.choices[0].delta.content, end="")


if __name__ == '__main__':
    pass
