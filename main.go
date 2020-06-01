package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/CodFrm/iotqq-plugins/command"
	"github.com/CodFrm/iotqq-plugins/config"
	"github.com/CodFrm/iotqq-plugins/db"
	"github.com/CodFrm/iotqq-plugins/model"
	"github.com/CodFrm/iotqq-plugins/utils"
	gosocketio "github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

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
	c, err := gosocketio.Dial(
		gosocketio.GetUrl(config.AppConfig.Addr, config.AppConfig.Port, false),
		transport.GetDefaultWebsocketTransport())
	if err != nil {
		log.Fatal(err)
	}
	err = c.On(gosocketio.OnDisconnection, func(h *gosocketio.Channel) {
		//log.Fatal("Disconnected")
	})
	if err != nil {
		log.Fatal(err)
	}
	err = c.On(gosocketio.OnConnection, func(h *gosocketio.Channel) {
		log.Println("Connected")
	})
	if err != nil {
		log.Fatal(err)
	}
	lastContent := make(map[int]string)
	lastNum := make(map[int]int)
	if err := c.On("OnGroupMsgs", func(h *gosocketio.Channel, args model.Message) {
		if args.CurrentPacket.Data.MsgType == "PicMsg" {
			val := make(map[string]interface{})
			if err := json.Unmarshal([]byte(args.CurrentPacket.Data.Content), &val); err != nil {
				return
			}
			content, ok := val["Content"].(string)
			if !ok {
				return
			}
			if strings.Index(content, "旋转图片") == 0 {
				cmd := strings.Split(strings.TrimFunc(content, func(r rune) bool {
					return r == '\r' || r == ' '
				}), " ")
				list, ok := val["GroupPic"].([]interface{})
				if !ok {
					return
				}
				picinfo := make([]*model.PicInfo, 0)
				for _, v := range list {
					m, ok := v.(map[string]interface{})
					if !ok {
						continue
					}
					url, ok := m["Url"].(string)
					if !ok {
						continue
					}
					picinfo = append(picinfo, &model.PicInfo{Url: url})
				}
				if len(picinfo) == 0 {
					return
				}
				image, err := command.RotatePic(cmd[1:], picinfo[0])
				if err != nil {
					utils.SendMsg(args.CurrentPacket.Data.FromGroupID, args.CurrentPacket.Data.FromUserID, " error:"+err.Error())
					return
				}
				if len(image) == 0 {
					return
				}
				msg := "@[GETUSERNICK(" + strconv.FormatInt(args.CurrentPacket.Data.FromUserID, 10) + ")]一共" + strconv.Itoa(len(image)) + "张图片,请准备接收~[PICFLAG]"
				base64Str, err := utils.ImageToBase64(image[0])
				if err != nil {
					msg += ",第1张发送失败," + err.Error()
				}
				utils.SendPicByBase64(args.CurrentPacket.Data.FromGroupID, args.CurrentPacket.Data.FromUserID, msg, base64Str)
				for k, v := range image[1:] {
					time.Sleep(time.Second * 2)
					base64Str, err := utils.ImageToBase64(v)
					msg := "@[GETUSERNICK(" + strconv.FormatInt(args.CurrentPacket.Data.FromUserID, 10) + ")]第" + strconv.Itoa(k+2) + "张图[PICFLAG]"
					if err != nil {
						msg = "@[GETUSERNICK(" + strconv.FormatInt(args.CurrentPacket.Data.FromUserID, 10) + ")]第" + strconv.Itoa(k+2) + "张发送失败," + err.Error()
					}
					utils.SendPicByBase64(args.CurrentPacket.Data.FromGroupID, args.CurrentPacket.Data.FromUserID, msg, base64Str)
				}
			}
		} else if args.CurrentPacket.Data.MsgType == "TextMsg" {
			regex := regexp.MustCompile("^来((\\d*)份|点)好[康|看]的(.*?)(图|$)")
			ret := regex.FindStringSubmatch(args.CurrentPacket.Data.Content)
			if len(ret) > 0 {
				hkd(args, "", ret)
				return
			}
			regex = regexp.MustCompile("^来.?份[色|涩]图")
			ret = regex.FindStringSubmatch(args.CurrentPacket.Data.Content)
			if len(ret) > 0 {
				utils.SendMsg(args.CurrentPacket.Data.FromGroupID, 0, "我又不是搞颜色的机器人")
				return
			}

			groupid := args.CurrentPacket.Data.FromGroupID
			if lastContent[groupid] == args.CurrentPacket.Data.Content {
				lastNum[groupid]++
			} else {
				lastNum[groupid] = 0
			}
			lastContent[groupid] = args.CurrentPacket.Data.Content
			if lastNum[groupid] == 2 {
				utils.SendMsg(args.CurrentPacket.Data.FromGroupID, 0, args.CurrentPacket.Data.Content)
			}
		}

	}); err != nil {
		log.Fatal(err)
	}
	if err := c.On("OnEvents", func(h *gosocketio.Channel, args interface{}) {
		//println(args)
	}); err != nil {
		log.Fatal(err)
	}
	SendJoin(c)
	for {
		select {
		case <-time.After(time.Second * 600):
			{
				SendJoin(c)
				println("doing...")
			}
		}
	}
}

func hkd(args model.Message, at string, commandstr []string) error {
	num, _ := strconv.Atoi(commandstr[2])
	if num <= 0 {
		num = 1
	} else if num > 4 {
		utils.SendMsg(args.CurrentPacket.Data.FromGroupID, args.CurrentPacket.Data.FromUserID, " 注意身体")
		return errors.New("注意身体")
	}
	utils.SendMsg(args.CurrentPacket.Data.FromGroupID, args.CurrentPacket.Data.FromUserID, " 图片搜索中...请稍后")
	go func() {
		for i := 0; i < num; i++ {
			img, err := command.HaoKangDe(commandstr[3])
			if err != nil {
				if err.Error() == "图片过少" {
					utils.SendMsg(args.CurrentPacket.Data.FromGroupID, args.CurrentPacket.Data.FromUserID, " "+err.Error())
					return
				}
				utils.SendMsg(args.CurrentPacket.Data.FromGroupID, args.CurrentPacket.Data.FromUserID, " 服务器开小差了,搜索失败T T,稍后再试一次吧")
				println(err.Error())
				return
			}
			//TODO:鉴黄,来源
			base64Str := base64.StdEncoding.EncodeToString(img)
			msg := "@[GETUSERNICK(" + strconv.FormatInt(args.CurrentPacket.Data.FromUserID, 10) + ")] 您的"
			if num == 1 {
				msg += commandstr[3] + "图收好\n[PICFLAG]"
			} else {
				msg += strconv.Itoa(num) + "份" + commandstr[3] + "图收好\n[PICFLAG]"
			}
			if i >= 1 {
				msg = ""
			}
			_, _ = utils.SendPicByBase64(args.CurrentPacket.Data.FromGroupID, args.CurrentPacket.Data.FromUserID, msg, base64Str)
			time.Sleep(time.Second * 3)
		}
	}()
	return nil
}

func SendJoin(c *gosocketio.Client) {
	log.Println("获取QQ号连接")
	result, err := c.Ack("GetWebConn", config.AppConfig.QQ, time.Second*5)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("emit", result)
	}
}
