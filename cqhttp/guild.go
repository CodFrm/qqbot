package cqhttp

import (
	"html"
	"strconv"
)

//map[channel_id:1.264896e+06 guild_id:3.5346341636528596e+16 message:123 message_id:272-451882151 message_type:guild post_type:message
//self_id:3.308923602e+09 self_tiny_id:1.4411521867741216e+17 sender:map[nickname:日落 user_id:1.4411521867809888e+17] sub_type:channel time:1.63724038e+09 user_id:1.4411521867809888e+17]

type GuildMsg struct {
	ChannelId  int64  `json:"channel_id"`
	GuildId    int64  `json:"guild_id"`
	Message_   string `json:"message"`
	MessageId  string `json:"message_id"`
	SelfId     int64  `json:"self_id"`
	SelfTinyId int64  `json:"self_tiny_id"`
	SubType    string `json:"sub_type"`
	Time_      int64  `json:"time"`
	UserId     int64  `json:"user_id"`
	Sender_    Sender `json:"sender"`
}

func (g *GuildMsg) Self() int64 {
	return g.SelfTinyId
}

func (g *GuildMsg) Message() string {
	return html.UnescapeString(g.Message_)
}
func (g *GuildMsg) Time() int64 {
	return g.Time_
}

func (g *GuildMsg) Sender() Sender {
	return g.Sender_
}

func (g *GuildMsg) MessageType() string {
	return "guild"
}

func (g *GuildMsg) ReplyText(s string) error {
	return SendGuildChannelMsg(g.GuildId, g.ChannelId, "[CQ:at,qq="+strconv.FormatInt(g.UserId, 10)+"]"+s)
}

func (g *GuildMsg) Group() int64 {
	return g.GuildId
}
