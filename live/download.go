package live

import (
	"errors"
	"fmt"

	"github.com/CodFrm/qqbot/live/aria2"
)

func DownloadToSource(guild, channel, user int64, url, platform string) error {
	_, ok := lives[fmt.Sprintf("%d:%d:%d", guild, channel, user)]
	if !ok {
		return errors.New("没有推流权限")
	}
	_, err := aria2.Download(url, platform, "./data/live/source")
	return err
}

func DownloadToFlv(guild, channel, user int64, url, platform string) error {
	_, ok := lives[fmt.Sprintf("%d:%d:%d", guild, channel, user)]
	if !ok {
		return errors.New("没有推流权限")
	}
	_, err := aria2.Download(url, platform, "./data/live/flv")
	return err
}
