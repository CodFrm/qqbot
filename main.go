package main

import (
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/CodFrm/qqbot/command"
	"github.com/CodFrm/qqbot/command/alimama"
	"github.com/CodFrm/qqbot/config"
	"github.com/CodFrm/qqbot/cqhttp"
	"github.com/CodFrm/qqbot/db"
	"github.com/CodFrm/qqbot/live"
	"github.com/CodFrm/qqbot/live/aria2"
	"github.com/CodFrm/qqbot/utils"
	"github.com/golang/glog"
)

var groupfile map[string]*os.File

func main() {
	if err := config.Init("config.yaml"); err != nil {
		log.Fatal(err)
	}
	aria2.DefaultRpc()
	if err := db.Init(); err != nil {
		log.Fatal(err)
	}
	if err := command.Init(); err != nil {
		log.Fatal(err)
	}
	if _, ok := config.AppConfig.FeatureMap["alimama"]; ok {
		if err := alimama.Init(); err != nil {
			log.Fatal(err)
		}
	}
	groupfile = make(map[string]*os.File)
	os.MkdirAll("data/group", os.ModeDir)

	c := cqhttp.NewClient(config.AppConfig.Addr)

	c.OnMessage(func(msg *cqhttp.MessageModel) {
		switch msg.MessageType() {
		case "guild":
			guild(msg)
		case "private":
			private(msg)
		}
	})

	if err := c.Start(); err != nil {
		glog.Fatalf("cqhttp: %v", err)
	}

}

// 查看当前频道id;打卡;查看video;查看可推流video;
// 播放(file);播放队列;添加到播放队列(file);
// 转码(file);转码进度;
// 直链下载到待转码 (url) (param);直链下载到可推流 (url) (param);
func guild(msg *cqhttp.MessageModel) {
	if _, ok := msg.CommandMatch("^查看当前频道id$"); ok {
		msg.ReplyText(fmt.Sprintf("guild: %v channel: %v user: %v", msg.Group(), msg.Message.(*cqhttp.GuildMsg).ChannelId, msg.Sender().UserId))
	} else if _, ok := msg.CommandMatch("^打卡$"); ok {
		if str, err := command.Sign(msg.Group(), msg.Message.(*cqhttp.GuildMsg).ChannelId, msg.Sender().UserId); err != nil {
			sendErr(msg, err)
		} else {
			if ok := command.IsWordGroup(msg.Group(), msg.Message.(*cqhttp.GuildMsg).ChannelId); ok {
				msg.ReplyText(str + ",请注意,本群后面将取消打卡指令,请分享/拍照进行打卡")
			} else {
				msg.ReplyText(str)
			}
		}
	} else if _, ok := msg.CommandMatch("^查看video$"); ok {
		m := ""
		for _, v := range live.ShowSource() {
			m += v + "; "
		}
		if m == "" {
			msg.ReplyText("没有可用的视频")
		} else {
			msg.ReplyText(m)
		}
	} else if _, ok := msg.CommandMatch("^查看可推流video"); ok {
		m := ""
		for _, v := range live.ShowLive() {
			m += v + "; "
		}
		if m == "" {
			msg.ReplyText("没有可用的视频")
		} else {
			msg.ReplyText(m)
		}
	} else if args, ok := msg.CommandMatch("^播放(.*?)$"); ok {
		if err := live.Play(msg.Group(), msg.Message.(*cqhttp.GuildMsg).ChannelId, msg.Sender().UserId, args[1]); err != nil {
			sendErr(msg, err)
		} else {
			msg.ReplyText("播放成功,请查看效果")
		}
	} else if args, ok := msg.CommandMatch("^播放队列$"); ok {
		if err := live.Play(msg.Group(), msg.Message.(*cqhttp.GuildMsg).ChannelId, msg.Sender().UserId, args[1]); err != nil {
			sendErr(msg, err)
		} else {
			msg.ReplyText("播放成功,请查看效果")
		}
	} else if args, ok := msg.CommandMatch("^添加到播放队列(.*?)"); ok {
		if err := live.Play(msg.Group(), msg.Message.(*cqhttp.GuildMsg).ChannelId, msg.Sender().UserId, args[1]); err != nil {
			sendErr(msg, err)
		} else {
			msg.ReplyText("播放成功,请查看效果")
		}
	} else if _, ok := msg.CommandMatch("^转码进度$"); ok {
		msg.ReplyText(live.TrProgress())
	} else if args, ok := msg.CommandMatch("^转码(.*?)$"); ok {
		if err := live.ToFlv(msg.Group(), msg.Message.(*cqhttp.GuildMsg).ChannelId, msg.Sender().UserId, args[1]); err != nil {
			sendErr(msg, err)
		} else {
			msg.ReplyText("转码中,输入\"转码进度\"查看进度")
		}
	} else if args, ok := msg.CommandMatch("^直链下载到待转码 (.*?) (.*?)$"); ok {
		if err := live.DownloadToSource(msg.Group(), msg.Message.(*cqhttp.GuildMsg).ChannelId, msg.Sender().UserId, args[1], args[2]); err != nil {
			sendErr(msg, err)
		} else {
			msg.ReplyText("下载中...请等待")
		}
	} else if args, ok := msg.CommandMatch("^直链下载到可推流 (.*?) (.*?)$"); ok {
		if err := live.DownloadToFlv(msg.Group(), msg.Message.(*cqhttp.GuildMsg).ChannelId, msg.Sender().UserId, args[1], args[2]); err != nil {
			sendErr(msg, err)
		} else {
			msg.ReplyText("下载中...请等待")
		}
	} else if _, ok := msg.CommandMatch("^查看下载进度$"); ok {
		if list, err := aria2.DownloadList(); err != nil {
			sendErr(msg, err)
		} else {
			msg.ReplyText("\n" + list.Table())
		}
	}
	return
}

func private(msg *cqhttp.MessageModel) {
	if _, ok := config.AppConfig.AdminQQMap[msg.Sender().UserId]; !ok {
		return
	}
	if args, ok := msg.CommandMatch("^帮(\\d+):(\\d+):(\\d+)推流 (.*?) (.*?)$"); ok {
		glog.Infof("帮%v推流: %v %v", msg.Self(), args[0], args[1])
		s, _ := url.PathUnescape(args[5])
		if err := live.AddLive(
			utils.StringToInt64(args[1]), utils.StringToInt64(args[2]), utils.StringToInt64(args[3]),
			args[4], s,
		); err != nil {
			sendErr(msg, err)
		} else {
			msg.ReplyText(fmt.Sprintf("推流绑定成功,帮助 %v 推流 绑定子频道: %v:%v 进行通知", args[3], args[1], args[2]))
		}
	}

}

func sendErr(m *cqhttp.MessageModel, err error) {
	m.ReplyText(err.Error())
}
