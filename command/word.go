package command

import (
	"bytes"
	"encoding/json"
	"errors"
	"image"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/CodFrm/qqbot/cqhttp"
	"github.com/CodFrm/qqbot/db"
	"github.com/CodFrm/qqbot/utils"
	"github.com/otiai10/gosseract/v2"
)

func IsWordGroup(group, channel int64) bool {
	return db.Redis.HExists("word:group", strconv.FormatInt(group, 10)+":"+strconv.FormatInt(channel, 10)).Val()
}

func IsWordOk(group, channel, qq int64, msg string) (string, error) {
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
		return Sign(group, channel, qq)
	}
	return "", nil
}

func IsWordImage(group, channel, qq int64, pic []byte) (string, error) {
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
		return Sign(group, channel, qq)
	}
	return "", nil
}

func BindWebsite(group, channel, qq int64, content string) (string, error) {
	m := utils.RegexMatch(content, "https:\\/\\/web\\.shanbay\\.com\\/web\\/wechat\\/calendar\\/\\?user_id=(.+?)\\\\")
	if len(m) > 0 {
		db.Redis.HSet("word:site:"+strconv.FormatInt(group, 10)+":"+strconv.FormatInt(channel, 10), qq, "https://apiv3.shanbay.com/uc/checkin/calendar/dates?user_id="+m[1])
		shanbayScan(group, channel, qq, "https://apiv3.shanbay.com/uc/checkin/calendar/dates?user_id="+m[1], time.Now())
		return "扇贝绑定成功", nil
	}
	m = utils.RegexMatch(content, "https:\\/\\/www\\.maimemo\\.com\\/share\\/page\\?(.*?)\\\\\\\"")
	if len(m) > 0 {
		db.Redis.HSet("word:site:"+strconv.FormatInt(group, 10)+":"+strconv.FormatInt(channel, 10), qq, "https://www.maimemo.com/share/page?"+strings.ReplaceAll(m[1], "\\u0026amp;", "&"))
		momoScan(group, channel, qq, "https://www.maimemo.com/share/page?"+strings.ReplaceAll(m[1], "\\u0026amp;", "&"), time.Now())
		return "墨墨绑定成功", nil
	}
	//m = utils.RegexMatch(content, "https:\\/\\/www\\.maimemo\\.com\\/share\\/page\\?(.*?)\\\\\\\"")
	//if len(m) > 0 {
	//	db.Redis.HSet("word:site:"+strconv.FormatInt(group, 10)+":"+strconv.FormatInt(channel, 10), qq, "https://www.maimemo.com/share/page?"+strings.ReplaceAll(m[1], "\\u0026amp;", "&"))
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
			s := strings.Split(group, ":")
			if len(s) != 2 {
				continue
			}
			igroup, channel, iqq := utils.StringToInt64(s[0]), utils.StringToInt64(s[1]), utils.StringToInt64(qq)
			scanSignalSite(igroup, channel, iqq, val)
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

func scanSignalSite(igroup, channel, iqq int64, site string) (bool, error) {
	if IsSign(igroup, channel, iqq) {
		return false, errors.New("今天签过到了")
	}
	//if ok, err := iotqq.IsInGroup(igroup, channel, iqq); err != nil {
	//	return false, err
	//} else if !ok {
	//	delSign(igroup, channel, iqq)
	//	return false, errors.New("不在群里?")
	//}
	if strings.Index(site, "apiv3.shanbay.com") != -1 {
		if shanbayScan(igroup, channel, iqq, site, time.Now()) {
			return true, nil
		}
	} else if strings.Index(site, "www.maimemo.com") != -1 {
		if momoScan(igroup, channel, iqq, site, time.Now()) {
			return true, nil
		}
	}
	return false, nil
}

func momoScan(group, channel, qq int64, url string, day time.Time) bool {
	ret, err := utils.HttpGet(url, nil, nil)
	if err != nil {
		return false
	}
	m := utils.RegexMatch(string(ret), "<p>学习天数：<span>(\\d+)</span>天</p>")
	if len(m) > 0 {
		d := db.Redis.HGet("sign:group:day:"+strconv.FormatInt(group, 10)+":"+strconv.FormatInt(channel, 10), strconv.FormatInt(qq, 10)).Val()
		if d != m[1] {
			db.Redis.HSet("sign:group:day:"+strconv.FormatInt(group, 10)+":"+strconv.FormatInt(channel, 10), strconv.FormatInt(qq, 10), m[1])
			if s, _ := Sign(group, channel, qq); s != "" {
				sendmsg(group, channel, qq, "墨墨检测成功,自动打卡,"+s)
				return true
			}
		}
	}
	return false
}

func shanbayScan(group, channel, qq int64, url string, day time.Time) bool {
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
			if s, _ := Sign(group, channel, qq); s != "" {
				sendmsg(group, channel, qq, "扇贝检测成功,自动打卡,"+s)
				return true
			}
			break
		}
	}
	return false
}

func sendmsg(group, channel, qq int64, text string) {
	if channel == 0 {
		//TODO:处理群
	} else {
		cqhttp.SendGuildChannelMsg(group, channel, "[CQ:at,qq="+strconv.FormatInt(qq, 10)+"]"+text)
	}
}
