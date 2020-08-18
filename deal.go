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
	} else if _, ok := args.CommandMatch(".\\w{10,}."); ok && args.CurrentPacket.Data.Content[:3] == "淘" {
		ret, tkl, err := alimama.DealTkl(args.CurrentPacket.Data.Content[3:])
		if err != nil {
			sendErr(args, err)
		} else if tkl == nil {
			args.SendMessage("没有发现淘口令")
		} else {
			args.SendMessage(ret + "\n" + "使用新的口令预计可反" + alimama.DealFl(tkl.Content[0].Tkfee3) + "￥")
			return true
		}
		return false
	} else if _, ok := args.CommandMatch("绑定(\\s|)(\\d+)($|\")"); ok {

	}
	return false
}
