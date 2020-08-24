package main

import (
	"github.com/CodFrm/iotqq-plugins/command/alimama"
	"github.com/CodFrm/iotqq-plugins/utils/iotqq"
	"log"
)

func dealUniversal(args iotqq.Message) bool {
	if cmd, ok := args.CommandMatch("有无(|.*?)($|\")"); ok {
		if str, err := alimama.Search(cmd[1]); err != nil {
			log.Println("有无", err)
			args.SendMessage("发生了一个系统错误")
		} else {
			args.SendMessage(str + "\n加群获取更多资讯,1131503629")
		}
		return true
	} else if _, ok := args.CommandMatch(".\\w{10,}."); ok && args.CurrentPacket.Data.Content[:3] == "淘" {
		ret, tkl, err := alimama.DealTkl(args.CurrentPacket.Data.Content[3:])
		if err != nil {
			log.Println("淘口令", err)
			args.SendMessage("发生了一个系统错误")
		} else if tkl == nil {
			args.SendMessage("没有发现淘口令")
		} else {
			args.SendMessage(ret + "\n" + "使用新的口令预计可反" + alimama.DealFl(tkl.Content[0].Tkfee3) + "￥")
			return true
		}
		return false
	} else if _, ok := args.CommandMatch("绑定(\\s|)(\\d+)($|\")"); ok {

	} else if cmd, ok := args.CommandMatch("订阅(\\s|)(.*?)($|\")"); ok && !args.Self() {
		if err := alimama.Subscribe(args.GetGroupId(), args.GetQQ(), cmd[2]); err != nil {
			log.Println("订阅", err)
			args.SendMessage("发生了一个系统错误")
		} else {
			args.SendMessage("订 阅 成 功")
		}
		return false
	} else if cmd, ok := args.CommandMatch("退订(\\s|)(.*?)($|\")"); ok && !args.Self() {
		if err := alimama.UnSubscribe(args.GetQQ(), cmd[2]); err != nil {
			log.Println("退订", err)
			args.SendMessage("发生了一个系统错误")
		} else {
			args.SendMessage("退 订 成 功")
		}
		return false
	} else if _, ok := args.CommandMatch("帮助"); ok {
		args.SendMessage("1.优惠购物,触发指令:'有无[物品名]',可获取商品列表和内部优惠券,选择你心爱的物品下单吧\n" +
			"2.淘口令转换,触发指令:'淘[淘宝口令]',可获取内部优惠券和优惠口令\n" +
			"3.订阅商品,触发指令:'订阅[关键字]',订阅指定的商品,如果有活动会直接私聊推送给你哦\n" +
			"饿了么每日红包:$nH3n1zNqDip$")
	}
	return false
}
