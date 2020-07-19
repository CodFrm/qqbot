package command

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/CodFrm/iotqq-plugins/db"
	"github.com/CodFrm/iotqq-plugins/utils"
	"github.com/CodFrm/iotqq-plugins/utils/iotqq"
	"github.com/otiai10/gosseract/v2"
	"image"
	"regexp"
	"strconv"
	"strings"
	"time"
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
	str = utils.RegexMatch(msg, "source name=\\\\\"不背单词\\\\\" icon=\\\\\"")
	if len(str) > 0 {
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

func BindWebsite(group int, qq int64, content string) (string, error) {
	m := utils.RegexMatch(content, "https:\\/\\/web\\.shanbay\\.com\\/web\\/wechat\\/calendar\\/\\?user_id=(.+?)\\\\")
	if len(m) > 0 {
		db.Redis.HSet("word:site:"+strconv.Itoa(group), qq, "https://apiv3.shanbay.com/uc/checkin/calendar/dates?user_id="+m[1])
		shanbayScan(group, qq, "https://apiv3.shanbay.com/uc/checkin/calendar/dates?user_id="+m[1], time.Now())
		return "扇贝绑定成功", nil
	}
	m = utils.RegexMatch(content, "https:\\/\\/www\\.maimemo\\.com\\/share\\/page\\?(.*?)\\\\\\\"")
	if len(m) > 0 {
		db.Redis.HSet("word:site:"+strconv.Itoa(group), qq, "https://www.maimemo.com/share/page?"+strings.ReplaceAll(m[1], "\\u0026amp;", "&"))
		momoScan(group, qq, "https://www.maimemo.com/share/page?"+strings.ReplaceAll(m[1], "\\u0026amp;", "&"), time.Now())
		return "墨墨绑定成功", nil
	}
	//m = utils.RegexMatch(content, "https:\\/\\/www\\.maimemo\\.com\\/share\\/page\\?(.*?)\\\\\\\"")
	//if len(m) > 0 {
	//	db.Redis.HSet("word:site:"+strconv.Itoa(group), qq, "https://www.maimemo.com/share/page?"+strings.ReplaceAll(m[1], "\\u0026amp;", "&"))
	//	return "百词斩绑定成功", nil
	//}
	return "", nil
}

func scanSite() {
	day := time.Now().Add(-time.Hour * 24).Format("2006:01:02")
	key := "sign:group:record:" + day
	list := db.Redis.HGetAll(key).Val()
	for group := range list {
		qqs := db.Redis.HGetAll("word:site:" + group).Val()
		for qq, val := range qqs {
			igroup, iqq := utils.StringToInt(group), utils.StringToInt64(qq)
			scanSignalSite(igroup, iqq, val)
		}
	}
}

type shanbayLog struct {
	Logs []struct {
		Date string `json:"date"`
	} `json:"logs"`
}

//
//func SignByWord(qqgroup int, qq int64) (string, error) {
//	group := strconv.Itoa(qqgroup)
//	s := db.Redis.HGet("word:site:"+group, strconv.FormatInt(qq, 10)).Val()
//	if s != "" {
//		if ok, err := scanSignalSite(qqgroup, qq, s); err != nil {
//			return "", err
//		} else if ok {
//			return "", nil
//		}
//	}
//	return Sign(qqgroup, qq)
//}

func scanSignalSite(igroup int, iqq int64, site string) (bool, error) {
	if IsSign(igroup, iqq) {
		return false, errors.New("今天签过到了")
	}
	if ok, err := iotqq.IsInGroup(igroup, iqq); err != nil {
		return false, err
	} else if !ok {
		delSign(igroup, iqq)
		return false, errors.New("不在群里?")
	}
	if strings.Index(site, "apiv3.shanbay.com") != -1 {
		if shanbayScan(igroup, iqq, site, time.Now()) {
			return true, nil
		}
	} else if strings.Index(site, "www.maimemo.com") != -1 {
		if momoScan(igroup, iqq, site, time.Now()) {
			return true, nil
		}
	}
	return false, nil
}

func momoScan(group int, qq int64, url string, day time.Time) bool {
	ret, err := utils.HttpGet(url, nil, nil)
	if err != nil {
		return false
	}
	m := utils.RegexMatch(string(ret), "<p>学习天数：<span>(\\d+)</span>天</p>")
	if len(m) > 0 {
		d := db.Redis.HGet("sign:group:day:"+strconv.Itoa(group), strconv.FormatInt(qq, 10)).Val()
		if d != m[1] {
			db.Redis.HSet("sign:group:day:"+strconv.Itoa(group), strconv.FormatInt(qq, 10), m[1])
			if s, _ := Sign(group, qq); s != "" {
				iotqq.QueueSendMsg(group, qq, "墨墨检测成功,自动打卡,"+s)
				return true
			}
		}
	}
	return false
}

func shanbayScan(group int, qq int64, url string, day time.Time) bool {
	ret, err := utils.HttpGet(url+"&start_date="+day.Add(-time.Hour*72).Format("2006-01-02")+
		"&end_date="+day.Add(time.Hour*72).Format("2006-01-02"), nil, nil)
	if err != nil {
		return false
	}
	logs := &shanbayLog{}
	if err := json.Unmarshal(ret, logs); err != nil {
		return false
	}
	for _, v := range logs.Logs {
		if v.Date == day.Format("2006-01-02") {
			if s, _ := Sign(group, qq); s != "" {
				iotqq.QueueSendMsg(group, qq, "扇贝检测成功,自动打卡,"+s)
				return true
			}
			break
		}
	}
	return false
}
