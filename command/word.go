package command

import (
	"bytes"
	"github.com/CodFrm/iotqq-plugins/db"
	"github.com/CodFrm/iotqq-plugins/utils"
	"github.com/otiai10/gosseract/v2"
	"image"
	"regexp"
	"strconv"
)

func IsWordGroup(group int) bool {
	return db.Redis.HExists("word:group", strconv.Itoa(group)).Val()
}

func IsWordOk(group int, qq int64, msg string) (string, error) {
	flag := false
	str := utils.RegexMatch(msg, "brief=\\\\\"\\[分享\\].+?百词斩.+?")
	if len(str) > 0 {
		//百词斩
		flag = true
	}
	if flag {
		return Sign(group, qq)
	}
	return "", nil
}

func IsWordImage(group int, qq int64, pic []byte) (string, error) {
	r := bytes.NewBuffer(pic)
	_, s, err := image.Decode(r)
	if err != nil {
		return "", nil
	}
	if s != "jpg" && s != "png" && s != "jpeg" {
		return "", nil
	}
	if len(pic) <= 1024*200 {
		return "", nil
	}
	client := gosseract.NewClient()
	defer client.Close()
	err = client.SetImageFromBytes(pic)
	if err != nil {
		return "", err
	}
	text, err := client.Text()
	if err != nil {
		return "", err
	}
	reg := regexp.MustCompile("[a-zA-Z]{3,}")
	list := reg.FindAllString(text, -1)
	strlen := 0
	for _, v := range list {
		strlen += len(v)
	}
	if strlen > 300 {
		return Sign(group, qq)
	}
	return "", nil
}
