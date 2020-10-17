package alimama

import (
	"github.com/CodFrm/iotqq-plugins/config"
	"github.com/CodFrm/iotqq-plugins/utils"
	"net/url"
)

func ShortUrl(u string) string {
	resp, err := utils.HttpGet("http://api.suolink.cn/api.htm?domain=mtw.so&url="+url.QueryEscape(u)+"&key="+config.AppConfig.Urlkey, nil, nil)
	if err != nil {
		return u
	}
	return string(resp)
}
