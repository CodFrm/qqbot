package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/CodFrm/qqbot/command"
	"github.com/CodFrm/qqbot/command/alimama"
	"github.com/CodFrm/qqbot/config"
	"github.com/CodFrm/qqbot/cqhttp"
	"github.com/CodFrm/qqbot/db"
	"github.com/CodFrm/qqbot/live"
	"github.com/CodFrm/qqbot/utils"
	"github.com/golang/glog"
)

var groupfile map[string]*os.File

func main() {
	if err := config.Init("config.yaml"); err != nil {
		log.Fatal(err)
	}
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

// 打卡;查看video;查看可推流video;播放(file);转码(file);转码进度;
func guild(msg *cqhttp.MessageModel) {
	if _, ok := msg.CommandMatch("^打卡$"); ok {
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
	} else if _, ok := msg.CommandMatch("^转码进度$"); ok {
		msg.ReplyText(live.TrProgress())
	} else if args, ok := msg.CommandMatch("^转码(.*?)$"); ok {
		if err := live.ToFlv(msg.Group(), msg.Message.(*cqhttp.GuildMsg).ChannelId, msg.Sender().UserId, args[1]); err != nil {
			sendErr(msg, err)
		} else {
			msg.ReplyText("转码中,输入\"转码进度\"查看进度")
		}
	} else if args, ok := msg.CommandMatch("^直链下载代转码 (.*?) (.*?)$"); ok {

		msg.ReplyText(args[2] + "下载失败,还不支持")
	} else if args, ok := msg.CommandMatch("^直链下载到可推流 (.*?) (.*?)$"); ok {

		msg.ReplyText(args[2] + "下载失败,还不支持")
	}
	return
}

func private(msg *cqhttp.MessageModel) {
	if _, ok := config.AppConfig.AdminQQMap[msg.Self()]; !ok {
		return
	}
	if args, ok := msg.CommandMatch("^帮(\\d+):(\\d+):(\\d+)推流 (.*?) (.*?)$"); ok {
		glog.Infof("帮%v推流: %v %v", msg.Self(), args[0], args[1])
		s, _ := url.PathUnescape(args[5])
		if err := live.AddLive(
			utils.StringToInt64(args[1]), utils.StringToInt64(args[2]), utils.StringToInt64(args[3]),
			args[4], strings.ReplaceAll(s, "&amp;", "&"),
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
