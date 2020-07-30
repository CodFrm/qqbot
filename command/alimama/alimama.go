package alimama

import (
	"github.com/CodFrm/iotqq-plugins/config"
	"github.com/CodFrm/iotqq-plugins/db"
	"github.com/CodFrm/iotqq-plugins/utils"
	"github.com/CodFrm/iotqq-plugins/utils/iotqq"
	"github.com/CodFrm/iotqq-plugins/utils/taobaoopen"
	"github.com/robfig/cron/v3"
	"net/url"
	"strconv"
)

var tb *taobaoopen.Taobao

func Init() error {
	c := cron.New(cron.WithSeconds())
	c.AddFunc("0 30 11 * * ?", notice)
	c.AddFunc("0 30 17 * * ?", notice)
	c.Start()
	tb = taobaoopen.NewTaobao(config.AppConfig.Taobao)
	return nil
}

func notice() {
	list, err := iotqq.GetGroupList()
	if err != nil {
		return
	}
	for _, v := range list {
		iotqq.QueueSendMsg(v.GroupId, 0, "快到饭点了,来一份外卖吧~\nhttps://sourl.cn/FhPLTD\n复制这条信息，$nH3n1zNqDip$，到【手机淘宝】即可查看."+
			"美团可使用此链接:https://sourl.cn/Kvz8Hk\n"+
			"后续将会通过QQ红包提供返现功能(预计外卖返现5%,购物0-10%不等)")
	}
}

func Search(keyword string) (string, error) {
	list, err := tb.MaterialSearch(keyword)
	if err != nil {
		return "", err
	}
	ret := &db.StringCache{}
	if err := db.GetOrSet("alimama:search:"+keyword, ret, func() (interface{}, error) {
		if s, err := GenCopywriting(list); err != nil {
			return nil, err
		} else {
			return &db.StringCache{String: s}, nil
		}
	}); err != nil {
		return "", nil
	}
	return ret.String, nil
}

func GenCopywriting(items []*taobaoopen.MaterialItem) (string, error) {
	if len(items) <= 0 {
		return "搜索结果为空", nil
	}
	ret := ""
	for _, v := range items {
		ret += v.ShortTitle + "\n"
		if v.CouponAmount == "" {
			tkl, err := tb.CreateTpwd(v.ShortTitle, "https:"+v.Url)
			if err != nil || tkl == "" {
				tkl = v.CouponShareUrl
			}
			ret += "价格:" + v.ZkFinalPrice + "￥ " + utils.ShortUrl("http://tb.icodef.com/tb.php?tkl="+url.QueryEscape(tkl)+"&pic="+url.QueryEscape(v.PictUrl)) + "\n"
		} else {
			coupon_start_fee, _ := strconv.ParseFloat(v.CouponStartFee, 64)
			zk_final_price, _ := strconv.ParseFloat(v.ZkFinalPrice, 64)
			if zk_final_price >= coupon_start_fee {
				coupon_amount, _ := strconv.ParseFloat(v.CouponAmount, 64)
				tkl, err := tb.CreateTpwd(v.ShortTitle, "https:"+v.CouponShareUrl)
				if err != nil || tkl == "" {
					tkl = v.CouponShareUrl
				}
				ret += "原价:" + v.ZkFinalPrice + "￥ 券后价:" + strconv.FormatFloat(zk_final_price-coupon_amount, 'G', 5, 64) + "￥ " + utils.ShortUrl("http://tb.icodef.com/tb.php?tkl="+url.QueryEscape(tkl)+"&pic="+url.QueryEscape(v.PictUrl)) + "\n"
			} else {
				tkl, err := tb.CreateTpwd(v.ShortTitle, "https:"+v.Url)
				if err != nil || tkl == "" {
					tkl = v.Url
				}
				ret += "价格:" + v.ZkFinalPrice + "￥ " + utils.ShortUrl("http://tb.icodef.com/tb.php?tkl="+url.QueryEscape(tkl)+"&pic="+url.QueryEscape(v.PictUrl)) + "\n"
			}
		}
	}
	return ret, nil
}
