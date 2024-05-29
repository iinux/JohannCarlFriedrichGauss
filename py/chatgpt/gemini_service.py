import requests
import my_config

proxies = {
  "http": "socks5://127.0.0.1:1081",
  "https": "socks5://127.0.0.1:1081",
}

url = 'https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent?key=' + my_config.gemini_key


def ask(msg):
    response = requests.post(url, headers={
        'Content-Type': 'application/json'
    }, json={
        'contents': [{
            'parts': [{
                'text': msg
            }]
        }]
    }, proxies=proxies)
    return response.json()['candidates'][0]['content']['parts'][0]['text']


if __name__ == '__main__':
    t = ask('你好')
    print(t)
    pass
