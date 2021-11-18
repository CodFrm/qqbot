package live

import (
	"errors"
	"fmt"

	"github.com/golang/glog"
	rtmp "github.com/zhangpeihao/gortmp"
)

var lives = make(map[string]*live)

func AddLive(guild, channel, user int64, url, secret string) error {
	live := newLive(guild, channel, user)

	client, err := rtmp.Dial(url, live, 100)
	if err != nil {
		return err
	}
	if err := client.Connect(); err != nil {
		return err
	}

	go func() {
		for {
			select {
			case stream := <-live.createStreamChan:
				stream.Attach(live)
				err = stream.Publish(secret, "live")
				if err != nil {
					glog.Errorf("Publish error: %s", err.Error())
				}
			}
		}
	}()
	lives[fmt.Sprintf("%d:%d:%d", guild, channel, user)] = live
	return nil
}

func Play(guild, channel, user int64, filename string) error {
	live, ok := lives[fmt.Sprintf("%d:%d:%d", guild, channel, user)]
	if !ok {
		return errors.New("没有推流权限")
	}
	return live.Play(filename)
}
