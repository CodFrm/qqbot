package main

import (
	"fmt"
	"log"
	"os"

	"github.com/CodFrm/iotqq-plugins/command"
	"github.com/CodFrm/iotqq-plugins/command/alimama"
	"github.com/CodFrm/iotqq-plugins/config"
	"github.com/CodFrm/iotqq-plugins/cqhttp"
	"github.com/CodFrm/iotqq-plugins/db"
	"github.com/CodFrm/iotqq-plugins/live"
	"github.com/CodFrm/iotqq-plugins/utils"
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
	} else if args, ok := msg.CommandMatch("^转码(.*?)$"); ok {

		msg.ReplyText(args[1] + "转码失败,还不支持")
	} else if args, ok := msg.CommandMatch("^直链下载代转码 (.*?) (.*?)$"); ok {

		msg.ReplyText(args[2] + "下载失败,还不支持")
	} else if args, ok := msg.CommandMatch("^直链下载到可推流 (.*?) (.*?)$"); ok {

		msg.ReplyText(args[2] + "下载失败,还不支持")
	}
	return
}

func private(msg *cqhttp.MessageModel) {
	if args, ok := msg.CommandMatch("^帮(\\d+):(\\d+):(\\d+)推流 (.*?) (.*?)$"); ok {
		glog.Infof("帮%v推流: %v %v", msg.Self(), args[0], args[1])
		if err := live.AddLive(
			utils.StringToInt64(args[1]), utils.StringToInt64(args[2]), utils.StringToInt64(args[3]),
			args[4], args[5],
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
