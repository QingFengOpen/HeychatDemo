package model

import (
	"encoding/json"
)

type GenericType struct {
	Type       string          `json:"type"`
	Sequence   int64           `json:"sequence"`
	Data       json.RawMessage `json:"data"`
	Timestamp  int64           `json:"timestamp"`
	NotifyType string          `json:"notify_type"`
}

const (
	MsgTypeCommand = "50"
)

type RoomBaseInfo struct {
	RoomAvatar string `json:"room_avatar"`
	RoomId     string `json:"room_id"`
	RoomName   string `json:"room_name"`
}

type SenderInfo struct {
	Avatar           string `json:"avatar"`
	AvatarDecoration struct {
		SrcType string `json:"src_type"`
		SrcUrl  string `json:"src_url"`
	} `json:"avatar_decoration"`
	Bot          bool     `json:"bot"`
	Level        int      `json:"level"`
	Medals       any      `json:"medals"`
	Nickname     string   `json:"nickname"`
	Roles        []string `json:"roles"`
	RoomNickname string   `json:"room_nickname"`
	Tag          any      `json:"tag"`
	UserId       int      `json:"user_id"`
}

type UseCommandData struct {
	BotId           int              `json:"bot_id"`
	ChannelBaseInfo *ChannelBaseInfo `json:"channel_base_info"`
	CommandInfo     *CommandInfo     `json:"command_info"`
	MsgId           string           `json:"msg_id"`
	RoomBaseInfo    *RoomBaseInfo    `json:"room_base_info"`
	SendTime        int64            `json:"send_time"`
	SenderInfo      *SenderInfo      `json:"sender_info"`
}
