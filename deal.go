package main

import (
	"encoding/json"
	"github.com/CodFrm/iotqq-plugins/command/alimama"
	"github.com/CodFrm/iotqq-plugins/config"
	"github.com/CodFrm/iotqq-plugins/utils"
	"github.com/CodFrm/iotqq-plugins/utils/iotqq"
	"github.com/CodFrm/iotqq-plugins/utils/taobaoopen"
	"log"
	"net/url"
)

func dealUniversal(args iotqq.Message) bool {
	if args.Self() {
		return false
	}
	pic := &iotqq.PicMsgContent{}
	if err := json.Unmarshal([]byte(args.CurrentPacket.Data.Content), pic); err == nil {
		args.CurrentPacket.Data.Content = pic.Content
	}
	if cmd, ok := args.CommandMatch("有无(|.*?)($|\")"); ok {
		if str, err := alimama.Search(cmd[1]); err != nil {
			log.Println("有无", err)
			args.SendMessage("发生了一个系统错误")
		} else {
			args.SendMessage(str + "\n另可以发送'订阅" + cmd[1] + "'来关注本类商品哦\n加群获取更多资讯,1131503629")
		}
		return true
	} else if _, ok := args.CommandMatch("([\\p{Sc}](\\w{8,12})[\\p{Sc}])|(.*?\\.jd\\.com\\/)|(.*?\\.(taobao|tmall)\\.com\\/)"); ok {
		var tkl *taobaoopen.ConverseTkl
		_, tkl, err := alimama.DealTklFl(args.CurrentPacket.Data.Content[3:])
		if err != nil {
			if err.Error() == "很抱歉！商品ID解析错误！！！" {
				args.SendMessage("此商品不支持,无法搜索!")
				return true
			}
			log.Println("淘口令", err)
			args.SendMessage("发生了一个系统错误")
		} else if tkl == nil {
			args.SendMessage("没有发现淘口令")
		} else if tkl.Content[0].Shorturl == "" {
			args.SendMessage("此商品不支持,无法搜索!")
		} else {
			msg := "约反:" + alimama.DealFl(tkl.Content[0].Tkfee3) + " "
			if tkl.Content[0].CouponInfoMoney != "" && tkl.Content[0].CouponInfoMoney != "0" {
				msg += "优惠券:" + tkl.Content[0].CouponInfoMoney + " 券后价:"
			} else {
				msg += "价格:"
			}
			msg += tkl.Content[0].QuanhouJiage
			kl := utils.RegexMatch(tkl.Content[0].Tkl, "[\\p{Sc}](\\w{8,12})[\\p{Sc}]")
			if len(kl) == 2 {
				msg += "\n" + alimama.ShortUrl("http://tb.icodef.com/tb.php?tkl="+url.QueryEscape(kl[1])+"&pic="+url.QueryEscape(tkl.Content[0].PictURL))
				args.SendMessage(msg)
			} else {
				args.SendMessage("系统发生了一个错误")
			}
		}
		return true
	} else if _, ok := args.CommandMatch("绑定(\\s|)(\\d+)($|\")"); ok {

	} else if cmd, ok := args.CommandMatch("订阅(\\s|)(.*?)($|\")"); ok && !args.Self() {
		if cmd[2] == "" {
			args.SendMessage("请输入订阅关键字")
			return true
		}
		if err := alimama.Subscribe(args.GetGroupId(), args.GetQQ(), cmd[2]); err != nil {
			log.Println("订阅", err)
			args.SendMessage("发生了一个系统错误")
		} else {
			args.SendMessage(cmd[2] + " 订 阅 成 功")
		}
		return true
	} else if cmd, ok := args.CommandMatch("退订(\\s|)(.*?)($|\")"); ok && !args.Self() {
		if err := alimama.UnSubscribe(args.GetQQ(), cmd[2]); err != nil {
			log.Println("退订", err)
			args.SendMessage("发生了一个系统错误")
		} else {
			args.SendMessage(cmd[2] + " 退 订 成 功")
		}
		return true
	} else if _, ok := args.CommandMatch("帮助"); ok {
		args.SendMessage("1.优惠购物,触发指令:'有无[物品名]',可获取商品列表和内部优惠券,选择你心爱的物品下单吧\n" +
			"2.优惠查券,触发指令:'[口令或者链接]',可查询优惠券和优惠口令\n" +
			"3.订阅商品,触发指令:'订阅[关键字]',订阅指定的商品,如果有活动会直接私聊推送给你哦\n" +
			"饿了么每日红包:$5YiUccAeTlY$")
		return true
	}
	//管理员命令
	if _, ok := config.AppConfig.AdminQQMap[args.CurrentPacket.Data.FromUin]; !ok {
		return false
	}
	if cmd, ok := args.CommandMatch("^(开启|关闭)转发$"); ok {
		args.Err(alimama.EnableGroupForward(cmd[1] == "开启"))
		return true
	} else if cmd, ok := args.CommandMatch("^(添加|删除)转发群(\\d+)$"); ok {
		group := utils.StringToInt(cmd[2])
		if cmd[1] == "添加" {
			args.Err(alimama.AddForwardGroup(group))
		} else {
			args.Err(alimama.RemoveForwardGroup(group))
		}
		return true
	}
	return false
}
