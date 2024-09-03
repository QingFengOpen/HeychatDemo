# This is a sample Python script.
import asyncio
import logging

from client import start
from conf import config

# Press ⇧F10 to execute it or replace it with your code.
# Press Double ⇧ to search everywhere for classes, files, tool windows, actions, and settings.


if __name__ == "__main__":
    token = config.HeyChatAPPToken
    logging.basicConfig(level=logging.INFO)
    asyncio.run(start(token))
