package cqhttp

//map[font:0 message:123 message_id:-1.883559911e+09 message_type:private post_type:message raw_message:123 self_id:3.308923602e+09
//sender:map[age:0 nickname:mhsj sex:unknown user_id:9.58139621e+08] sub_type:friend target_id:3.308923602e+09 time:1.637241196e+09 user_id:9.58139621e+08]

type PrivateMsg struct {
	Message_  string `json:"message"`
	MessageId int64  `json:"message_id"`
	SelfId    int64  `json:"self_id"`
	SubType   string `json:"sub_type"`
	Time_     int64  `json:"time"`
	UserId    int64  `json:"user_id"`
	Sender_   Sender `json:"sender"`
}

func (p *PrivateMsg) Self() int64 {
	return p.SelfId
}

func (p *PrivateMsg) Message() string {
	return p.Message_
}

func (p *PrivateMsg) Time() int64 {
	return p.Time_
}

func (p *PrivateMsg) Sender() Sender {
	return p.Sender_
}

func (p *PrivateMsg) MessageType() string {
	return "private"
}

func (p *PrivateMsg) Group() int64 {
	return 0
}

func (p *PrivateMsg) ReplyText(s string) error {
	return SendPrivateMsg(p.UserId, 0, s)
}
