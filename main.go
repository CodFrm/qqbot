package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/CodFrm/iotqq-plugins/command"
	"github.com/CodFrm/iotqq-plugins/config"
	"github.com/CodFrm/iotqq-plugins/db"
	"github.com/CodFrm/iotqq-plugins/handler"
	"github.com/CodFrm/iotqq-plugins/model"
	"github.com/CodFrm/iotqq-plugins/utils"
	"github.com/CodFrm/iotqq-plugins/utils/iotqq"
	gosocketio "github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	"io/ioutil"
	"log"
	"math/rand"
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
	if err := c.On("OnGroupMsgs", func(h *gosocketio.Channel, args iotqq.Message) {
		if err := command.IsBlackList(strconv.FormatInt(args.CurrentPacket.Data.FromUserID, 10)); err != nil {
			return
		}
		if err := command.IsBlackList("group" + strconv.Itoa(args.CurrentPacket.Data.FromGroupID)); err != nil {
			return
		}
		if args.CurrentPacket.Data.MsgType == "XmlMsg" {
			handler.HandlerXmlMsg(args)
			return
		}
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
			if _, ok := config.AppConfig.ManageGroupMap[args.CurrentPacket.Data.FromGroupID]; ok && args.CurrentPacket.Data.FromUserID != args.CurrentQQ {
				for _, v := range picinfo {
					resp, err := http.Get(v.Url)
					if err != nil {
						continue
					}
					defer resp.Body.Close()
					v.Byte, _ = ioutil.ReadAll(resp.Body)
					if resp.ContentLength > 1024*1024 {
						continue
					}
					if ok, _, err := command.IsAdult(args.CurrentPacket.Data, v); err != nil {
						if ok == 1 {
							println(err)
						} else if ok == 2 {
							iotqq.SendMsg(args.CurrentPacket.Data.FromGroupID, 0, err.Error())
							return
						} else if ok == 3 {
							iotqq.SendMsg(args.CurrentPacket.Data.FromGroupID, 0, err.Error())
							iotqq.RevokeMsg(args.CurrentPacket.Data.FromGroupID, args.CurrentPacket.Data.MsgSeq, args.CurrentPacket.Data.MsgRandom)
							return
						}
					}
				}
			}
			content, _ := val["Content"].(string)
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
				args.SendMessage("进行中,请稍后...")
				image, err := command.RotatePic(cmd[1:], picinfo[0])
				time.Sleep(time.Second * 2)
				if err != nil {
					args.SendMessage("error:" + err.Error())
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
				iotqq.SendPicByBase64(args.CurrentPacket.Data.FromGroupID, args.CurrentPacket.Data.FromUserID, msg, base64Str)
				for k, v := range image[1:] {
					time.Sleep(time.Second * 2)
					base64Str, err := utils.ImageToBase64(v)
					msg := "@[GETUSERNICK(" + strconv.FormatInt(args.CurrentPacket.Data.FromUserID, 10) + ")]第" + strconv.Itoa(k+2) + "张图[PICFLAG]"
					if err != nil {
						msg = "@[GETUSERNICK(" + strconv.FormatInt(args.CurrentPacket.Data.FromUserID, 10) + ")]第" + strconv.Itoa(k+2) + "张发送失败," + err.Error()
					}
					iotqq.SendPicByBase64(args.CurrentPacket.Data.FromGroupID, args.CurrentPacket.Data.FromUserID, msg, base64Str)
				}
			} else if strings.Index(content, "图片鉴") == 0 && (strings.Index(content, "黄") != -1 || strings.Index(content, "色") != -1) {
				if ok, _, err := command.IsAdult(args.CurrentPacket.Data, picinfo[0]); err != nil {
					if ok == 1 {
						println(err)
						iotqq.SendMsg(args.CurrentPacket.Data.FromGroupID, args.CurrentPacket.Data.FromUserID, "服务器开小差了,鉴图失败")
					} else if ok == 2 {
						iotqq.SendMsg(args.CurrentPacket.Data.FromGroupID, args.CurrentPacket.Data.FromUserID, "疑似色图")
					} else if ok == 3 {
						iotqq.SendMsg(args.CurrentPacket.Data.FromGroupID, args.CurrentPacket.Data.FromUserID, "就是色图,铐起来")
					} else if ok == 4 {
						iotqq.SendMsg(args.CurrentPacket.Data.FromGroupID, args.CurrentPacket.Data.FromUserID, err.Error())
					}
				} else {
					if strings.Index(content, "色") != -1 {
						str := utils.FileBase64("./data/img/1.jpg")
						iotqq.SendPicByBase64(args.CurrentPacket.Data.FromGroupID, args.CurrentPacket.Data.FromUserID, "", str)
					} else {
						iotqq.SendMsg(args.CurrentPacket.Data.FromGroupID, args.CurrentPacket.Data.FromUserID, "正常图片")
					}
				}
			} else if ok := command.IsWordGroup(args.CurrentPacket.Data.FromGroupID); ok && !command.IsSign(args.CurrentPacket.Data.FromGroupID, args.CurrentPacket.Data.FromUserID) {
				if s, err := command.IsWordOk(args.CurrentPacket.Data.FromGroupID, args.CurrentPacket.Data.FromUserID, args.CurrentPacket.Data.Content); err != nil {
					sendErr(args, err)
					return
				} else if s != "" {
					args.SendMessage(s)
					return
				}
				if s, err := command.IsWordImage(args.CurrentPacket.Data.FromGroupID, args.CurrentPacket.Data.FromUserID, picinfo[0].Byte); err != nil {
					sendErr(args, err)
					return
				} else if s != "" {
					args.SendMessage(s)
					return
				}
			}
		} else if args.CurrentPacket.Data.MsgType == "TextMsg" {
			regex := regexp.MustCompile("^来((\\d*)份|点)好[康|看]的(.*?)(图|$)")
			ret := regex.FindStringSubmatch(args.CurrentPacket.Data.Content)
			if len(ret) > 0 {
				hkd(args, "", ret)
				return
			}
			if cmd := commandMatch(args.CurrentPacket.Data.Content, "^来(点|丶|份)(.*?)(图|$)$"); len(cmd) > 0 {
				hkd(args, "", []string{
					"", "", "", cmd[2],
				})
			} else if cmd := commandMatch(args.CurrentPacket.Data.Content, "^关联tag (.+?) (.+?)$"); len(cmd) > 0 {
				if _, ok := config.AppConfig.AdminQQMap[args.CurrentPacket.Data.FromUserID]; !ok {
					return
				}
				if err := command.RelateTag(cmd[1], cmd[2]); err != nil {
					iotqq.SendMsg(args.CurrentPacket.Data.FromGroupID, args.CurrentPacket.Data.FromUserID, err.Error())
					return
				}
				iotqq.SendMsg(args.CurrentPacket.Data.FromGroupID, args.CurrentPacket.Data.FromUserID, "OK")
				return
			} else if cmd := commandMatch(args.CurrentPacket.Data.Content, "黑名单 (.*?) (\\d)"); len(cmd) > 0 {
				if _, ok := config.AppConfig.AdminQQMap[args.CurrentPacket.Data.FromUserID]; !ok {
					return
				}
				if err := command.BlackList(cmd[1], cmd[2], ""); err != nil {
					sendErr(args, err)
					return
				}
				sendErr(args, errors.New("OK"))
			} else if cmd, ok := args.CommandMatch("^清理缓存\\s?(.*?$|$)"); ok && args.IsAdmin() {
				if err := command.CleanCache(cmd[1]); err != nil {
					sendErr(args, err)
					return
				}
				args.SendMessage("OK")
				return
			} else if _, ok := args.CommandMatch("^(当前|本群)场景$"); ok {
				list, err := command.QueryGroupScenes(args.CurrentPacket.Data.FromGroupID)
				if err != nil {
					sendErr(args, err)
					return
				}
				if len(list) == 0 {
					list = []string{"当前无场景"}
				}
				args.SendMessage(strings.Join(list, ","), func(o *iotqq.Options) {
					o.NotAt = true
				})
				return
			} else if cmd, ok := args.CommandMatch("^查询场景( p(\\d+)$| (.+?)$|$)"); ok {
				list, err := command.ScenesList(cmd[3], utils.StringToInt(cmd[2]))
				if err != nil {
					sendErr(args, err)
					return
				}
				if len(list) == 0 {
					list = []string{"没有了"}
				}
				args.SendMessage(strings.Join(list, ","), args.NotAt())
				return
			} else if cmd, ok := args.CommandMatch("^查看场景 (.+?)$"); ok {
				m, err := command.QueryScenesTag(cmd[1])
				if err != nil {
					sendErr(args, err)
					return
				}
				if len(m) == 0 {
					sendErr(args, errors.New("场景内容为空"))
					return
				}
				content := ""
				for k, v := range m {
					content += k + "=>" + v + " "
				}
				args.SendMessage(content, args.NotAt())
				return
			} else if cmd, ok := args.CommandMatch("^(添加|移除)场景 (.+?)$"); ok {
				if ok, _ := iotqq.IsAdmin(args.CurrentPacket.Data.FromGroupID, args.CurrentPacket.Data.FromUserID); !ok {
					args.SendMessage("你没有权限")
					return
				}
				s := strings.Split(cmd[2], ",")
				if cmd[1] == "添加" {
					if err := command.AddScenes(args.CurrentPacket.Data.FromGroupID, s); err != nil {
						sendErr(args, err)
						return
					}
				} else if cmd[1] == "移除" {
					if err := command.RemoveScenes(args.CurrentPacket.Data.FromGroupID, s); err != nil {
						sendErr(args, err)
						return
					}
				}
				args.SendMessage("OK")
				return
			} else if cmd, ok := args.CommandMatch("^映射tag (.+?) (.+?)$"); ok {
				if ok, _ := command.IsScenesOk(args.CurrentPacket.Data.FromGroupID); !ok {
					args.SendMessage("本群无权限")
					return
				}
				if ok, _ := iotqq.IsAdmin(args.CurrentPacket.Data.FromGroupID, args.CurrentPacket.Data.FromUserID); !ok {
					args.SendMessage("你没有权限")
					return
				}
				if err := command.CreateScenes("group"+strconv.Itoa(args.CurrentPacket.Data.FromGroupID), 2); err == nil {
					if err := command.AddScenes(args.CurrentPacket.Data.FromGroupID, []string{"group" + strconv.Itoa(args.CurrentPacket.Data.FromGroupID)}); err != nil {
						sendErr(args, err)
						return
					}
				}
				if err := command.AddScenesMap("group"+strconv.Itoa(args.CurrentPacket.Data.FromGroupID), cmd[1], cmd[2]); err != nil {
					sendErr(args, err)
					return
				}
				args.SendMessage("OK")
				return
			} else if cmd, ok := args.CommandMatch("^场景映射 (.+?) (.+?) (.+?)$"); ok {
				if args.CurrentPacket.Data.FromGroupID != 974381109 {
					args.SendMessage("无权限")
					return
				}
				if ok, _ := iotqq.IsAdmin(args.CurrentPacket.Data.FromGroupID, args.CurrentPacket.Data.FromUserID); !ok {
					args.SendMessage("你没有权限")
					return
				}
				if err := command.CreateScenes(cmd[1], 1); err != nil {
					if err.Error() != "场景存在" {
						sendErr(args, err)
						return
					}
				}
				if err := command.AddScenesMap(cmd[1], cmd[2], cmd[3]); err != nil {
					sendErr(args, err)
					return
				}
				args.SendMessage("OK")
				return
			} else if _, ok := args.CommandMatch("^打卡$"); ok {
				if str, err := command.Sign(args.CurrentPacket.Data.FromGroupID, args.CurrentPacket.Data.FromUserID); err != nil {
					sendErr(args, err)
				} else {
					args.SendMessage(str)
				}
				return
			} else if cmd, ok := args.CommandMatch("^(添加|删除)奖惩 (.+?)( (.*?)|)$"); ok {
				reargs := strings.Split(cmd[4], " ")
				if err := command.SetRewards(strconv.Itoa(args.CurrentPacket.Data.FromGroupID), args.CurrentPacket.Data.FromUserID, cmd[1] == "删除", cmd[2], reargs...); err != nil {
					sendErr(args, err)
				} else {
					args.SendMessage("OK")
				}
				return
			} else if cmd, ok := args.CommandMatch("^添加全群奖惩 (.+?)( (.*?)|)$"); ok {
				if ok, err := iotqq.IsAdmin(args.CurrentPacket.Data.FromGroupID, args.CurrentPacket.Data.FromUserID); err != nil {
					sendErr(args, err)
					return
				} else if !ok {
					args.SendMessage("你没有权限")
					return
				}
				reargs := strings.Split(cmd[3], " ")
				if err := command.AdminGroupReward(strconv.Itoa(args.CurrentPacket.Data.FromGroupID), false, cmd[1], reargs...); err != nil {
					sendErr(args, err)
					return
				}
				args.SendMessage("OK")
				return
			} else if cmd, ok := args.CommandMatch("^查看奖惩(| (.*?))$"); ok {
				if cmd[2] == "" {
					cmd[2] = strconv.Itoa(args.CurrentPacket.Data.FromGroupID)
				}
				val, err := command.GetRewards(cmd[2], args.CurrentPacket.Data.FromUserID)
				if err != nil {
					sendErr(args, err)
				} else {
					s := make([]string, 0)
					for _, v := range val {
						s = append(s, v.Command)
					}
					if len(s) <= 0 {
						args.SendMessage("你没有设置奖惩")
					} else {
						args.SendMessage(strings.Join(s, ","))
					}
				}
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
				iotqq.SendMsg(args.CurrentPacket.Data.FromGroupID, 0, args.CurrentPacket.Data.Content)
			}
		} else if args.CurrentPacket.Data.MsgType == "ReplayMsg" || args.CurrentPacket.Data.MsgType == "AtMsg" {
			if strings.Index(args.CurrentPacket.Data.Content, "求原图") != -1 {
				reg := regexp.MustCompile(`pixiv:(\d+)`)
				cmd := reg.FindStringSubmatch(args.CurrentPacket.Data.Content)
				if len(cmd) > 0 {
					args.SendMessage("原图较大,请耐心等待")
					imgbyte, err := command.GetPixivImg(cmd[1])
					if err != nil {
						time.Sleep(time.Second)
						args.SendMessage("系统错误,发送失败:" + err.Error())
						return
					}
					base64Str := base64.StdEncoding.EncodeToString(imgbyte)
					_, _ = iotqq.SendPicByBase64(args.CurrentPacket.Data.FromGroupID, args.CurrentPacket.Data.FromUserID, "原图收好\n[PICFLAG]", base64Str)
				}
			} else if m := commandMatch(args.CurrentPacket.Data.Content, "再来(一|亿)(点|份)"); len(m) > 0 {
				reg := regexp.MustCompile(`pixiv:(\d+)`)
				cmd := reg.FindStringSubmatch(args.CurrentPacket.Data.Content)
				if len(cmd) > 0 {
					args.SendMessage(" 图片检索中...请稍后")
					n := 1
					if m[1] == "亿" {
						n = rand.Intn(3) + 2
					}
					for i := 0; n > i; i++ {
						img, imgInfo, err := command.ZaiLaiYiDian(cmd[1])
						if err != nil {
							if err.Error() == "我真的一张都没有了" {
								args.SendMessage(" " + err.Error())
								return
							}
							args.SendMessage(" 服务器开小差了,搜索失败T T,稍后再试一次吧")
							return
						}
						base64Str := base64.StdEncoding.EncodeToString(img)
						msg := "pixiv:" + imgInfo.Id + " " + imgInfo.Title + " 画师:" + imgInfo.UserName + "\n" + "https://www.pixiv.net/artworks/" + imgInfo.Id + "\n[PICFLAG]"
						_, _ = iotqq.SendPicByBase64(args.CurrentPacket.Data.FromGroupID, args.CurrentPacket.Data.FromUserID, msg, base64Str)
						time.Sleep(time.Second * 3)
					}
				}
			} else if cmd := commandMatch(args.CurrentPacket.Data.Content, "黑名单(.*?)(\\d)"); len(cmd) > 0 {
				if _, ok := config.AppConfig.AdminQQMap[args.CurrentPacket.Data.FromUserID]; !ok {
					return
				}
				m := &struct {
					UserID []int64 `json:"UserID"`
				}{}
				if err := json.Unmarshal([]byte(args.CurrentPacket.Data.Content), m); err != nil {
					sendErr(args, err)
					return
				}
				for _, v := range m.UserID {
					command.BlackList(strconv.FormatInt(v, 10), cmd[2], "")
				}
				sendErr(args, errors.New("OK"))
			} else if cmd := commandMatch(args.CurrentPacket.Data.Content, "给我(康康|看看)"); len(cmd) > 0 {
				cmd := commandMatch(args.CurrentPacket.Data.Content, "图片已撤回,证据已保留ID:(\\w+)")
				if len(cmd) > 0 {
					if b, err := command.Gwkk(cmd[1]); err != nil {
						sendErr(args, err)
						return
					} else {
						base64Str := base64.StdEncoding.EncodeToString(b)
						if _, err := iotqq.SendFriendPicMsg(args.CurrentPacket.Data.FromUserID, "", base64Str); err != nil {
							sendErr(args, err)
							return
						}
						time.Sleep(time.Second)
						sendErr(args, errors.New("已私聊发送"))
					}
				}
			} else if strings.Index(args.CurrentPacket.Data.Content, "nmsl") != -1 {
				args.SendMessage(" " + utils.Nmsl())
				return
			} else if (strings.Index(args.CurrentPacket.Data.Content, "help") != -1 || strings.Index(args.CurrentPacket.Data.Content, "功能") != -1 ||
				strings.Index(args.CurrentPacket.Data.Content, "帮助") != -1 || strings.Index(args.CurrentPacket.Data.Content, "菜单") != -1) && args.CurrentPacket.Data.FromUserID != args.CurrentQQ {
				if strings.Index(args.CurrentPacket.Data.Content, "帮助 图片场景") != -1 {
					args.SendMessage("可设置当前群的场景内容,将tag映射到另外一个tag实现更加有趣的图片搜索\n" +
						"1.3.1.当前场景,触发指令:'当前/本群场景',查询本群场景\n" +
						"1.3.2.添加和移除场景,触发指令:'添加/移除场景 [场景名]'*,对本群添加或者删除场景\n" +
						"1.3.3.查询场景,触发指令:'查询场景 [场景关键字(可选)]/p[页码(可选)]',寻找有趣的场景设置到群\n" +
						"1.3.4.查看场景,触发指令:'查看场景 [场景名]',查看场景中tag的映射表\n" +
						"1.3.5.映射tag,触发指令:'映射tag [映射tag] [被映射tag]'*,自定义群内tag映射表(暂未开放)\n" +
						"1.3.6.场景映射,触发指令:'场景映射 [场景名] [映射tag] [被映射tag]'*,需要高级权限")
					return
				} else if strings.Index(args.CurrentPacket.Data.Content, "帮助 每日打卡") != -1 {
					args.SendMessage("帮助你坚持完成一件事情\n" +
						"5.1.打卡,触发指令:'打卡',完成今天的任务\n" +
						"5.2.添加/删除奖惩,触发指令:'添加/删除奖惩 [奖惩方案] [参数]',打卡成功或者失败后将进行奖惩,让行动更有动力\n" +
						"5.3.查看奖惩,触发指令:'查看奖惩',查看我选择的奖惩方案\n" +
						"奖惩方案:1.设置名片 打卡成功后帮你修改群内名片,2.nmsl,3.nmsl单词特供版(陆续开发中)\n" +
						"eg.添加奖惩 设置名片 打卡第N天")
					return
				}
				iotqq.SendMsg(args.CurrentPacket.Data.FromGroupID, 0, "1.来点好康的,触发指令:'来1份好康的,来点好看的,来点好看的风景图',享受生活的美好\n"+
					"1.1.求原图,触发指令:'回复+求原图',可获得原图内容\n"+
					"1.2.再来一点,触发指令:'回复+再来一/亿点',可获得更多好康的\n"+
					"1.3.场景功能,触发指令:'帮助 图片场景',让好康的更好玩")
				time.Sleep(time.Second)
				iotqq.SendMsg(args.CurrentPacket.Data.FromGroupID, 0, "2.旋转图片,触发指令:'旋转图片 垂直/镜像/翻转/放大/缩小/灰白/颜色反转/高清重制 [图片]',更方便快捷的图片编辑\n"+
					"3.图片鉴黄,触发指令:'图片鉴黄/色 [图片]',让我们来猎杀那些色批(默认不会开启自动鉴黄功能)\n"+
					"3.1给我康康,触发指令:'回复+给我康康/看看',成为专业鉴黄师\n"+
					"4.清理潜水,触发指令:'踢潜水 人数 舔狗/面子/普通模式'*,更方便快捷的清人工具,需要有管理员权限\n"+
					"5.每日打卡,触发指令:'帮助 每日打卡'\n"+
					"还有更多神秘功能待你探索.")
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

func sendErr(m iotqq.Message, err error) {
	iotqq.SendMsg(m.CurrentPacket.Data.FromGroupID, m.CurrentPacket.Data.FromUserID, err.Error())
}

func commandMatch(content string, command string) []string {
	reg := regexp.MustCompile(command)
	return reg.FindStringSubmatch(content)
}

func hkd(args iotqq.Message, at string, commandstr []string) error {
	num, _ := strconv.Atoi(commandstr[2])
	if num <= 0 {
		num = 1
	} else if num > 4 {
		args.SendMessage(" 注意身体")
		return errors.New("注意身体")
	}
	args.SendMessage(" 图片搜索中...请稍后")
	go func() {
		mapTag, err := command.QueryMap(args.CurrentPacket.Data.FromGroupID, commandstr[3])
		if err != nil || mapTag == "" {
			mapTag = commandstr[3]
		}
		for i := 0; i < num; i++ {
			img, imgInfo, err := command.HaoKangDe(mapTag)
			if err != nil {
				if err.Error() == "图片过少" {
					args.SendMessage(" " + err.Error())
					return
				}
				args.SendMessage(" 服务器开小差了,搜索失败T T,稍后再试一次吧")
				println(err.Error())
				return
			}
			n, flag, err := command.IsAdult(args.CurrentPacket.Data, &model.PicInfo{
				Url:  "",
				Byte: img,
			})
			if n == 2 || n == 3 {
				_ = command.BanR18(commandstr[3])
				if flag {
					return
				}
				args.SendMessage("触发颜色警告,一定限度后进入黑名单")
				return
			}
			base64Str := base64.StdEncoding.EncodeToString(img)
			msg := "您的"
			if num == 1 {
				msg += commandstr[3] + "图收好\n"
			} else {
				msg += strconv.Itoa(num) + "份" + commandstr[3] + "图收好\n"
			}
			if i >= 1 {
				msg = ""
			}
			db.Redis.Set("pixiv:send:qq:"+imgInfo.Id, args.CurrentPacket.Data.FromUserID, time.Hour)
			msg += "pixiv:" + imgInfo.Id + " " + commandstr[3] + " " + imgInfo.Title + " 画师:" + imgInfo.UserName + "\n" + "https://www.pixiv.net/artworks/" + imgInfo.Id + "\n[PICFLAG]"
			_, _ = iotqq.SendPicByBase64(args.CurrentPacket.Data.FromGroupID, args.CurrentPacket.Data.FromUserID, msg, base64Str)
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
