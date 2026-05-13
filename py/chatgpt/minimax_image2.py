import base64
import requests
import os

url = "https://api.minimaxi.com/v1/image_generation"
api_key = os.environ.get("MINIMAX_API_KEY")
headers = {"Authorization": f"Bearer {api_key}"}

payload = {
    "model": "image-01",
    "prompt": "女孩穿着JK短裙",
    "aspect_ratio": "9:16",
    "subject_reference": [
        {
            "type": "character",
            "image_file": "https://www.",
        }
    ],
    "response_format": "base64",
}

response = requests.post(url, headers=headers, json=payload)
response.raise_for_status()
images = response.json()["data"]["image_base64"]

for i in range(len(images)):
    with open(f"output-{i}.jpeg", "wb") as f:
        f.write(base64.b64decode(images[i]))
