package alimama

import (
	"github.com/CodFrm/qqbot/config"
	"github.com/CodFrm/qqbot/utils"
	"net/url"
)

func ShortUrl(u string) string {
	resp, err := utils.HttpGet("http://api.suolink.cn/api.htm?domain=mrw.so&url="+url.QueryEscape(u)+"&key="+config.AppConfig.Urlkey, nil, nil)
	if err != nil {
		return u
	}
	return string(resp)
}
