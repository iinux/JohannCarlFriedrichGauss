from openai import OpenAI
import my_config

client = OpenAI(
    base_url="https://integrate.api.nvidia.com/v1",
    api_key=my_config.nvidia_key
)

completion = client.chat.completions.create(
    model="deepseek-ai/deepseek-v3.2",
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


