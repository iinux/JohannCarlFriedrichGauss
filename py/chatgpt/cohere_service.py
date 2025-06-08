import my_config
import cohere


def ask_from_wechat(msg, user_model=None):
    return ask(msg.content, user_model)


def ask(content, user_model=None):
    co = cohere.Client(my_config.cohere_key)

    response = co.chat(
        chat_history=[
            {"role": "USER", "message": "Who discovered gravity?"},
            {
                "role": "CHATBOT",
                "message": "The man who is widely credited with discovering gravity is Sir Isaac Newton",
            },
        ],
        message="What year was he born?",
        # perform web search before answering the question. You can also use your own custom connector.
        connectors=[{"id": "web-search"}],
    )

    print(response)


if __name__ == '__main__':
    t = ask('你好')
    print(t)
    pass
