package model

const (
	TypeSubCommand      int = 1
	TypeSubCommandGroup int = 2
	TypeString          int = 3
	TypeInteger         int = 4
	TypeBoolean         int = 5
	TypeUser            int = 6
	TypeChannel         int = 7
	TypeRole            int = 8
)

type ChannelBaseInfo struct {
	ChannelId   string `json:"channel_id"`
	ChannelName string `json:"channel_name"`
	ChannelType int    `json:"channel_type"`
}

type CommandInfo struct {
	Id      string     `json:"id"`
	Name    string     `json:"name"`
	Type    int        `json:"type"`
	Options []*Options `json:"options,omitempty" bson:"options"`
}

type Options struct {
	Name    string     `form:"name" json:"name" bson:"name" `
	Type    int        `form:"type" json:"type" bson:"type"`
	Value   string     `form:"value" json:"value,omitempty" bson:"value"`
	Choices []*Options `form:"choices" json:"choices,omitempty" bson:"choices"`
}
