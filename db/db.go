package db

import (
	"fmt"
	"github.com/CodFrm/iotqq-plugins/config"
	goRedis "github.com/go-redis/redis/v7"
)

var Redis *goRedis.Client

func Init() error {
	Redis = goRedis.NewClient(&goRedis.Options{
		Addr:     config.AppConfig.Redis.Addr,
		Password: config.AppConfig.Redis.Password,
		DB:       config.AppConfig.Redis.DB,
	})
	if _, err := Redis.Ping().Result(); err != nil {
		return fmt.Errorf("redis open error: %v", err)
	}
	return nil
}
