package handler

import (
	"github.com/CodFrm/iotqq-plugins/command"
	"github.com/CodFrm/iotqq-plugins/utils/iotqq"
)

func HandlerXmlMsg(args iotqq.Message) {
	if s, _ := command.IsWordOk(args.CurrentPacket.Data.FromGroupID, args.CurrentPacket.Data.FromUserID, args.CurrentPacket.Data.Content); s != "" {
		args.SendMessage(s)
		return
	}
	s, err := command.BindWebsite(args.CurrentPacket.Data.FromGroupID, args.CurrentPacket.Data.FromUserID, args.CurrentPacket.Data.Content)
	if err != nil {
		args.SendMessage(err.Error())
		return
	} else if s != "" {
		args.SendMessage(s)
		return
	}
}
