package model

import (
	"bytes"
	"encoding/json"
	"github.com/CodFrm/iotqq-plugins/config"
	"io/ioutil"
	"net/http"
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
