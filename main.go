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
	"io/ioutil"
	"log"
	"net/http"
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
			list, ok := val["GroupPic"].([]interface{})
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
			if _, ok := config.AppConfig.ManageGroupMap[args.CurrentPacket.Data.FromGroupID]; ok {
				for _, v := range picinfo {
					resp, err := http.Get(v.Url)
					if err != nil {
						continue
					}
					defer resp.Body.Close()
					picinfo[0].Byte, _ = ioutil.ReadAll(resp.Body)
					if resp.ContentLength > 1024*1024 {
						continue
					}
					if ok, err := command.IsAdult(v); err != nil {
						if ok == 1 {
							println(err)
						} else if ok == 2 {
							utils.SendMsg(args.CurrentPacket.Data.FromGroupID, args.CurrentPacket.Data.FromUserID, err.Error())
						} else if ok == 3 {
							utils.SendMsg(args.CurrentPacket.Data.FromGroupID, args.CurrentPacket.Data.FromUserID, err.Error())
							utils.RevokeMsg(args.CurrentPacket.Data.FromGroupID, args.CurrentPacket.Data.MsgSeq, args.CurrentPacket.Data.MsgRandom)
						}
					}
				}
			}
			content, ok := val["Content"].(string)
			if !ok {
				return
			}
			if picinfo[0].Byte == nil {
				resp, err := http.Get(picinfo[0].Url)
				if err != nil {
					return
				}
				defer resp.Body.Close()
				picinfo[0].Byte, _ = ioutil.ReadAll(resp.Body)
			}
			if strings.Index(content, "旋转图片") == 0 {
				cmd := strings.Split(strings.TrimFunc(content, func(r rune) bool {
					return r == '\r' || r == ' '
				}), " ")
				if !ok {
					return
				}
				utils.SendMsg(args.CurrentPacket.Data.FromGroupID, args.CurrentPacket.Data.FromUserID, "进行中,请稍后...")
				image, err := command.RotatePic(cmd[1:], picinfo[0])
				time.Sleep(time.Second * 2)
				if err != nil {
					utils.SendMsg(args.CurrentPacket.Data.FromGroupID, args.CurrentPacket.Data.FromUserID, "error:"+err.Error())
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
			} else if strings.Index(content, "图片鉴") == 0 && (strings.Index(content, "黄") != -1 || strings.Index(content, "色") != -1) {
				if ok, err := command.IsAdult(picinfo[0]); err != nil {
					if ok == 1 {
						println(err)
						utils.SendMsg(args.CurrentPacket.Data.FromGroupID, args.CurrentPacket.Data.FromUserID, "服务器开小差了,鉴图失败")
					} else if ok == 2 {
						utils.SendMsg(args.CurrentPacket.Data.FromGroupID, args.CurrentPacket.Data.FromUserID, "疑似色图")
					} else if ok == 3 {
						utils.SendMsg(args.CurrentPacket.Data.FromGroupID, args.CurrentPacket.Data.FromUserID, "就是色图,铐起来")
					} else if ok == 4 {
						utils.SendMsg(args.CurrentPacket.Data.FromGroupID, args.CurrentPacket.Data.FromUserID, err.Error())
					}
				} else {
					if strings.Index(content, "色") != -1 {
						str := utils.FileBase64("./data/img/1.jpg")
						utils.SendPicByBase64(args.CurrentPacket.Data.FromGroupID, args.CurrentPacket.Data.FromUserID, "", str)
					} else {
						utils.SendMsg(args.CurrentPacket.Data.FromGroupID, args.CurrentPacket.Data.FromUserID, "正常图片")
					}
				}
			}
		} else if args.CurrentPacket.Data.MsgType == "TextMsg" {
			regex := regexp.MustCompile("来((\\d*)份|点)好[康|看]的(.*?)(图|$)")
			ret := regex.FindStringSubmatch(args.CurrentPacket.Data.Content)
			if len(ret) > 0 {
				hkd(args, "", ret)
				return
			}
			regex = regexp.MustCompile("^来.?份[色|涩]图")
			ret = regex.FindStringSubmatch(args.CurrentPacket.Data.Content)
			if len(ret) > 0 {
				utils.SendMsg(args.CurrentPacket.Data.FromGroupID, args.CurrentPacket.Data.FromUserID, "我又不是搞颜色的机器人")
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
		} else if args.CurrentPacket.Data.MsgType == "AtMsg" {
			if strings.Index(args.CurrentPacket.Data.Content, "help") != -1 || strings.Index(args.CurrentPacket.Data.Content, "功能") != -1 ||
				strings.Index(args.CurrentPacket.Data.Content, "帮助") != -1 || strings.Index(args.CurrentPacket.Data.Content, "菜单") != -1 {
				utils.SendMsg(args.CurrentPacket.Data.FromGroupID, 0, "1.来点好康的,触发指令:'来1份好康的,来点好看的,来点好看的风景图',享受生活的美好\n"+
					"2.旋转图片,触发指令:'旋转图片 垂直/镜像/翻转/放大/缩小/高清重制 [图片]',更方便快捷的图片编辑\n"+
					"3.图片鉴黄,触发指令:'图片鉴黄/色 [图片]',让我们来猎杀那些色批\n"+
					"4.清理潜水,触发指令:'踢潜水 人数 舔狗/面子/普通模式',更方便快捷的清人工具,需要有管理员权限"+"还有更多神秘功能待你探索.")
				return
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
