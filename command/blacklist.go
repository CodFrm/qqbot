package command

import (
	"errors"
	"github.com/CodFrm/iotqq-plugins/db"
)

func BlackList(user, remove string) error {
	return db.Redis.Set("blacklist:"+user, remove, 0).Err()
}

func IsBlackList(user string) error {
	if val := db.Redis.Get("blacklist:" + user).Val(); val == "1" {
		return errors.New("黑名单中")
	}
	return nil
}
