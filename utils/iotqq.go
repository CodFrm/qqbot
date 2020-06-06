package utils

import (
	"bytes"
	"encoding/json"
	"github.com/CodFrm/iotqq-plugins/config"
	"io/ioutil"
	"net/http"
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
