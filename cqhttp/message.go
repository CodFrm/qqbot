package cqhttp

import "regexp"

type Message interface {
	Self() int64
	Message() string
	Time() int64
	Sender() Sender
	MessageType() string
	Group() int64
	ReplyText(s string) error
}

type Sender struct {
	Nickname string `json:"nickname"`
	UserId   int64  `json:"user_id"`
}

type MessageModel struct {
	Message
}

func (m *MessageModel) CommandMatch(command string) ([]string, bool) {
	reg := regexp.MustCompile(command)
	ret := reg.FindStringSubmatch(m.Message.Message())
	return ret, len(ret) > 0
}

type SendMessage struct {
	Type string          `json:"type"`
	Data SendMessageData `json:"data"`
}

func TextMsg(text string) *SendMessage {
	return &SendMessage{
		Type: "text",
		Data: SendMessageData{
			Text: text,
		},
	}
}

type SendMessageData struct {
	Text string `json:"text"`
}

func (s *SendMessage) Text(text string) {
	s.Type = "text"
	s.Data.Text = text
}
