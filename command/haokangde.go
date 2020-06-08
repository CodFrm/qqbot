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

var PicIsNil = errors.New("我真的一张都没有了")

func HaoKangDe(command string) ([]byte, *model.PixivPicItem, error) {
	page, _ := db.Redis.Get("pixivlist" + command + ":page").Int()
	if page <= 0 {
		page = 1
	}
	img, err := getPixivPicByCommand(command, page)
	if err != nil {
		return nil, nil, err
	}
	imgbyte, err := downloadPixivPic(img)
	return imgbyte, img, err
}

func GetPixivImg(id string) ([]byte, error) {
	_, err := os.Stat("./data/pixiv/" + id + "_big.jpg")
	if err == nil {
		return nil, errors.New("已经发送过了啦")
	}
	b, err := ioutil.ReadFile("./data/pixiv/" + id + ".json")
	if err != nil {
		return nil, errors.New("图片缓存不存在,请给我看过的图片")
	}
	m := &model.PixivIllust{}
	if err := json.Unmarshal(b, m); err != nil {
		return nil, errors.New("错误的缓存")
	}
	data, err := utils.HttpGet(m.Body.Urls.Original, map[string]string{
		"Cookie":  config.AppConfig.Pixiv.Cookie,
		"Referer": "https://www.pixiv.net/artworks/" + id,
	}, proxy)
	if err != nil {
		return nil, errors.New("网络错误,请稍后重试")
	}
	if err := ioutil.WriteFile("./data/pixiv/"+id+"_big.jpg", data, 0755); err != nil {
		return nil, err
	}
	return data, nil
}

func getPixivPicByCommand(command string, page int) (*model.PixivPicItem, error) {
	var data []*model.PixivPicItem
	var err error
	if command == "" {
		data, err = pixivRankList(page)
	} else {
		data, err = pixivList(command, page)
	}
	if err != nil {
		return nil, err
	}
	img, err := uniqueRand(command, data)
	if err != nil {
		if err == PicIsNil {
			page = page + 1
			page = SetPage(command, page)
			return getPixivPicByCommand(command, page)
		}
		return nil, err
	}
	return img, err
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

//30天内不再重复
func uniqueRand(tag string, data []*model.PixivPicItem) (*model.PixivPicItem, error) {
	randList := make([]*model.PixivPicItem, 0)
	for _, v := range data {
		if !db.Redis.HExists("uniqueRand"+tag, v.Id).Val() {
			randList = append(randList, v)
		}
	}
	if len(randList) == 0 {
		return nil, PicIsNil
	}
	ret := randList[rand.Intn(len(randList))]
	db.Redis.HSet("uniqueRand"+tag, ret.Id, "1")
	if db.Redis.HLen("uniqueRand"+tag).Val() <= 1 {
		db.Redis.Expire("uniqueRand"+tag, time.Second*86400*30)
	}
	return ret, nil
}

var Hots = []string{"30000", "20000", "10000", "5000"}

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
	if err := db.Set("pixivTag"+tag, ret, db.WithTTL(time.Second*86400*7)); err != nil {
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

func pixivRankList(page int) ([]*model.PixivPicItem, error) {
	picList := make([]*model.PixivPicItem, 0)
	if err := db.GetOrSet("pixivList:rank:"+strconv.Itoa(page), &picList, func() (interface{}, error) {
		str, err := utils.HttpGet("https://www.pixiv.net/ranking.php?mode=weekly&content=illust&p="+strconv.Itoa(page)+"&format=json",
			map[string]string{
				"Cookie":  config.AppConfig.Pixiv.Cookie,
				"Referer": "https://www.pixiv.net/ranking.php?mode=weekly&content=illust",
			}, proxy)
		if err != nil {
			return "", err
		}
		m := &model.PixivRankList{}
		if err := json.Unmarshal(str, m); err != nil {
			return nil, err
		}
		ret := make([]*model.PixivPicItem, 0)
		for _, v := range m.Contents {
			ret = append(ret, &model.PixivPicItem{
				Id:              strconv.Itoa(v.IllustId),
				ProfileImageUrl: v.ProfileImg,
				Url:             v.Url,
				UserId:          strconv.Itoa(v.UserId),
				UserName:        v.UserName,
				Title:           v.Title,
			})
		}
		return ret, nil
	}, db.WithTTL(time.Second*86400)); err != nil {
		return picList, err
	}
	return picList, nil
}

func SetPage(tag string, page int) int {
	end, _ := db.Redis.Get("pixivlist" + tag + ":page:expire").Int64()
	if end > time.Now().Unix() && end > 0 {
	} else {
		page = 1
		db.Redis.Set("pixivlist"+tag+":page:expire", time.Now().Unix()+86400*3, time.Second*86400*3)
	}
	db.Redis.Set("pixivlist"+tag+":page", page, time.Second*86400*3)
	return page
}

func getPicList(tag string, hot int, page int) func() (i interface{}, err error) {
	return func() (i interface{}, err error) {
		var relateTag string
		relateTag, hot = getTagCache(tag)
		if relateTag == "" {
			relateTag = tag
		}
		str, err := utils.HttpGet("https://www.pixiv.net/ajax/search/illustrations/"+tagurlencode(relateTag, hot)+
			"?word="+tagurlencode(relateTag, hot)+"&order=date_d&mode=safe&p="+strconv.Itoa(page)+"&s_mode=s_tag&type=illust_and_ugoira&lang=zh",
			map[string]string{
				"Cookie":  config.AppConfig.Pixiv.Cookie,
				"Referer": "https://www.pixiv.net/tags/" + tagurlencode(relateTag, hot) + "/illustrations?s_mode=s_tag",
			}, proxy)
		if err != nil {
			return "", err
		}
		m := &model.IllustRespond{}
		if err := json.Unmarshal(str, m); err != nil {
			return "", err
		}
		//图片过少
		if len(m.Body.Illust.Data) <= 10 {
			if hot >= 3 {
				if tag != relateTag {
					return nil, errors.New("图片过少")
				}
				relateTags, _ := getRelateTags(tag)
				if len(relateTags) == 0 {
					return nil, errors.New("图片过少")
				} else {
					relateTag = relateTags[0]
				}
				setTagCache(tag, relateTag, 0)
				SetPage(tag, 1)
				return getPicList(tag, 0, 1)()
			}
			time.Sleep(time.Second * 1)
			setTagCache(tag, relateTag, hot+1)
			SetPage(tag, 1)
			return getPicList(tag, hot+1, 1)()
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
