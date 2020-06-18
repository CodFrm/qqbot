package iotqq

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/CodFrm/iotqq-plugins/config"
	"github.com/CodFrm/iotqq-plugins/db"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var url1, qq string

type QQinfo struct {
	Code    int    `json:"code"`
	Data    Data1  `json:"data"`
	Default int    `json:"default"`
	Message string `json:"message"`
	Subcode int    `json:"subcode"`
}
type Data1 struct {
	AvatarURL     string `json:"avatarUrl"`
	Bitmap        string `json:"bitmap"`
	Commfrd       int    `json:"commfrd"`
	Friendship    int    `json:"friendship"`
	Greenvip      int    `json:"greenvip"`
	IntimacyScore int    `json:"intimacyScore"`
	IsFriend      int    `json:"isFriend"`
	Logolabel     string `json:"logolabel"`
	Nickname      string `json:"nickname"`
	Qqvip         int    `json:"qqvip"`
	Qzone         int    `json:"qzone"`
	Realname      string `json:"realname"`
	Redvip        int    `json:"redvip"`
	Smartname     string `json:"smartname"`
	Uin           int    `json:"uin"`
}
type QQ struct {
	Cont int
}
type PSkey struct {
	Connect     string `json:"connect"`
	Docs        string `json:"docs"`
	Docx        string `json:"docx"`
	Game        string `json:"game"`
	Gamecenter  string `json:"gamecenter"`
	Imgcache    string `json:"imgcache"`
	MTencentCom string `json:"m.tencent.com"`
	Mail        string `json:"mail"`
	Mma         string `json:"mma"`
	Now         string `json:"now"`
	Office      string `json:"office"`
	Openmobile  string `json:"openmobile"`
	Qqweb       string `json:"qqweb"`
	Qun         string `json:"qun"`
	Qzone       string `json:"qzone"`
	QzoneCom    string `json:"qzone.com"`
	TenpayCom   string `json:"tenpay.com"`
	Ti          string `json:"ti"`
	Vip         string `json:"vip"`
	Weishi      string `json:"weishi"`
}
type Cook struct {
	ClientKey string `json:"ClientKey"`
	Cookies   string `json:"Cookies"`
	Gtk       string `json:"Gtk"`
	Gtk32     string `json:"Gtk32"`
	PSkey     PSkey  `json:"PSkey"`
	Skey      string `json:"Skey"`
}
type Conf struct {
	Enable bool
	GData  map[string]int
}
type Data2 struct {
	Date   string `json:"date"`
	City   string `json:"city"`
	Adcode string `json:"adcode"`
	Min    string `json:"min"`
	Max    string `json:"max"`
	Type   string `json:"type"`
	Air    string `json:"air"`
	Wind   string `json:"wind"`
}
type Weather struct {
	Code int   `json:"code"`
	Data Data2 `json:"data"`
}
type CurrentPacket struct {
	Data      Data   `json:"Data"`
	WebConnID string `json:"WebConnId"`
}
type Data struct {
	Content       string      `json:"Content"`
	FromGroupID   int         `json:"FromGroupId"`
	FromGroupName string      `json:"FromGroupName"`
	FromNickName  string      `json:"FromNickName"`
	FromUserID    int64       `json:"FromUserId"`
	MsgRandom     int         `json:"MsgRandom"`
	MsgSeq        int         `json:"MsgSeq"`
	MsgTime       int         `json:"MsgTime"`
	MsgType       string      `json:"MsgType"`
	RedBaginfo    interface{} `json:"RedBaginfo"`
}
type Message struct {
	CurrentPacket CurrentPacket `json:"CurrentPacket"`
	CurrentQQ     int64         `json:"CurrentQQ"`
}
type Channel struct {
	Channel string `json:"channel"`
}

func (d *Message) IsAdmin() bool {
	_, ok := config.AppConfig.AdminQQMap[d.CurrentPacket.Data.FromUserID]
	return ok
}

type Options struct {
	NotAt bool
}

type Option func(o *Options)

func (d *Message) buildOptions(args ...Option) *Options {
	ret := &Options{}
	for _, v := range args {
		v(ret)
	}
	return ret
}

func (d *Message) SendMessage(msg string, args ...Option) error {
	o := d.buildOptions(args...)
	var err error
	if o.NotAt {
		_, err = SendMsg(d.CurrentPacket.Data.FromGroupID, 0, msg)
	} else {
		_, err = SendMsg(d.CurrentPacket.Data.FromGroupID, d.CurrentPacket.Data.FromUserID, msg)
	}
	return err
}

func (d *Message) CommandMatch(command string) ([]string, bool) {
	reg := regexp.MustCompile(command)
	ret := reg.FindStringSubmatch(d.CurrentPacket.Data.Content)
	return ret, len(ret) > 0
}

func (d *Message) Equal(s string) bool {
	return s == d.CurrentPacket.Data.Content
}

func (d *Message) NotAt() Option {
	return func(o *Options) {
		o.NotAt = true
	}
}

func (d *Data) SendPicByBase64(At int64, Content string, Base64 string) (string, error) {
	//发送图文信息
	tmp := make(map[string]interface{})
	tmp["toUser"] = d.FromGroupID
	tmp["sendToType"] = 2
	tmp["sendMsgType"] = "PicMsg"
	tmp["picBase64Buf"] = ""
	tmp["fileMd5"] = ""
	tmp["picUrl"] = ""
	tmp["picBase64Buf"] = Base64
	tmp["content"] = Content
	tmp["groupid"] = 0
	tmp["atUser"] = At
	tmp1, _ := json.Marshal(tmp)
	resp, err := http.Post("http://"+config.AppConfig.Url+"/v1/LuaApiCaller?funcname=SendMsg&timeout=10&qq="+config.AppConfig.QQ, "application/json", bytes.NewBuffer(tmp1))
	if err != nil {
		return "", nil
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

func (d *Data) SendMsg(At int64, Content string) (string, error) {
	tmp := make(map[string]interface{})
	tmp["toUser"] = d.FromGroupID
	tmp["sendToType"] = 2
	tmp["sendMsgType"] = "TextMsg"
	tmp["picBase64Buf"] = ""
	tmp["fileMd5"] = ""
	tmp["picUrl"] = ""
	tmp["content"] = Content
	tmp["groupid"] = 0
	tmp["atUser"] = At
	tmp1, _ := json.Marshal(tmp)
	resp, err := (http.Post("http://"+config.AppConfig.Url+"/v1/LuaApiCaller?funcname=SendMsg&timeout=10&qq="+config.AppConfig.QQ, "application/json", bytes.NewBuffer(tmp1)))
	if err != nil {
		return "", nil
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

func SendPicByBase64(qqgroup int, At int64, Content string, Base64 string) (string, error) {
	//发送图文信息
	tmp := make(map[string]interface{})
	tmp["toUser"] = qqgroup
	tmp["sendToType"] = 2
	tmp["sendMsgType"] = "PicMsg"
	tmp["picBase64Buf"] = ""
	tmp["fileMd5"] = ""
	tmp["picUrl"] = ""
	tmp["picBase64Buf"] = Base64
	tmp["content"] = Content
	tmp["groupid"] = 0
	tmp["atUser"] = At
	tmp1, _ := json.Marshal(tmp)
	resp, err := http.Post("http://"+config.AppConfig.Url+"/v1/LuaApiCaller?funcname=SendMsg&timeout=10&qq="+config.AppConfig.QQ, "application/json", bytes.NewBuffer(tmp1))
	if err != nil {
		return "", nil
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

func SendXML(qqgroup int, Content string) (string, error) {
	//发送图文信息
	tmp := make(map[string]interface{})
	tmp["toUser"] = qqgroup
	tmp["sendToType"] = 2
	tmp["sendMsgType"] = "XmlMsg"
	tmp["content"] = Content
	tmp["groupid"] = 0
	tmp["atUser"] = 0
	tmp1, _ := json.Marshal(tmp)
	resp, err := http.Post("http://"+config.AppConfig.Url+"/v1/LuaApiCaller?funcname=SendMsg&timeout=10&qq="+config.AppConfig.QQ, "application/json", bytes.NewBuffer(tmp1))
	if err != nil {
		return "", nil
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

func SendPicByUrl(qqgroup int, At int64, Content string, picUrl string) (string, error) {
	//发送图文信息
	tmp := make(map[string]interface{})
	tmp["toUser"] = qqgroup
	tmp["sendToType"] = 2
	tmp["sendMsgType"] = "PicMsg"
	tmp["picBase64Buf"] = ""
	tmp["fileMd5"] = ""
	tmp["picUrl"] = picUrl
	tmp["picBase64Buf"] = ""
	tmp["content"] = Content
	tmp["groupid"] = 0
	tmp["atUser"] = At
	tmp1, _ := json.Marshal(tmp)
	resp, err := http.Post("http://"+config.AppConfig.Url+"/v1/LuaApiCaller?funcname=SendMsg&timeout=10&qq="+config.AppConfig.QQ, "application/json", bytes.NewBuffer(tmp1))
	if err != nil {
		return "", nil
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

func SendMsg(qqgroup int, At int64, Content string) (string, error) {
	tmp := make(map[string]interface{})
	tmp["toUser"] = qqgroup
	tmp["sendToType"] = 2
	tmp["sendMsgType"] = "TextMsg"
	tmp["picBase64Buf"] = ""
	tmp["fileMd5"] = ""
	tmp["picUrl"] = ""
	tmp["content"] = Content
	tmp["groupid"] = 0
	tmp["atUser"] = At
	tmp1, _ := json.Marshal(tmp)
	resp, err := (http.Post("http://"+config.AppConfig.Url+"/v1/LuaApiCaller?funcname=SendMsg&timeout=10&qq="+config.AppConfig.QQ, "application/json", bytes.NewBuffer(tmp1)))
	if err != nil {
		return "", nil
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

func ShutUp(qqgroup int, user, t int64) (string, error) {
	tmp := make(map[string]interface{})
	tmp["ShutUpType"] = 1
	tmp["GroupID"] = qqgroup
	tmp["ShutUid"] = user
	tmp["ShutTime"] = t
	tmp1, _ := json.Marshal(tmp)
	resp, err := (http.Post("http://"+config.AppConfig.Url+"/v1/LuaApiCaller?funcname=ShutUp&timeout=10&qq="+config.AppConfig.QQ, "application/json", bytes.NewBuffer(tmp1)))
	if err != nil {
		return "", nil
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

func SendFriendPicMsg(qq int64, Content string, Base64 string) (string, error) {
	//发送图文信息
	tmp := make(map[string]interface{})
	tmp["toUser"] = qq
	tmp["sendToType"] = 1
	tmp["sendMsgType"] = "PicMsg"
	tmp["picBase64Buf"] = ""
	tmp["fileMd5"] = ""
	tmp["picUrl"] = ""
	tmp["picBase64Buf"] = Base64
	tmp["content"] = Content
	tmp["groupid"] = 0
	tmp["atUser"] = 0
	tmp1, _ := json.Marshal(tmp)
	resp, err := http.Post("http://"+config.AppConfig.Url+"/v1/LuaApiCaller?funcname=SendMsg&timeout=10&qq="+config.AppConfig.QQ, "application/json", bytes.NewBuffer(tmp1))
	if err != nil {
		return "", nil
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if strings.Index(string(body), `"Ret":0`) == -1 {
		return "", errors.New(string(body))
	}
	return string(body), nil
}

func RevokeMsg(qqgroup int, MsgSeq int, MsgRandom int) (string, error) {
	tmp := make(map[string]interface{})
	tmp["GroupID"] = qqgroup
	tmp["MsgSeq"] = MsgSeq
	tmp["MsgRandom"] = MsgRandom
	tmp1, _ := json.Marshal(tmp)
	resp, err := (http.Post("http://"+config.AppConfig.Url+"/v1/LuaApiCaller?funcname=RevokeMsg&timeout=10&qq="+config.AppConfig.QQ, "application/json", bytes.NewBuffer(tmp1)))
	if err != nil {
		return "", nil
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

func GetGroupList() ([]*GroupInfo, error) {
	ret := make([]*GroupInfo, 0)
	token := ""
	for {
		tmp := make(map[string]interface{})
		tmp["NextToken"] = token
		tmp1, _ := json.Marshal(tmp)
		resp, err := (http.Post("http://"+config.AppConfig.Url+"/v1/LuaApiCaller?funcname=GetGroupList&timeout=10&qq="+config.AppConfig.QQ, "application/json", bytes.NewBuffer(tmp1)))
		if err != nil {
			return nil, nil
		}
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		s := struct {
			NextToken string
			TroopList []*GroupInfo
		}{}
		if err := json.Unmarshal(body, &s); err != nil {
			return nil, err
		}
		token = s.NextToken
		ret = append(ret, s.TroopList...)
		if s.NextToken == "" {
			break
		}
	}
	return ret, nil
}

func GetGroupUserList(group int) ([]*GroupMemberInfo, error) {
	ret := make([]*GroupMemberInfo, 0)
	LastUin := 0
	for {
		tmp := make(map[string]interface{})
		tmp["GroupUin"] = group
		tmp["LastUin"] = LastUin
		tmp1, _ := json.Marshal(tmp)
		resp, err := (http.Post("http://"+config.AppConfig.Url+"/v1/LuaApiCaller?funcname=GetGroupUserList&timeout=10&qq="+config.AppConfig.QQ, "application/json", bytes.NewBuffer(tmp1)))
		if err != nil {
			return nil, nil
		}
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		s := struct {
			LastUin    int
			MemberList []*GroupMemberInfo
		}{}
		if err := json.Unmarshal(body, &s); err != nil {
			return nil, err
		}
		LastUin = s.LastUin
		ret = append(ret, s.MemberList...)
		if s.LastUin == 0 {
			break
		}
	}
	return ret, nil
}

func GetGroupAdminList(group int) ([]int64, error) {
	ret := make([]int64, 0)
	if err := db.GetOrSet("qqgroup:admin:list"+strconv.Itoa(group), &ret, func() (interface{}, error) {
		list, err := GetGroupList()
		if err != nil {
			return nil, err
		}
		for _, v := range list {
			if v.GroupId == group {
				ret = append(ret, v.GroupOwner)
			}
		}
		members, err := GetGroupUserList(group)
		if err != nil {
			return nil, err
		}
		for _, v := range members {
			if v.GroupAdmin == 1 {
				ret = append(ret, v.MemberUin)
			}
		}
		return ret, nil
	}, db.WithTTL(time.Hour)); err != nil {
		return nil, err
	}
	return ret, nil
}

func IsAdmin(group int, user int64) (bool, error) {
	list, err := GetGroupAdminList(group)
	if err != nil {
		return false, err
	}
	for _, v := range list {
		if v == user {
			return true, nil
		}
	}
	return false, nil
}

func ModifyGroupCard(qqgroup int, UserID int64, NewNick string) (string, error) {
	tmp := make(map[string]interface{})
	tmp["GroupID"] = qqgroup
	tmp["UserID"] = UserID
	tmp["NewNick"] = NewNick
	tmp1, _ := json.Marshal(tmp)
	resp, err := (http.Post("http://"+config.AppConfig.Url+"/v1/LuaApiCaller?funcname=ModifyGroupCard&timeout=10&qq="+config.AppConfig.QQ, "application/json", bytes.NewBuffer(tmp1)))
	if err != nil {
		return "", nil
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}
