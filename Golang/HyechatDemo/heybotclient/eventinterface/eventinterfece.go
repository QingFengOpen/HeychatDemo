package eventinterface

import (
	"context"

	"HyechatDemo/heybotclient/model"
)

type EventHandler interface {
	OnMessage(c context.Context, msg *model.GenericType) error
}
