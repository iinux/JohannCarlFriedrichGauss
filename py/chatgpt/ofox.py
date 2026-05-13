import my_config
from openai import OpenAI

client = OpenAI(
    api_key=my_config.ofox_key,
    base_url="https://api.ofox.ai/v1",
)

response = client.chat.completions.create(
    model="openai/gpt-5.4-mini",
    messages=[{"role": "user", "content": "生命的意义是什么？"}],
)
print(response.choices[0].message.content)
