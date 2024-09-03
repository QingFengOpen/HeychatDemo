import common.common
import conf.command as command
import conf.model as model

RepeaterCommandID = "1830851200470908928"


def on_use_bot_command(data):
    meta = {}
    if isinstance(data, dict):
        meta = model.UseCommandData(**data)
    if meta and meta.command_info:
        command_id = meta.command_info.id
        if command_id == RepeaterCommandID:
            on_repeater(meta)


def on_repeater(meta):
    if len(meta.command_info.options) == 1:
        option = meta.command_info.options[0]
        if option.type == command.TYPE_STRING:
            req = model.ChannelImSendReq(
                msg=option.value,
                msg_type=model.MSG_TYPE_MDTEXT,
                channel_id=meta.channel_base_info.channel_id,
                room_id=meta.room_base_info.room_id,
            )
            common.common.SendMessage(req)


class EventHandler:
    async def on_message(self, data):
        message_type = data["type"]
        message_data = data["data"]
        if message_type == model.MSG_TYPE_USECOMMAND:
            on_use_bot_command(message_data)
