package eventhandler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"HyechatDemo/heybotclient/common"
	"HyechatDemo/heybotclient/model"
)

type Domain struct {
	httpClient  *http.Client
	clientToken string
}

const (
	RepeaterCommandID = "9527"
)

func New(token string) *Domain {
	domain := &Domain{
		httpClient:  &http.Client{},
		clientToken: token,
	}
	return domain
}

func (d *Domain) OnMessage(c context.Context, msg *model.GenericType) error {
	switch msg.Type {
	case model.MsgTypeCommand:
		return d.OnUseBotCommand(c, msg.Data)
	default:
	}
	return nil
}

func (d *Domain) OnUseBotCommand(_ context.Context, data json.RawMessage) error {
	command := &model.UseCommandData{}
	if err := json.Unmarshal(data, command); err != nil {
		return err
	}
	if command != nil && command.CommandInfo != nil {
		commandID := command.CommandInfo.Id
		if commandID == RepeaterCommandID {
			if len(command.CommandInfo.Options) == 1 {
				option := command.CommandInfo.Options[0]
				if option.Type == model.TypeString {
					return common.SendChannelIM(d.httpClient, &model.ChannelImSendReq{
						AtUserId:     strconv.Itoa(command.SenderInfo.UserId),
						Msg:          option.Value,
						ChannelId:    command.ChannelBaseInfo.ChannelId,
						MsgType:      10,
						ChannelType:  1,
						RoomId:       command.RoomBaseInfo.RoomId,
						Addition:     "{\"img_files_info\":[]}",
						HeychatAckId: "",
					}, d.clientToken)
				}
			}
		}
	}
	return nil
}
