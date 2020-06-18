package command

import (
	"encoding/json"
	"errors"
	"github.com/CodFrm/iotqq-plugins/db"
	"github.com/CodFrm/iotqq-plugins/utils/iotqq"
	"github.com/go-redis/redis/v7"
	"github.com/pkumza/numcn"
	"strconv"
	"strings"
	"time"
)

var rewardsMap = map[string]func(group int, qq int64, rewards bool, continuous int, args ...string){
	"设置名片": rewardGroupName,
}

func Sign(qqgroup int, qq int64) error {
	key := "sign:day:" + strconv.Itoa(qqgroup) + ":"
	val, err := db.Redis.HGet(key+time.Now().Format("2006:01:02"), strconv.FormatInt(qq, 10)).Result()
	if err != nil && err != redis.Nil {
		return err
	}
	if val == "1" {
		return errors.New("今天签过到了")
	}
	if err := db.Redis.HSet(key+time.Now().Format("2006:01:02"), qq, "1").Err(); err != nil {
		return err
	}
	continuous := 1
	if val, err := db.Redis.HGet(key+time.Now().Add(-time.Hour*24).Format("2006:01:02"), strconv.FormatInt(qq, 10)).Result(); err != nil && err != redis.Nil {
		return err
	} else if val == "1" {
		continuous = int(db.Redis.HIncrBy("sign:record:"+strconv.Itoa(qqgroup), strconv.FormatInt(qq, 10), 1).Val())
	} else {
		continuous = 1
		db.Redis.HSet("sign:record:"+strconv.Itoa(qqgroup), qq, 1)
	}
	db.Redis.Expire(key, time.Hour*72)
	go execRewards(qqgroup, qq, true, continuous)
	return nil
}

type Reward struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
}

func SetRewards(qqgroup int, qq int64, rm bool, command string, args ...string) error {
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
	if _, ok := rewardsMap[command]; !ok {
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
	key := "sign:rewards:" + strconv.Itoa(qqgroup)
	return db.Redis.HSet(key, strconv.FormatInt(qq, 10), s).Err()
}

func GetRewards(qqgroup int, qq int64) ([]*Reward, error) {
	key := "sign:rewards:" + strconv.Itoa(qqgroup)
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

func execRewards(qqgroup int, qq int64, rewards bool, continuous int) {
	list, _ := GetRewards(qqgroup, qq)
	for _, v := range list {
		f := rewardsMap[v.Command]
		if f != nil {
			f(qqgroup, qq, rewards, continuous, v.Args...)
		}
	}
}

func rewardGroupName(group int, qq int64, rewards bool, continuous int, args ...string) {
	if len(args) < 1 {
		return
	}
	s := strings.Join(args, " ")
	s = strings.Replace(s, "N", numcn.EncodeFromInt64(int64(continuous)), 1)
	iotqq.ModifyGroupCard(group, qq, s)
	if rewards {
		time.Sleep(time.Second * 2)
		iotqq.SendMsg(group, qq, "奖励你新id:"+s)
	}
}
