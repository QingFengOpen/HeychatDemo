package model

import "fmt"

type ChannelImSendReq struct {
	Msg              string `json:"msg"`
	MsgType          int    `json:"msg_type"`
	HeychatAckId     string `json:"heychat_ack_id"`
	ReplyId          string `json:"reply_id"`
	RoomId           string `json:"room_id"`
	Addition         string `json:"addition"`
	AtUserId         string `json:"at_user_id"`
	AtRoleId         string `json:"at_role_id"`
	MentionChannelId string `json:"mention_channel_id"`
	ChannelId        string `json:"channel_id"`
	ChannelType      int    `json:"channel_type"`
}

func GetWssUrl(token string) string {
	return fmt.Sprintf("%s/chatroom/ws/connect?chat_os_type=bot&client_type=heybox_chat&chat_version=999.0.0&token=%s", getHost(), token)
}

func getHost() string {
	return "wss://chat.xiaoheihe.cn"
}

const HttpReqHost = "https://chat.xiaoheihe.cn"
