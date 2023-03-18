package alimama

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/CodFrm/qqbot/config"
	"github.com/CodFrm/qqbot/db"
	"github.com/CodFrm/qqbot/utils"
	"github.com/CodFrm/qqbot/utils/iotqq"
	"github.com/CodFrm/qqbot/utils/jdunion"
	"github.com/CodFrm/qqbot/utils/taobaoopen"
)

var Tb *taobaoopen.Taobao
var jd *jdunion.JdUnion
var mq *broker
var tbfl *taobaoopen.Taobao
var topicList map[string]int

func Init() error {
	//c := cron.New(cron.WithSeconds())
	//c.AddFunc("0 30 7 * * ?", notice("美好的一天开始了,记得吃早餐哦."))
	//c.AddFunc("0 30 11 * * ?", notice("该吃中饭啦,要不要点外卖~"))
	//c.AddFunc("0 30 14 * * ?", notice("下午茶时间到啦,来份奶茶吧~"))
	//c.AddFunc("0 30 17 * * ?", notice("该吃晚饭啦,来一份外卖吧~"))
	//c.AddFunc("0 30 22 * * ?", notice("夜宵时间,来撸串"))
	//c.Start()
	topicList = make(map[string]int)
	Tb = taobaoopen.NewTaobao(config.AppConfig.Taobao)
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
					iotqq.QueueSendPrivateMsg(int(group), qq, "0-您关注的'"+keyword.topic+"'有新消息\n"+info+"\n回复'退订"+keyword.topic+"'可退订该消息.小程序/APP自助查券和返利https://m3w.cn/tyq")
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
			iotqq.QueueSendMsg(utils.StringToInt(v), 0, "0-"+t+"\n复制这条信息(建议收藏)，$5YiUccAeTlY$，到【手机淘宝】即可查看."+
				"美团可使用此链接:https://sourl.cn/Kvz8Hk\n小程序/APP自助查券和返利https://m3w.cn/tyq"+
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
	if tkl := utils.RegexMatch(msg, "[^\\w](\\w{8,12})[^\\w]"); len(tkl) >= 2 {
		if ret, err := tbfl.ConversionTkl(tkl[1]); err != nil {
			return msg, nil, err
		} else {
			if len(ret.Content) < 1 {
				return msg, nil, nil
			}
			newtkl := utils.RegexMatch(ret.Content[0].Tkl, "[^\\w](\\w{8,12})[^\\w]")
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
				newtkl := utils.RegexMatch(ret.Content[0].Tkl, "[^\\w](\\w{8,12})[^\\w]")
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
	if tkl := utils.RegexMatch(msg, "[^\\w](\\w{8,12})[^\\w]"); len(tkl) >= 2 {
		if ret, err := Tb.ConversionTkl(tkl[1]); err != nil {
			if err.Error() == "很抱歉！商品ID解析错误！！！" {
				// 获取口令链接和类型,判断是否为活动口令
				ret, err := Tb.ResolveTklAddress(tkl[1])
				if err != nil {
					return msg, nil, err
				}
				if ret.URLType == "10" || ret.URLType == "3" {
					//处理活动链接
					activeId := ret.URLID
					ret, err := Tb.GetActiveInfo(activeId)
					if err != nil {
						return msg, nil, err
					}
					if len(ret.Response.Data.PageName) <= 5 {
						ret.Response.Data.PageName = "省钱不吃土"
					}
					mytkl, err := Tb.CreateTpwd(ret.Response.Data.PageName, ret.Response.Data.ClickURL)
					if err != nil {
						return msg, nil, err
					}
					newtkl := utils.RegexMatch(mytkl, "[^\\w](\\w{8,12})[^\\w]")
					if len(newtkl) == 2 {
						msg = strings.ReplaceAll(msg, tkl[1], newtkl[1])
						re := regexp.MustCompile("(https?|ftp|file)://[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|]")
						msg = re.ReplaceAllString(msg, ret.Response.Data.ShortClickURL)
						return msg, &taobaoopen.ConverseTkl{
							Status: 0,
							Content: []taobaoopen.ConverseTklContent{
								{
									TaoID:    activeId,
									Tkl:      newtkl[1],
									Shorturl: ret.Response.Data.ShortClickURL,
								},
							},
						}, nil
					}
					return msg, nil, nil
				}
			} else {
				return msg, nil, err
			}
		} else {
			if len(ret.Content) < 1 {
				return msg, nil, nil
			}
			newtkl := utils.RegexMatch(ret.Content[0].Tkl, "[^\\w](\\w{8,12})[^\\w]")
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
		resp, err := utils.HttpGet(ret.Data.ShortURL, nil, nil)
		if err != nil {
			retTkl.Content = append(retTkl.Content, taobaoopen.ConverseTklContent{
				TaoID:    "error",
				Shorturl: ret.Data.ShortURL,
			})
		} else {
			tourl := utils.RegexMatch(string(resp), "hrl='(.*?)';")
			if len(tourl) > 0 {
				resp, err := utils.HttpGet(tourl[1], nil, nil)
				if err != nil {
					retTkl.Content = append(retTkl.Content, taobaoopen.ConverseTklContent{
						TaoID:    "error",
						Shorturl: ret.Data.ShortURL,
					})
				} else {
					ids := utils.RegexMatch(string(resp), "sku[iI]d:[\"\\s]{1,2}(\\d+)[\",]")
					if len(ids) > 0 {
						retTkl.Content = append(retTkl.Content, taobaoopen.ConverseTklContent{
							TaoID:    ids[1],
							Shorturl: ret.Data.ShortURL,
						})
					}
				}
			}
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
	return "搜索功能存在问题，暂未开放", nil
	list, err := Tb.MaterialSearch(keyword)
	if err != nil {
		return "", err
	}
	ret := &db.StringCache{}
	if err := db.GetOrSet("alimama:search:"+keyword+":v2", ret, func() (interface{}, error) {
		//ret := "网站上线啦,更强大更好用的搜索方式:" + "https://gw.icodef.com/pages/search/search?keyword=" + url.QueryEscape(keyword)
		//",直接访问此链接即可查看搜索结果:" + ShortUrl("https://gw.icodef.com/pages/search/search?keyword="+url.QueryEscape(keyword))
		//return &db.StringCache{String: ret}, nil
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
			tkl, err := Tb.CreateTpwd(v.ShortTitle, "https:"+v.Url)
			if err != nil || tkl == "" {
				tkl = v.CouponShareUrl
			}
			kl := utils.RegexMatch(tkl, "[^\\w](\\w{8,12})[^\\w]")
			if len(kl) == 2 {
				tkl = kl[1]
			}
			ret += "价格:" + v.ZkFinalPrice + "￥ " + ShortUrl("https://gw.icodef.com/tb.html?tkl="+url.QueryEscape(tkl)+"&pic="+url.QueryEscape(v.PictUrl)) + "\n"
		} else {
			coupon_start_fee, _ := strconv.ParseFloat(v.CouponStartFee, 64)
			zk_final_price, _ := strconv.ParseFloat(v.ZkFinalPrice, 64)
			if zk_final_price >= coupon_start_fee {
				coupon_amount, _ := strconv.ParseFloat(v.CouponAmount, 64)
				tkl, err := Tb.CreateTpwd(v.ShortTitle, "https:"+v.CouponShareUrl)
				if err != nil || tkl == "" {
					tkl = v.CouponShareUrl
				}
				kl := utils.RegexMatch(tkl, "[^\\w](\\w{8,12})[^\\w]")
				if len(kl) == 2 {
					tkl = kl[1]
				}
				ret += "原价:" + v.ZkFinalPrice + "￥ 券后价:" + strconv.FormatFloat(zk_final_price-coupon_amount, 'G', 5, 64) + "￥ " + ShortUrl("http://gw.icodef.com/tb.html?tkl="+url.QueryEscape(tkl)+"&pic="+url.QueryEscape(v.PictUrl)) + "\n"
			} else {
				tkl, err := Tb.CreateTpwd(v.ShortTitle, "https:"+v.Url)
				if err != nil || tkl == "" {
					tkl = v.Url
				}
				kl := utils.RegexMatch(tkl, "[^\\w](\\w{8,12})[^\\w]")
				if len(kl) == 2 {
					tkl = kl[1]
				}
				ret += "价格:" + v.ZkFinalPrice + "￥ " + ShortUrl("https://gw.icodef.com/tb.html?tkl="+url.QueryEscape(tkl)+"&pic="+url.QueryEscape(v.PictUrl)) + "\n"
			}
		}
	}
	return ret + "复制价格后方的口令到淘宝即可享受优惠", nil
}
