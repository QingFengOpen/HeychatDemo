import asyncio
import json
import logging
import ssl
from typing import Optional

import websockets
from websockets.exceptions import ConnectionClosedError, ConnectionClosedOK

import conf.config
from eventhandler.eventhandler import EventHandler


# Assuming these are defined elsewhere or as part of your project
class EventType:
    PING = "PING"
    PONG = "PONG"
    TAG_PONG = 'PONG, tag: []'


class GenericType:
    def __init__(self, data):
        self.data = data


# Constants
PING_INTERVAL = 30
CHECK_INTERVAL = 30
SLEEP_TIME = 1
MAX_SLEEP_TIME = 60


def get_headers():
    return {
        "Accept": "application/json, text/plain, */*",
        "Accept-Language": "zh-CN,zh;q=0.9",
        "Cache-Control": "no-cache",
        "Pragma": "no-cache",
    }


def get_wss_url(token):
    """
    根据当前配置和令牌构建WebSocket连接URL.
    """
    base_url = conf.config.WSS_URL
    common_params = conf.config.COMMON_PARAMS
    token_params = conf.config.TOKEN_PARAMS
    return f"{base_url}{common_params}{token_params}{token}"


class WebSocketClient:
    def __init__(self, token: str, event_handler: Optional[EventHandler] = None):
        self.token = token
        self.event_handler = event_handler if event_handler is not None else EventHandler()
        self.conn = None
        self.is_connected = False
        self.close = False
        self.message_queue = asyncio.Queue()
        self.ctx = None
        self.ping_task = None
        self.logger = logging.getLogger(__name__)

    async def connect(self):
        uri = get_wss_url(self.token)
        context = ssl.create_default_context()
        context.check_hostname = False
        context.verify_mode = ssl.CERT_NONE
        try:
            self.conn = await websockets.connect(uri, ssl=context, extra_headers=get_headers())
            self.is_connected = True
            self.logger.info("Connected to the server.")
            self.ctx = asyncio.get_running_loop()
            self.ping_task = self.ctx.create_task(self.heartbeat())
            self.ctx.create_task(self.receive())
            self.ctx.create_task(self.handle())

        except Exception as e:
            self.logger.error(f"Connection failed: {e}")
            await self.reconnect()

    async def send_ping(self):
        await self.conn.send(EventType.PING)

    async def heartbeat(self):
        self.logger.info("start heartbeat")
        while not self.close:
            await asyncio.sleep(PING_INTERVAL)
            if not self.is_connected:
                continue
            try:
                await self.send_ping()
            except ConnectionClosedError:
                self.logger.error("Connection closed unexpectedly during heartbeat.")
                await self.reconnect()
                break

    async def receive(self):
        self.logger.info("receive data")
        while not self.close:
            try:
                message = await self.conn.recv()
                await self.message_queue.put(message)
            except (ConnectionClosedOK, ConnectionClosedError):
                self.logger.error("Connection closed.")
                await self.reconnect()
                break

    async def reconnect(self):
        self.is_connected = False
        sleep_time = SLEEP_TIME
        while not self.close and not self.is_connected:
            await asyncio.sleep(sleep_time)
            if sleep_time < MAX_SLEEP_TIME:
                sleep_time *= 2
            else:
                sleep_time = SLEEP_TIME
            self.logger.info("Reconnecting to the server...")
            await self.connect()

    async def handle(self):
        self.logger.info("handle data")
        while not self.close:
            message = await self.message_queue.get()
            if message == EventType.PONG or message == EventType.PONG.lower() or str.startswith(message,
                                                                                                EventType.PONG):
                continue
            try:
                data = json.loads(message)
                await self.event_handler.on_message(data)
            except json.JSONDecodeError as e:
                self.logger.error(f"Failed to decode JSON: {e}")

    async def close_client(self):
        self.close = True
        if self.conn:
            await self.conn.close()
        if self.ctx:
            self.ctx.stop()
        if self.ping_task:
            self.ping_task.cancel()
            try:
                await self.ping_task
            except asyncio.CancelledError:
                pass


# Example usage
async def start(token):
    client = WebSocketClient(token)
    await client.connect()
    try:
        while True:
            await asyncio.sleep(1)
    except KeyboardInterrupt:
        await client.close_client()
