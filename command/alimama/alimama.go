package alimama

import (
	"fmt"
	"github.com/CodFrm/iotqq-plugins/config"
	"github.com/CodFrm/iotqq-plugins/db"
	"github.com/CodFrm/iotqq-plugins/utils"
	"github.com/CodFrm/iotqq-plugins/utils/iotqq"
	"github.com/CodFrm/iotqq-plugins/utils/jdunion"
	"github.com/CodFrm/iotqq-plugins/utils/taobaoopen"
	"github.com/robfig/cron/v3"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

var tb *taobaoopen.Taobao
var jd *jdunion.JdUnion
var mq *broker
var tbfl *taobaoopen.Taobao
var topicList map[string]int

func Init() error {
	c := cron.New(cron.WithSeconds())
	c.AddFunc("0 30 7 * * ?", notice("美好的一天开始了,记得吃早餐哦."))
	c.AddFunc("0 30 11 * * ?", notice("该吃中饭啦,要不要点外卖~"))
	c.AddFunc("0 30 14 * * ?", notice("下午茶时间到啦,来份奶茶吧~"))
	c.AddFunc("0 30 17 * * ?", notice("该吃晚饭啦,来一份外卖吧~"))
	c.AddFunc("0 30 22 * * ?", notice("夜宵时间,来撸串"))
	c.Start()
	tb = taobaoopen.NewTaobao(config.AppConfig.Taobao)
	tbfl = taobaoopen.NewTaobao(config.AppConfig.TaobaoFl)
	jd = jdunion.NewJdUnion(config.AppConfig.Jd)
	mq = NewBroker()
	//初始化mq订阅
	mlist, err := db.Redis.HGetAll("alimama:subscribe:list").Result()
	if err != nil {
		return err
	}
	for qq, group := range mlist {
		tlist, err := db.Redis.SMembers("alimama:subscribe:topic:" + qq).Result()
		if err != nil {
			return err
		}
		for _, topic := range tlist {
			topicList[topic]++
			mq.subscribe(topic, qq, &subscribe{
				handler: func(info string, keyword *publisher) {
					group, _ := strconv.ParseInt(keyword.param, 10, 64)
					qq, _ := strconv.ParseInt(keyword.tag, 10, 64)
					iotqq.QueueSendPrivateMsg(int(group), qq, "您关注的'"+keyword.topic+"'有新消息\n"+info+"\n回复'退订"+keyword.topic+"'可退订该消息.不加关键字退订全部")
				},
				param: group,
			})
		}
	}
	return nil
}

func notice(t string) func() {
	return func() {
		list, err := db.Redis.SMembers("alimama:group:list").Result()
		if err != nil {
			return
		}
		for _, v := range list {
			iotqq.QueueSendMsg(utils.StringToInt(v), 0, t+"\nhttps://sourl.cn/FhPLTD\n复制这条信息，$nH3n1zNqDip$，到【手机淘宝】即可查看."+
				"美团可使用此链接:https://sourl.cn/Kvz8Hk"+
				"")
		}
	}
}

func AddGroup(qqgroup string, rm bool) error {
	if rm {
		return db.Redis.SRem("alimama:group:list", qqgroup).Err()
	}
	return db.Redis.SAdd("alimama:group:list", qqgroup).Err()
}

func DealTklFl(msg string) (string, *taobaoopen.ConverseTkl, error) {
	if tkl := utils.RegexMatch(msg, "[\\p{Sc}](\\w{8,12})[\\p{Sc}]"); len(tkl) >= 2 {
		if ret, err := tbfl.ConversionTkl(tkl[1]); err != nil {
			return msg, nil, err
		} else {
			if len(ret.Content) < 1 {
				return msg, nil, nil
			}
			newtkl := utils.RegexMatch(ret.Content[0].Tkl, "[\\p{Sc}](\\w{8,12})[\\p{Sc}]")
			if len(newtkl) == 2 {
				msg = strings.ReplaceAll(msg, tkl[1], newtkl[1])
				re := regexp.MustCompile("(https?|ftp|file)://[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|]")
				msg = re.ReplaceAllString(msg, ret.Content[0].Shorturl)
				return msg, ret, nil
			}
			return msg, nil, nil
		}
	}
	//处理淘宝链接
	links := utils.RegexMatchs(msg, "http[s]://[\\w-]+\\.((taobao|tmall).com)[_:\\.\\/\\w\\?=&%]+")
	retTkl := &taobaoopen.ConverseTkl{
		Content: make([]taobaoopen.ConverseTklContent, 0),
	}
	if len(links) > 0 {
		for _, v := range links {
			id := utils.RegexMatch(v[0], "(\\?|&)id=(.*?)(&|$)")
			if len(id) == 0 {
				continue
			}
			if ret, err := tbfl.ConversionShopId(id[2]); err != nil {
				continue
			} else {
				if len(ret.Content) < 1 {
					continue
				}
				newtkl := utils.RegexMatch(ret.Content[0].Tkl, "[\\p{Sc}](\\w{8,12})[\\p{Sc}]")
				if len(newtkl) == 2 {
					msg = strings.ReplaceAll(msg, v[0], ret.Content[0].Shorturl)
					retTkl = ret
				}
			}
		}
		if len(retTkl.Content) == 0 {
			retTkl.Content = append(retTkl.Content, taobaoopen.ConverseTklContent{})
		}
		return msg, retTkl, nil
	}
	//处理京东链接
	for _, v := range utils.RegexMatchs(msg, "http[s]://[\\w-]+\\.(jd.com)[\\.\\/\\w\\?=&%]+") {
		ret, err := jd.ConversionLink(v[0])
		if err != nil {
			continue
		}
		msg = strings.ReplaceAll(msg, v[0], ret.Data.ShortURL)
		ids := utils.RegexMatch(ret.Data.ClickURL, "e=(.*?)&")
		if len(ids) > 0 {
			retTkl.Content = append(retTkl.Content, taobaoopen.ConverseTklContent{
				TaoID:    ids[1],
				Shorturl: ret.Data.ShortURL,
			})
		}
	}
	if len(retTkl.Content) == 0 {
		retTkl.Content = append(retTkl.Content, taobaoopen.ConverseTklContent{})
	}
	return msg, retTkl, nil
}

func DealTkl(msg string) (string, *taobaoopen.ConverseTkl, error) {
	if tkl := utils.RegexMatch(msg, "[\\p{Sc}](\\w{8,12})[\\p{Sc}]"); len(tkl) >= 2 {
		if ret, err := tb.ConversionTkl(tkl[1]); err != nil {
			return msg, nil, err
		} else {
			if len(ret.Content) < 1 {
				return msg, nil, nil
			}
			newtkl := utils.RegexMatch(ret.Content[0].Tkl, "[\\p{Sc}](\\w{8,12})[\\p{Sc}]")
			if len(newtkl) == 2 {
				msg = strings.ReplaceAll(msg, tkl[1], newtkl[1])
				re := regexp.MustCompile("(https?|ftp|file)://[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|]")
				msg = re.ReplaceAllString(msg, ret.Content[0].Shorturl)
				return msg, ret, nil
			}
			return msg, nil, nil
		}
	}
	//处理京东链接
	retTkl := &taobaoopen.ConverseTkl{
		Content: make([]taobaoopen.ConverseTklContent, 0),
	}
	for _, v := range utils.RegexMatchs(msg, "http[s]://[\\w-]+\\.(jd.com)[\\/\\w\\?=&%]+") {
		ret, err := jd.ConversionLink(v[0])
		if err != nil {
			continue
		}
		msg = strings.ReplaceAll(msg, v[0], ret.Data.ShortURL)
		ids := utils.RegexMatch(ret.Data.ClickURL, "e=(.*?)&")
		if len(ids) > 0 {
			retTkl.Content = append(retTkl.Content, taobaoopen.ConverseTklContent{
				TaoID:    ids[1],
				Shorturl: ret.Data.ShortURL,
			})
		}
	}
	if len(retTkl.Content) > 0 && retTkl.Content[0].TaoID != "" {
		return msg, retTkl, nil
	}
	return msg, nil, nil
}

func DealFl(fl string) string {
	ret, err := strconv.ParseFloat(fl, 64)
	if err != nil {
		return "0"
	}
	ret = ret * 0.55
	return fmt.Sprintf("%.2f", ret)
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
			ret += "价格:" + v.ZkFinalPrice + "￥ " + ShortUrl("http://tb.icodef.com/tb.php?tkl="+url.QueryEscape(tkl)+"&pic="+url.QueryEscape(v.PictUrl)) + "\n"
			//ret += "价格:" + v.ZkFinalPrice + "￥ " + tkl + "\n"
		} else {
			coupon_start_fee, _ := strconv.ParseFloat(v.CouponStartFee, 64)
			zk_final_price, _ := strconv.ParseFloat(v.ZkFinalPrice, 64)
			if zk_final_price >= coupon_start_fee {
				coupon_amount, _ := strconv.ParseFloat(v.CouponAmount, 64)
				tkl, err := tb.CreateTpwd(v.ShortTitle, "https:"+v.CouponShareUrl)
				if err != nil || tkl == "" {
					tkl = v.CouponShareUrl
				}
				ret += "原价:" + v.ZkFinalPrice + "￥ 券后价:" + strconv.FormatFloat(zk_final_price-coupon_amount, 'G', 5, 64) + "￥ " + ShortUrl("http://tb.icodef.com/tb.php?tkl="+url.QueryEscape(tkl)+"&pic="+url.QueryEscape(v.PictUrl)) + "\n"
				//ret += "原价:" + v.ZkFinalPrice + "￥ 券后价:" + strconv.FormatFloat(zk_final_price-coupon_amount, 'G', 5, 64) + "￥ " +
				//	tkl + "\n"
			} else {
				tkl, err := tb.CreateTpwd(v.ShortTitle, "https:"+v.Url)
				if err != nil || tkl == "" {
					tkl = v.Url
				}
				ret += "价格:" + v.ZkFinalPrice + "￥ " + ShortUrl("http://tb.icodef.com/tb.php?tkl="+url.QueryEscape(tkl)+"&pic="+url.QueryEscape(v.PictUrl)) + "\n"
				//ret += "价格:" + v.ZkFinalPrice + "￥ " + tkl + "\n"
			}
		}
	}
	return ret + "复制价格后方的口令到淘宝即可享受优惠", nil
}
