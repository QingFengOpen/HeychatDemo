import json

import conf.config as conf
import requests


def SendMessage(payload):
    url = f"{conf.HTTP_HOST}{conf.SEND_MSG_URL}{conf.COMMON_PARAMS}"

    headers = {
        'Content-Type': 'application/json;charset=UTF-8',
        'token': conf.HeyChatAPPToken,
    }
    response = requests.request("POST", url, headers=headers, data=payload.json())
    content = response.content.decode('utf-8')
    print(content)





