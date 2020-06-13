package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/CodFrm/iotqq-plugins/config"
	"io/ioutil"
	"net/http"
	"strings"
)

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
