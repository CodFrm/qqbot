package command

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/CodFrm/iotqq-plugins/config"
	"github.com/CodFrm/iotqq-plugins/db"
	"github.com/CodFrm/iotqq-plugins/model"
	"github.com/CodFrm/iotqq-plugins/utils"
	"io"
	"io/ioutil"
	"math/rand"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

func HaoKangDe(command string) ([]byte, error) {
	data, err := pixivList(command, 1)
	if err != nil {
		return nil, err
	}
	//随机选取一张,保存到本地
	img, err := uniqueRand(command, data)
	return downloadPixivPic(img)
}

func downloadPixivPic(pic *model.PixivPicItem) ([]byte, error) {
	_, err := os.Stat("./data/pixiv/" + pic.Id + ".jpg")
	var r io.Reader
	if err != nil {
		if !os.IsExist(err) {
			os.MkdirAll("./data/pixiv", os.ModePerm)
			//download file
			data, err := utils.HttpGet("https://www.pixiv.net/ajax/illust/"+pic.Id+"?lang=zh", map[string]string{
				"Cookie": config.AppConfig.Pixiv.Cookie,
			}, proxy)
			if err != nil {
				return nil, err
			}

			m := &model.PixivIllust{}
			if err := json.Unmarshal(data, m); err != nil {
				return nil, err
			}

			if err := ioutil.WriteFile("./data/pixiv/"+pic.Id+".json", data, 0755); err != nil {
				return nil, err
			}

			data, err = utils.HttpGet(m.Body.Urls.Small, map[string]string{
				"Cookie":  config.AppConfig.Pixiv.Cookie,
				"Referer": "https://www.pixiv.net/artworks/" + pic.Id,
			}, proxy)
			if err := ioutil.WriteFile("./data/pixiv/"+pic.Id+".jpg", data, 0755); err != nil {
				return nil, err
			}
			r = bytes.NewReader(data)
		}
	} else {
		file, err := os.Open("./data/pixiv/" + pic.Id + ".jpg")
		if err != nil {
			return nil, err
		}
		defer file.Close()
		r = file
	}
	return ioutil.ReadAll(r)
}

//一天内不再重复
func uniqueRand(tag string, data []*model.PixivPicItem) (*model.PixivPicItem, error) {
	randList := make([]*model.PixivPicItem, 0)
	for _, v := range data {
		if !db.Redis.HExists("uniqueRand"+tag, v.Id).Val() {
			randList = append(randList, v)
		}
	}
	if len(randList) == 0 {
		return nil, errors.New("我真的一张都没有了")
	}
	ret := randList[rand.Intn(len(randList))]
	db.Redis.HSet("uniqueRand"+tag, ret.Id, "1")
	if db.Redis.HLen("uniqueRand"+tag).Val() <= 1 {
		db.Redis.Expire("uniqueRand"+tag, time.Second*86400*30)
	}
	return ret, nil
}

var Hots = []string{"30000", "20000", "10000"}

func tagurlencode(tag string, hot int) string {
	return strings.ReplaceAll(url.QueryEscape(tag+" "+Hots[hot]+"users入り"), "+", "%20")
}

func getTagCache(tag string) (string, int) {
	ret := &struct {
		db.StringCache
		db.IntCache
	}{}
	if err := db.Get("pixivTag"+tag, ret); err != nil {
		return "", 0
	}
	return ret.String, ret.Int
}

func setTagCache(tag, RelatedTag string, hot int) error {
	ret := &struct {
		db.StringCache
		db.IntCache
	}{}
	ret.String = RelatedTag
	ret.Int = hot
	if err := db.Set("pixivTag"+tag, ret, db.WithTTL(time.Second*86400*30)); err != nil {
		return err
	}
	return nil
}

func pixivList(tag string, page int) ([]*model.PixivPicItem, error) {
	picList := make([]*model.PixivPicItem, 0)
	if err := db.GetOrSet("pixivList"+tag+":"+strconv.Itoa(page), &picList, getPicList(tag, 0, page), db.WithTTL(time.Second*86400)); err != nil {
		return picList, err
	}
	return picList, nil
}

func getPicList(tag string, hot int, page int) func() (i interface{}, err error) {
	return func() (i interface{}, err error) {
		var relateTag string
		relateTag, hot = getTagCache(tag)
		if relateTag == "" {
			relateTag = tag
		}
		str, err := utils.HttpGet("https://www.pixiv.net/ajax/search/illustrations/"+tagurlencode(relateTag, hot)+
			"?word="+tagurlencode(relateTag, hot)+"&order=date_d&mode=safe&p=1&s_mode=s_tag&type=illust_and_ugoira&lang=zh",
			map[string]string{
				"Cookie": config.AppConfig.Pixiv.Cookie,
			}, proxy)
		if err != nil {
			return "", err
		}
		m := &model.IllustRespond{}
		if err := json.Unmarshal(str, m); err != nil {
			return "", err
		}
		//图片过少
		if m.Body.Illust.Total <= 10 {
			if hot >= 2 {
				if tag != relateTag {
					return nil, errors.New("图片过少")
				}
				relateTags, _ := getRelateTags(tag)
				if len(relateTags) == 0 {
					relateTag = tag
				} else {
					relateTag = relateTags[0]
				}
				setTagCache(tag, relateTag, hot)
				return getPicList(tag, hot+1, page)()
			}
			time.Sleep(time.Second * 4)
			setTagCache(tag, relateTag, hot+1)
			return getPicList(tag, hot+1, page)()
		}
		return m.Body.Illust.Data, nil
	}
}

func getRelateTags(tag string) ([]string, error) {
	str, err := utils.HttpGet("https://www.pixiv.net/rpc/cps.php?keyword="+url.QueryEscape(tag)+"&lang=zh",
		map[string]string{
			"Cookie":  config.AppConfig.Pixiv.Cookie,
			"Referer": "https://www.pixiv.net/",
		}, proxy)
	if err != nil {
		return nil, err
	}
	m := &model.PixivTags{}
	if err := json.Unmarshal(str, m); err != nil {
		return nil, err
	}
	ret := make([]string, 0)
	for _, v := range m.Candidates {
		ret = append(ret, v.TagName)
	}
	return ret, nil
}
