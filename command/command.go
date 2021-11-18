package command

import (
	"context"
	"github.com/CodFrm/qqbot/config"
	"github.com/CodFrm/qqbot/db"
	"github.com/mzz2017/shadowsocksR/client"
	pxy "github.com/nadoo/glider/proxy"
	"github.com/robfig/cron/v3"
	"net"
	"time"
)

var proxy func(ctx context.Context, network, addr string) (conn net.Conn, err error)

func Init() error {
	dia, err := client.NewSSRDialer(config.AppConfig.Ssr, pxy.Default)
	if err != nil {
		return err
	}
	proxy = func(ctx context.Context, network, addr string) (conn net.Conn, err error) {
		return dia.Dial(network, addr)
	}
	scenes = db.NewScenes()
	SignInit()
	c := cron.New(cron.WithSeconds())
	c.AddFunc("0 10 8 * * ?", cronGenWordCloud)
	c.Start()
	return nil
}

func CleanCache(key string) error {
	if key != "" {
		return delKeys(key + "*")
	}
	if err := delKeys("pixivList*"); err != nil {
		return err
	}
	if err := delKeys("pixivlist*"); err != nil {
		return err
	}
	if err := scanKey("pixivTag*", func(key string) error {
		db.Redis.Get(key)
		ret := &struct {
			db.StringCache
			db.IntCache
		}{}
		if err := db.Get(key, ret); err != nil {
			return err
		}
		ret.Int = 0
		db.Set(key, ret, db.WithTTL(time.Second*86400*7))
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func delKeys(key string) error {
	if err := scanKey(key, func(key string) error {
		return db.Redis.Del(key).Err()
	}); err != nil {
		return err
	}
	return nil
}

func scanKey(key string, deal func(key string) error) error {
	var cur uint64
	var keys []string
	var err error
	for {
		keys, cur, err = db.Redis.Scan(cur, key, 100).Result()
		if err != nil {
			return err
		}
		for _, v := range keys {
			if err := deal(v); err != nil {
				return err
			}
		}
		if cur == 0 {
			break
		}
	}
	return nil
}
