package command

import (
	"errors"
	"github.com/CodFrm/iotqq-plugins/db"
	"strconv"
	"time"
)

func BlackList(user, remove, t string) error {
	d := time.Duration(0)
	if t != "" {
		tmp, _ := strconv.ParseInt(t, 10, 64)
		d = time.Duration(tmp)
	}
	return db.Redis.Set("blacklist:"+user, remove, d).Err()
}

func IsBlackList(user string) error {
	if val := db.Redis.Get("blacklist:" + user).Val(); val == "1" {
		return errors.New("黑名单中")
	}
	return nil
}
