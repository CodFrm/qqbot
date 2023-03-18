package taoke

import (
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

func (t *TaoKe) admin() {
	zero.OnCommand("开启转发", zero.OnlyPrivate, zero.AdminPermission).Handle(func(ctx *zero.Ctx) {
		if err := t.repo.EnableForward(); err != nil {
			ctx.SendChain(message.Text("转发开启失败"))
			return
		}
		ctx.SendChain(message.Text("转发开启成功"))
	})
	zero.OnCommand("关闭转发", zero.OnlyPrivate, zero.AdminPermission).Handle(func(ctx *zero.Ctx) {
		if err := t.repo.DisableForward(); err != nil {
			ctx.SendChain(message.Text("转发关闭失败"))
			return
		}
		ctx.SendChain(message.Text("转发关闭成功"))
	})
	zero.OnCommand("转发列表", zero.OnlyPrivate, zero.AdminPermission).
		SetBlock(true).SetPriority(20).
		Handle(func(ctx *zero.Ctx) {
			groups, err := t.repo.ForwardGroupList()
			if err != nil {
				ctx.SendChain(message.Text("获取转发列表失败"))
				return
			}
			china := make([]message.MessageSegment, 0)
			china = append(china, message.Text("转发列表: "))
			for _, group := range groups {
				china = append(china, message.Text(group), message.Text(" "))
			}
			ctx.SendChain(china...)
		})
	zero.OnCommand("添加转发", zero.OnlyPrivate, zero.AdminPermission).Handle(func(ctx *zero.Ctx) {
		group := ctx.State["args"].(string)
		if err := t.repo.AddForwardGroup(group); err != nil {
			ctx.SendChain(message.Text("添加转发失败"))
			return
		}
		ctx.SendChain(message.Text("添加转发成功"))
	})
	zero.OnCommand("删除转发", zero.OnlyPrivate, zero.AdminPermission).Handle(func(ctx *zero.Ctx) {
		group := ctx.State["args"].(string)
		if err := t.repo.RemoveForwardGroup(group); err != nil {
			ctx.SendChain(message.Text("删除转发失败"))
			return
		}
		ctx.SendChain(message.Text("删除转发成功"))
	})
	zero.OnCommand("添加来源", zero.OnlyPrivate, zero.AdminPermission).Handle(func(ctx *zero.Ctx) {
		group := ctx.State["args"].(string)
		if err := t.repo.AddSourceGroup(group); err != nil {
			ctx.SendChain(message.Text("添加来源失败"))
			return
		}
		ctx.SendChain(message.Text("添加来源成功"))
	})
	zero.OnCommand("删除来源", zero.OnlyPrivate, zero.AdminPermission).Handle(func(ctx *zero.Ctx) {
		group := ctx.State["args"].(string)
		if err := t.repo.RemoveSourceGroup(group); err != nil {
			ctx.SendChain(message.Text("删除来源失败"))
			return
		}
		ctx.SendChain(message.Text("删除来源成功"))
	})
	zero.OnCommand("来源列表", zero.OnlyPrivate, zero.AdminPermission).Handle(func(ctx *zero.Ctx) {
		groups, err := t.repo.SourceGroupList()
		if err != nil {
			ctx.SendChain(message.Text("获取来源列表失败"))
			return
		}
		china := make([]message.MessageSegment, 0)
		china = append(china, message.Text("来源列表: "))
		for _, group := range groups {
			china = append(china, message.Text(group), message.Text(" "))
		}
		ctx.SendChain(china...)
	})
	zero.OnCommand("转", zero.OnlyPrivate, zero.AdminPermission).
		SetPriority(30).
		Handle(t.forwardToGroup)
}
