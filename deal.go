package main

import (
	"github.com/CodFrm/iotqq-plugins/command/alimama"
	"github.com/CodFrm/iotqq-plugins/utils/iotqq"
)

func dealUniversal(args iotqq.Message) bool {
	if cmd, ok := args.CommandMatch("有无(|.*?)($|\")"); ok {
		if str, err := alimama.Search(cmd[1]); err != nil {
			sendErr(args, err)
		} else {
			args.SendMessage(str + "\n加群获取更多资讯,1131503629")
		}
		return true
	} else if _, ok := args.CommandMatch("绑定(\\s|)(\\d+)($|\")"); ok {

	}
	return false
}
