package taoke

import (
	"fmt"
	"strings"

	"github.com/CodFrm/qqbot/command/alimama"
	"github.com/CodFrm/qqbot/utils"
	"github.com/CodFrm/qqbot/utils/taobaoopen"
	"github.com/codfrm/cago/pkg/logger"
	"github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"go.uber.org/zap"
)

func (t *TaoKe) forward() {
	// 转发指定群消息
	zero.OnMessage(zero.OnlyGroup, func(ctx *zero.Ctx) bool {
		if ok, err := t.repo.IsEnableForward(); err != nil {
			return false
		} else if !ok {
			return false
		}
		if ok, err := t.repo.IsOriginGroup(ctx.Event.GroupID); err != nil {
			return false
		} else if !ok {
			return false
		}
		content := ctx.MessageString()
		//匹配淘口令发送
		if tkl := utils.RegexMatch(content, "([^\\w](\\w{8,12})[^\\w])|(.*?\\.jd\\.com\\/)"); len(tkl) > 0 {
			if strings.Index(content, "自助") != -1 || strings.Index(content, "网站") != -1 ||
				strings.Index(content, "luxbk.cn") != -1 || strings.Index(content, "群号") != -1 || strings.Index(content, "进分群") != -1 ||
				strings.Index(content, "本群") != -1 || (strings.Index(content, "群") != -1 && strings.Index(content, "进") != -1) {
				return false
			}
			if strings.Index(content, "饿了么") != -1 || strings.Index(content, "美团") != -1 {
				return false
			}
			return true
		}
		return false
	}).Handle(t.forwardToGroup)

	zero.OnMessage(zero.OnlyPrivate, func(ctx *zero.Ctx) bool {
		if tkl := utils.RegexMatch(ctx.MessageString(), "([^\\w](\\w{8,12})[^\\w])|(.*?\\.jd\\.com\\/)|(.*?\\.(taobao|tmall)\\.com\\/)"); len(tkl) > 0 {
			return true
		}
		return false
	}).SetPriority(30).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		ret, tkl, err := alimama.DealTkl(ctx.MessageString())
		if err != nil {
			logger.Default().Error("转链错误", zap.Error(err), zap.String("message", ctx.MessageString()))
			if err.Error() == "很抱歉！商品ID解析错误！！！" {
				ctx.Send("此商品不支持,无法搜索!")
				return
			}
			ctx.Send("发生了一个系统错误: " + err.Error())
			return
		} else if tkl == nil {
			ctx.Send("没有发现淘口令")
		} else if tkl.Content[0].Shorturl == "" {
			ctx.Send("此商品不支持,无法搜索!")
		} else {
			msg := ret + "\n约反:" + alimama.DealFl(tkl.Content[0].Tkfee3) + " "
			if tkl.Content[0].CouponInfoMoney != "" && tkl.Content[0].CouponInfoMoney != "0" {
				msg += "优惠券:" + tkl.Content[0].CouponInfoMoney + " 券后价:"
			} else {
				msg += "价格:"
			}
			msg += tkl.Content[0].QuanhouJiage + "￥"
			ctx.Send(msg)
		}
	})

	zero.OnCommand("订阅", zero.OnlyPrivate).Handle(func(ctx *zero.Ctx) {
		if err := t.repo.Subscribe(ctx.Event.UserID, ctx.Event.GroupID, ctx.State["args"].(string)); err != nil {
			ctx.Send("订阅失败")
			return
		}
		ctx.Send("订阅成功")
	})

	zero.OnCommand("取消订阅", zero.OnlyPrivate).Handle(func(ctx *zero.Ctx) {
		if err := t.repo.UnSubscribe(ctx.Event.UserID, ctx.State["args"].(string)); err != nil {
			ctx.Send("取消订阅失败")
			return
		}
		ctx.Send("取消订阅成功")
	})

	zero.OnCommand("订阅列表", zero.OnlyPrivate).Handle(func(ctx *zero.Ctx) {
		list, err := t.repo.SubscribeQQ()
		if err != nil {
			ctx.Send("获取订阅列表失败")
			return
		}
		msg := ""
		for qq := range list {
			msg += fmt.Sprintf("%d", qq) + "\n"
			topics, err := t.repo.SubscribeTopic(qq)
			if err != nil {
				logger.Default().Error("获取订阅列表失败", zap.Error(err), zap.Int64("qq", qq))
			}
			for _, topic := range topics {
				msg += topic + " "
			}
			msg += "\n"
		}
		ctx.Send(msg)
	})
}

func (t *TaoKe) forwardToGroup(ctx *zero.Ctx) {
	content := ctx.MessageString()
	if strings.HasPrefix(content, "转") {
		content = strings.TrimSpace(strings.TrimPrefix(content, "转"))
	}
	//非图片,直接转发
	groups, err := t.repo.ForwardGroupList()
	if err != nil {
		logger.Default().Error("获取转发列表失败", zap.Error(err), zap.String("content", content))
		return
	}
	//单独的口令
	cmd := utils.RegexMatch(content, "^[^\\w](\\w{8,12})[^\\w]$")
	subscribeMsg := ""
	if len(cmd) > 0 {
		_, tkl, err := t.DealTkl(content)
		if err != nil {
			logger.Default().Error("转发口令失败", zap.Error(err), zap.String("content", content))
			return
		}
		url := tkl.Content[0].PictURL
		content := "0-" + tkl.Content[0].TaoTitle + " " + tkl.Content[0].QuanhouJiage + "￥" + "\n" + tkl.Content[0].Tkl
		for _, v := range groups {
			if url == "" {
				ctx.SendGroupMessage(v, message.Text(content))
			} else {
				ctx.SendGroupMessage(v, []message.MessageSegment{
					message.Text(content),
					message.Image(url),
				})
			}
		}
		subscribeMsg = content
	} else {
		msg := make([]message.MessageSegment, 0)
		var tkl *taobaoopen.ConverseTkl
		for _, v := range ctx.Event.Message {
			if v.Type == "text" {
				if strings.HasPrefix(v.Data["text"], "/转") {
					v.Data["text"] = strings.TrimPrefix(v.Data["text"], "/转")
					v.Data["text"] = strings.TrimSpace(v.Data["text"])
				}
				if v.Data["text"] == "" {
					continue
				}
			}
			if v.Type == "at" {
				if v.Data["qq"] == "all" {
					msg = append(msg, message.Text("@全体成员"))
				}
				continue
			}
			msg = append(msg, v)
			if v.Type != "text" {
				continue
			}
			content, tmpTkl, err := t.DealTkl(v.Data["text"])
			if err != nil && err.Error() != "很抱歉！商品ID解析错误！！！" {
				logger.Default().Error("转发口令失败", zap.Error(err), zap.String("content", v.String()))
				return
			}
			if strings.HasPrefix(content, "转") {
				content = strings.TrimSpace(strings.TrimPrefix(content, "转"))
			}
			if tmpTkl != nil {
				tkl = tmpTkl
				v.Data["text"] = content
			}
			subscribeMsg += content + "\n"
		}
		if tkl != nil && t.repo.IsTklSend(tkl) {
			logger.Default().Error("重复发送", zap.String("content", content))
			return
		}
		for _, v := range groups {
			ctx.SendGroupMessage(v, msg)
		}
	}
	// 转发给订阅者
	go t.queueSend(ctx, subscribeMsg)
	return
}

func (t *TaoKe) queueSend(ctx *zero.Ctx, subscribeMsg string) {
	// 获取订阅者
	subscribe, err := t.repo.SubscribeQQ()
	if err != nil {
		logger.Default().Error("获取订阅者失败", zap.Error(err))
		return
	}
	// 发送订阅消息
	for qq, group := range subscribe {
		topics, err := t.repo.SubscribeTopic(qq)
		if err != nil {
			logger.Default().Error("获取订阅主题失败", zap.Error(err), zap.Int64("qq", qq))
			return
		}
		for _, topic := range topics {
			if strings.Contains(ctx.MessageString(), topic) {
				if group == 0 {
					ctx.SendPrivateMessage(qq, message.Text(
						fmt.Sprintf("0-您关注的'%s'有新消息\n%s\n回复'/退订 %s'可退订该消息.", topic, subscribeMsg, topic),
					))
					continue
				}
				t.SendPrivateMessage(ctx, qq, group, message.Message{
					message.Text(
						fmt.Sprintf("0-您关注的'%s'有新消息\n%s\n回复'/退订 %s'可退订该消息.", topic, subscribeMsg, topic),
					),
				})
			}
		}
	}

}

// SendPrivateMessage 发送私聊消息
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#send_private_msg-%E5%8F%91%E9%80%81%E7%A7%81%E8%81%8A%E6%B6%88%E6%81%AF
func (t *TaoKe) SendPrivateMessage(ctx *zero.Ctx, userID int64, groupId int64, message message.Message) int64 {
	rsp := ctx.CallAction("send_private_msg", zero.Params{
		"user_id":  userID,
		"group_id": groupId,
		"message":  message,
	}).Data.Get("message_id")
	if rsp.Exists() {
		logrus.Infof("[api] 发送私聊消息(%v): %v (id=%v)", userID, message.String(), rsp.Int())
		return rsp.Int()
	}
	return 0 // 无法获取返回值
}
