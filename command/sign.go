package command

import (
	"encoding/json"
	"errors"
	"github.com/CodFrm/iotqq-plugins/db"
	"github.com/CodFrm/iotqq-plugins/utils"
	"github.com/CodFrm/iotqq-plugins/utils/iotqq"
	"github.com/go-redis/redis/v7"
	"github.com/pkumza/numcn"
	"github.com/robfig/cron/v3"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var RewardsMap = map[string]func(group int, qq int64, rewards bool, day time.Time, continuous int, args ...string){
	"设置名片": rewardGroupName, "nmsl": rewardNmsl,
	"nmsl单词特供版": rewardNmsl2, "踢出本群": rewardKick,
	"温柔词典": rewardRainbowFart,
}

var nmslEnglish []string
var rainbowFart []string

func SignInit() {
	c := cron.New(cron.WithSeconds())
	c.AddFunc("0 0 7 * * ?", everyDay)
	c.AddFunc("0 0/15 * * * ?", scanSite)
	c.Start()
	rand.Seed(time.Now().UnixNano())
	nmslEnglish = strings.Split(`你妈逼你今天学习了吗？废物
连一个单词都背不下来，天天躺着做你妈春秋大梦呢？
我家狗都会说英文，你竟然不会？
牌牌琦都会上Youtube打出来Chinese Kongfu Yao，而你28个英文字符都拼不全。
单词都不背，天天活你妈的有个屁意思。
你说只要你努力，全世界都会为你让步，其实你狗屁不是，10个单词都背不下来。
现在的大学生真的不行，几十个单词都背不下来还天天嘲讽我们大专生。
当年我认识的一个喜欢吃骨灰拌饭的妹妹都能一天背十个单词，再看看你这个废物？
老子拿脚踩一脚四级卷子考出来的都比你分高，你也能挺个逼脸不背单词？
你说你女神喜欢洋人是婊子，其实你不知道她跟外国人处对象是为了学英文，而你只会说卧槽。
你笑印度人说英文有股咖喱味，印度人笑你连用英语说咖喱都不会说。
你笑特朗普是傻逼，却不知道人家说着你这辈子都学不会的语言。你也配说他？`, "\n")
	rainbowFart = strings.Split(`我喜欢背单词，因为我喜欢你，而你喜欢背单词，所以请让我一直喜欢背单词好吗？
你知道吗？我朋友一直很奇怪，为什么我高考58分却能过四级，只有我知道，是因为我不想你看低我
我学英语的唯一动力，就是希望写一首英文情诗给你，所以答应我，请你好好背单词，可以看懂好嘛？
你知道我为什么想让你学英语么？因为只有你学英语的时候，才会对我说i love you
为什么我想让你背单词，因为我想有一天，我问你我好孤独怎么说，你可以对我说i love you
你不可以这么懒惰的！再不背单词我就叫警察叔叔给你抓走了哦，那样你就再也看不到我了！
我当初没有跟心爱的人考到一个学校，最大的鸿沟就是我英语58，而她110，所以不要再重复我的悲剧，好嘛
如果你每天都更努力一点，那我就能喜欢你多一点！所以答应我继续背单词好不好嘛！
你说如果我们一起考过六级，就让我跟你在一起，所以我一直在努力，因为放弃你我做不到，请你也不要放弃我好不好
我背了好多好多单词，只希望有一天能漂洋过海去看他，在过年的时候跟他说一声happy new year 希望你也努力，不要将来有一天来一句过年好 二妞
阿巴阿巴阿巴阿巴阿巴阿巴阿巴阿巴阿巴阿巴阿巴`, "\n")
}

func Sign(qqgroup int, qq int64) (string, error) {
	key := "sign:day:" + strconv.Itoa(qqgroup) + ":"
	val, err := db.Redis.HGet(key+time.Now().Format("2006:01:02"), strconv.FormatInt(qq, 10)).Result()
	if err != nil && err != redis.Nil {
		return "", err
	}
	if val == "1" {
		return "", errors.New("今天签过到了")
	}
	autoAddReward(strconv.Itoa(qqgroup), qq)
	if err := db.Redis.HSet(key+time.Now().Format("2006:01:02"), qq, "1").Err(); err != nil {
		return "", err
	}
	continuous := 1
	if val, err := db.Redis.HGet(key+time.Now().Add(-time.Hour*24).Format("2006:01:02"), strconv.FormatInt(qq, 10)).Result(); err != nil && err != redis.Nil {
		return "", err
	} else if val == "1" {
		continuous = int(db.Redis.HIncrBy("sign:record:"+strconv.Itoa(qqgroup), strconv.FormatInt(qq, 10), 1).Val())
	} else {
		continuous = 1
		db.Redis.HSet("sign:record:"+strconv.Itoa(qqgroup), qq, 1)
	}
	db.Redis.Expire(key+time.Now().Format("2006:01:02"), time.Hour*72)
	go execRewards(qqgroup, qq, true, time.Now(), continuous)
	db.Redis.HSet("sign:group:record:"+time.Now().Format("2006:01:02"), strconv.Itoa(qqgroup), "1")
	db.Redis.HSet("sign:end:record:"+strconv.Itoa(qqgroup), qq, time.Now().Format("2006:01:02"))
	return "打卡成功,你连续打卡了" + numcn.EncodeFromInt64(int64(continuous)) + "天", nil
}

func SetContinuousDay(qqgroup int, qq int64, day int) error {
	return db.Redis.HSet("sign:record:"+strconv.Itoa(qqgroup), strconv.FormatInt(qq, 10), strconv.Itoa(day)).Err()
}

func IsSign(qqgroup int, qq int64) bool {
	key := "sign:day:" + strconv.Itoa(qqgroup) + ":"
	val, err := db.Redis.HGet(key+time.Now().Format("2006:01:02"), strconv.FormatInt(qq, 10)).Result()
	if err != nil && err != redis.Nil {
		return false
	}
	if val == "1" {
		return true
	}
	return false
}

func autoAddReward(group string, qq int64) {
	list, _ := GetRewards("group"+group, 8888)
	for _, v := range list {
		list, _ := GetRewards(group, qq)
		flag := false
		for _, v2 := range list {
			if v2.Command == v.Command {
				flag = true
				break
			}
		}
		if !flag {
			SetRewards(group, qq, false, v.Command, v.Args...)
		}
	}
}

type Reward struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
}

func everyDay() {
	day := time.Now().Add(-time.Hour * 24).Format("2006:01:02")
	key := "sign:group:record:" + day
	list := db.Redis.HGetAll(key).Val()
	for group := range list {
		qqs := db.Redis.HGetAll("sign:end:record:" + group).Val()
		for qq, val := range qqs {
			igroup, iqq := utils.StringToInt(group), utils.StringToInt64(qq)
			if ok, err := iotqq.IsInGroup(igroup, iqq); err != nil {
				continue
			} else if !ok {
				delSign(igroup, iqq)
				continue
			}
			if val != day && val != time.Now().Format("2006:01:02") {
				//惩罚
				t, _ := time.Parse("2006:01:02", val)
				go execRewards(utils.StringToInt(group), utils.StringToInt64(qq), false, t, 0)
			}
		}
	}
}

func AdminGroupReward(qqgroup string, rm bool, command string, args ...string) error {
	if err := SetRewards("group"+qqgroup, 8888, rm, command, args...); err != nil {
		return err
	}
	list := db.Redis.HGetAll("sign:end:record:" + qqgroup).Val()
	for k := range list {
		autoAddReward(qqgroup, utils.StringToInt64(k))
	}
	return nil
}

func SetRewards(qqgroup string, qq int64, rm bool, command string, args ...string) error {
	rs, err := GetRewards(qqgroup, qq)
	if err != nil {
		return err
	}
	flag := false
	for k, v := range rs {
		if v.Command == command {
			if rm == true {
				if k == 0 {
					rs = rs[1:]
				} else if k == len(rs)-1 {
					rs = rs[0 : len(rs)-1]
				} else {
					rs = append(rs[k:], rs[k+1:]...)
				}
			} else {
				v.Args = args
			}
			flag = true
			break
		}
	}
	if _, ok := RewardsMap[command]; !ok {
		return errors.New("不存在的奖惩方案")
	}
	if !flag && rm == false {
		rs = append(rs, &Reward{
			Command: command,
			Args:    args,
		})
	}
	s, err := json.Marshal(rs)
	if err != nil {
		return err
	}
	key := "sign:rewards:" + qqgroup
	return db.Redis.HSet(key, strconv.FormatInt(qq, 10), s).Err()
}

func GetRewards(qqgroup string, qq int64) ([]*Reward, error) {
	key := "sign:rewards:" + qqgroup
	val, err := db.Redis.HGet(key, strconv.FormatInt(qq, 10)).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	if val == "" {
		return nil, nil
	}
	rs := make([]*Reward, 0)
	if err := json.Unmarshal([]byte(val), &rs); err != nil {
		return nil, err
	}
	return rs, nil
}

func execRewards(qqgroup int, qq int64, rewards bool, day time.Time, continuous int) {
	list, _ := GetRewards(strconv.Itoa(qqgroup), qq)
	for _, v := range list {
		f := RewardsMap[v.Command]
		if f != nil {
			f(qqgroup, qq, rewards, day, continuous, v.Args...)
		}
	}
}

func rewardGroupName(group int, qq int64, rewards bool, day time.Time, continuous int, args ...string) {
	if len(args) < 1 {
		return
	}
	if continuous < 1 {
		continuous = 1
	}
	if !rewards {
		return
	}
	s := strings.Join(args, " ")
	s = strings.Replace(s, "N", numcn.EncodeFromInt64(int64(continuous)), 1)
	iotqq.ModifyGroupCard(group, qq, s)
	time.Sleep(time.Second * 2)
	iotqq.SendMsg(group, qq, "奖励你新id:"+s)
}

func rewardNmsl(group int, qq int64, rewards bool, day time.Time, continuous int, args ...string) {
	if continuous > 0 || rewards {
		return
	}
	iotqq.QueueSendMsg(group, qq, utils.Nmsl())
}

//英语特供版
func rewardNmsl2(group int, qq int64, rewards bool, day time.Time, continuous int, args ...string) {
	if rewards {
		return
	}
	if rand.Intn(100) < 2 {
		str := utils.FileBase64("./data/img/3.jpg")
		iotqq.SendPicByBase64(group, qq, "", str)
		return
	}
	iotqq.QueueSendMsg(group, qq, strings.ReplaceAll(nmslEnglish[rand.Intn(len(nmslEnglish))], "妈", "🐴"))
}

func rewardKick(group int, qq int64, rewards bool, day time.Time, continuous int, args ...string) {
	if rewards {
		return
	}
	t := time.Now().Sub(day)
	d := t.Hours() / 24
	if d >= 3 {
		iotqq.QueueSendMsg(group, qq, "超过3天未打卡,将自动移除本群")
		delSign(group, qq)
		iotqq.Kick(group, qq)
		return
	}
	iotqq.QueueSendMsg(group, qq, "提示:超过3天未打卡,将自动移除本群")
}

func delSign(group int, qq int64) {
	db.Redis.HDel("sign:end:record:"+strconv.Itoa(group), strconv.FormatInt(qq, 10))
}

func rewardRainbowFart(group int, qq int64, rewards bool, day time.Time, continuous int, args ...string) {
	iotqq.QueueSendMsg(group, qq, rainbowFart[rand.Intn(len(rainbowFart))])
}
