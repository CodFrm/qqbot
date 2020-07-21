package alimama

import (
	"github.com/CodFrm/iotqq-plugins/config"
	"github.com/CodFrm/iotqq-plugins/utils/iotqq"
	"github.com/CodFrm/iotqq-plugins/utils/taobaoopen"
	"github.com/robfig/cron/v3"
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
		iotqq.QueueSendMsg(v.GroupId, 0, "快到饭点了,来一份外卖吧~\nhttps://m.tb.cn/h.VsVRwwj\n复制这条信息，$nH3n1zNqDip$，到【手机淘宝】即可查看."+
			"美团可使用此链接:https://sourl.cn/Kvz8Hk\n"+
			"后续将会通过QQ红包提供返现功能")
	}
}

func Search(keyword string) ([]*taobaoopen.MaterialItem, error) {
	list, err := tb.MaterialSearch(keyword)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func GenCopywriting(items []*taobaoopen.MaterialItem) (string, error) {
	if len(items) <= 0 {
		return "搜索结果为空", nil
	}
	ret := ""
	urls := make([]string, 0)
	couponUrls := make([]string, 0)
	for _, v := range items {
		if v.Url != "" {
			urls = append(urls, "https:"+v.Url)
		}
		if v.CouponShareUrl != "" {
			couponUrls = append(couponUrls, "https:"+v.CouponShareUrl)
		}
	}
	var err error
	var urlsSpread []*taobaoopen.SpreadItem
	var couponUrlsSpread []*taobaoopen.SpreadItem
	if len(urls) > 0 {
		urlsSpread, err = tb.GetSpread(urls)
		if err != nil {
			return "", err
		}
	}
	if len(couponUrls) > 0 {
		couponUrlsSpread, err = tb.GetSpread(couponUrls)
		if err != nil {
			return "", err
		}
	}
	urlsIndex := 0
	couponUrlsIndex := 0
	for _, v := range items {
		ret += v.ShortTitle + "\n"
		if v.CouponAmount == "" {
			ret += "价格:" + v.ZkFinalPrice + "￥ " + urlsSpread[urlsIndex].Content + "\n"
		} else {
			coupon_start_fee, _ := strconv.ParseFloat(v.CouponStartFee, 64)
			zk_final_price, _ := strconv.ParseFloat(v.ZkFinalPrice, 64)
			if zk_final_price >= coupon_start_fee {
				coupon_amount, _ := strconv.ParseFloat(v.CouponAmount, 64)
				ret += "原价:" + v.ZkFinalPrice + "￥ 券后价:" + strconv.FormatFloat(zk_final_price-coupon_amount, 'G', 5, 64) + "￥ " + couponUrlsSpread[couponUrlsIndex].Content + "\n"
			} else {
				ret += "价格:" + v.ZkFinalPrice + "￥ " + urlsSpread[urlsIndex].Content + "\n"
			}
		}
		if v.Url != "" {
			urlsIndex++
		}
		if v.CouponShareUrl != "" {
			couponUrlsIndex++
		}
	}
	return ret + "返利功能后续开发", nil
}
