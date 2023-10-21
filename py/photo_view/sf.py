import requests
import re
import os
import sf_data
from bs4 import BeautifulSoup

index = "/?ch=smov&op=free_movie&class_name=%E5%85%8D%E8%B4%B9%E4%B8%93%E5%8C%BA"
ua = 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) ' \
     'Chrome/118.0.0.0 Safari/537.36'
domain = sf_data.domain
session = sf_data.session


def request_index():
    index_content = requests.get(domain + index)
    print('index_content: ' + index_content.text)

    soup = BeautifulSoup(index_content.text, "html.parser")
    titles = soup.find_all("span", {"class": "icon-caret-right"})
    dates = soup.find_all("span", {"class": "icon-upload-alt"})
    codes = soup.find_all("span", {"class": "icon-barcode"})

    i = 0
    for title in titles:
        print(title.text)
        print(dates[i].text)
        print(str(dates[i]))
        print(codes[i].text)
        ds = re.findall(r"(\d{4}-\d{2}-\d{2})", str(dates[i]))
        file_name = "upload-%s-free-%s-code-%s-title-%s.mp4" % (ds[0], ds[1], codes[i].text, title.text)
        print(file_name)
        request_j_index(file_name, codes[i].text)
        i += 1


def request_j_index(file_name, code):
    res = requests.post(domain + '/jindex.php', headers={
        'cookie': 'PHPSESSID=%s' % session,
        'User-Agent': ua,
    }, data={
        "op": "do_playts",
        "func": "new_play",
        "post_id": session,
        "SPCode": code
    })
    print(res.text)

    video_url = res.json()['video_url']
    print(video_url)
    soup = BeautifulSoup(video_url, "html.parser")
    iframe = soup.find("iframe")
    video_url = iframe["src"]
    print(video_url)
    request_mpd(video_url, file_name)


def request_mpd(video_url, file_name):
    res = requests.get(domain + '/' + video_url, headers={
        'cookie': 'PHPSESSID=%s' % session,
        'User-Agent': ua,
    })
    print(res.text)
    mr = re.findall(r"https.*mpd", res.text)
    print(mr[0])
    print(file_name)
    download_mpd(mr[0], file_name)


def download_mpd(mpd, file_name):
    if not os.path.exists(file_name):
        cmd = "ffmpeg -i '%s' -c copy %s" % (mpd, file_name)
        print(cmd)
        os.system(cmd)


request_index()
