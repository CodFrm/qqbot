package live

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/CodFrm/qqbot/cqhttp"
	transition "github.com/CodFrm/qqbot/live/ffmpeg"
)

var lives = make(map[string]*live)

func AddLive(guild, channel, user int64, url, secret string) error {
	live := newLive(guild, channel, user, url, secret)

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

func PlayQueue(guild, channel, user int64, filename string) error {
	live, ok := lives[fmt.Sprintf("%d:%d:%d", guild, channel, user)]
	if !ok {
		return errors.New("没有推流权限")
	}
	return live.PlayQueue(filename)
}

var trProgress float32
var inProgress bool

func ToFlv(guild, channel, user int64, filename string) error {
	if inProgress {
		return fmt.Errorf("上个视频还在转码中: %.2f", trProgress)
	}
	_, ok := lives[fmt.Sprintf("%d:%d:%d", guild, channel, user)]
	if !ok {
		return errors.New("没有推流权限")
	}
	i := strings.LastIndex(filename, ".")
	if i == -1 {
		return errors.New("错误的文件名")
	}
	info, _ := os.Stat("./data/live/flv/" + filename[:i] + ".flv")
	if info != nil {
		return errors.New("已转码过了")
	}

	go func() {
		inProgress = true
		progress := make(chan float32)
		ctx, cancel := context.WithCancel(context.Background())
		defer func() {
			cancel()
			trProgress = 0
			inProgress = false
		}()
		go func() {
			for {
				select {
				case <-ctx.Done():
				case trProgress = <-progress:
				}
			}
		}()
		err := transition.ToFlv("./data/live/source/"+filename, "./data/live/flv/"+filename[:i]+".flv", progress)
		if err != nil {
			cqhttp.SendGuildChannelMsg(guild, channel, "[CQ:at,qq="+strconv.FormatInt(user, 10)+"] "+"转码失败: "+err.Error())
		}
	}()
	return nil
}

func TrProgress() string {
	if trProgress == 0 {
		return "暂无转码任务"
	}
	return fmt.Sprintf("%.2f", trProgress)
}
