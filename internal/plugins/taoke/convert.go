package taoke

import (
	"regexp"
	"strings"

	"github.com/CodFrm/qqbot/utils"
	"github.com/CodFrm/qqbot/utils/taobaoopen"
)

func (t *TaoKe) DealTkl(msg string) (string, *taobaoopen.ConverseTkl, error) {
	if tkl := utils.RegexMatch(msg, "[^\\w](\\w{8,12})[^\\w]"); len(tkl) >= 2 {
		if ret, err := t.tb.ConversionTkl(tkl[1]); err != nil {
			if err.Error() == "很抱歉！商品ID解析错误！！！" {
				// 获取口令链接和类型,判断是否为活动口令
				ret, err := t.tb.ResolveTklAddress(tkl[1])
				if err != nil {
					return msg, nil, err
				}
				if ret.URLType == "10" || ret.URLType == "3" {
					//处理活动链接
					activeId := ret.URLID
					ret, err := t.tb.GetActiveInfo(activeId)
					if err != nil {
						return msg, nil, err
					}
					if len(ret.Response.Data.PageName) <= 5 {
						ret.Response.Data.PageName = "省钱不吃土"
					}
					mytkl, err := t.tb.CreateTpwd(ret.Response.Data.PageName, ret.Response.Data.ClickURL)
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
		ret, err := t.jd.ConversionLink(v[0])
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
