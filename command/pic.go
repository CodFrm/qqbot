package command

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/CodFrm/iotqq-plugins/config"
	"github.com/CodFrm/iotqq-plugins/db"
	"github.com/CodFrm/iotqq-plugins/model"
	"github.com/CodFrm/iotqq-plugins/utils"
	"github.com/CodFrm/iotqq-plugins/utils/iotqq"
	"github.com/nfnt/resize"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func RotatePic(command []string, pic *model.PicInfo) ([]image.Image, error) {
	if len(command) > 4 {
		return nil, errors.New("命令过多")
	}
	r := bytes.NewBuffer(pic.Byte)
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}
	if img.Bounds().Dx() > 2048 || img.Bounds().Dy() > 2048 {
		return nil, errors.New("图片过大(max:2048*2048)")
	}
	retImage := make([]image.Image, 0)
	var hd_deal = false
	for _, v := range command {
		command2 := strings.Split(v, "+")
		if len(command2) > 4 {
			continue
		}
		tmpimg := copyImg(img)
		for _, v := range command2 {
			switch v {
			case "镜像":
				tmpimg = mirror(tmpimg)
			case "垂直":
				tmpimg = rotate90(tmpimg)
			case "翻转":
				tmpimg = rotate180(tmpimg)
			case "放大":
				tmpimg = narrow(tmpimg, 1.1)
			case "缩小":
				tmpimg = narrow(tmpimg, 0.9)
			case "高清重制":
				if hd_deal {
					continue
				}
				var err error
				tmpimg, err = hd(tmpimg)
				if err != nil {
					fmt.Printf("%v\n", err)
					return nil, errors.New("高清重制失败")
				}
				hd_deal = true
			case "灰白", "RIP", "R.I.P.":
				tmpimg = gray(tmpimg)
			case "颜色反转":
				tmpimg = reverse(tmpimg)
			default:
				continue
			}
		}
		retImage = append(retImage, tmpimg)
	}
	return retImage, nil
}

func reverse(m image.Image) image.Image {
	reverse := image.NewRGBA(image.Rect(0, 0, m.Bounds().Dx(), m.Bounds().Dy()))
	for x := m.Bounds().Min.X; x < m.Bounds().Max.X; x++ {
		for y := m.Bounds().Min.Y; y < m.Bounds().Max.Y; y++ {
			colorRgb := m.At(x, y)
			r, g, b, a := colorRgb.RGBA()
			reverse.Set(x, y, color.RGBA{uint8(255 - (r >> 8)), uint8(255 - (g >> 8)), uint8(255 - (b >> 8)), uint8(a >> 8)})
		}
	}
	return reverse
}

func gray(m image.Image) image.Image {
	gray := image.NewRGBA(image.Rect(0, 0, m.Bounds().Dx(), m.Bounds().Dy()))
	for x := m.Bounds().Min.X; x < m.Bounds().Max.X; x++ {
		for y := m.Bounds().Min.Y; y < m.Bounds().Max.Y; y++ {
			colorRgb := m.At(x, y)
			_, g, _, a := colorRgb.RGBA()
			g_uint8 := uint8(g >> 8)
			a_uint8 := uint8(a >> 8)
			gray.Set(x, y, color.RGBA{g_uint8, g_uint8, g_uint8, a_uint8})
		}
	}
	return gray
}

type hdRespond struct {
	OutputUrl string `json:"output_url"`
}

func hd(m image.Image) (image.Image, error) {
	url := "https://api.deepai.org/api/waifu2x"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	buffer := bytes.NewBuffer(nil)
	if err := png.Encode(buffer, m); err != nil {
		return nil, err
	}
	part1, err := writer.CreateFormFile("image", filepath.Base("hd.jpg"))
	_, err = io.Copy(part1, buffer)
	if err != nil {
		return nil, err
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	client := &http.Client{
		Transport: &http.Transport{
			DialContext: proxy,
		},
		Timeout: time.Second * 60,
	}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("api-key", config.AppConfig.Hdkey)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	resp := &hdRespond{}
	if err := json.Unmarshal(body, resp); err != nil {
		return nil, err
	}
	imgResp, err := http.Get(resp.OutputUrl)
	if err != nil {
		return nil, err
	}
	defer imgResp.Body.Close()
	img, _, err := image.Decode(imgResp.Body)
	return img, err
}

func copyImg(m image.Image) image.Image {
	new := image.NewRGBA(image.Rect(0, 0, m.Bounds().Dx(), m.Bounds().Dy()))
	for x := m.Bounds().Min.X; x < m.Bounds().Max.X; x++ {
		for y := m.Bounds().Min.Y; y < m.Bounds().Max.Y; y++ {
			new.Set(x, y, m.At(x, y))
		}
	}
	return new
}

func rotate90(m image.Image) image.Image {
	rotate90 := image.NewRGBA(image.Rect(0, 0, m.Bounds().Dy(), m.Bounds().Dx()))
	for x := m.Bounds().Min.Y; x < m.Bounds().Max.Y; x++ {
		for y := m.Bounds().Max.X - 1; y >= m.Bounds().Min.X; y-- {
			rotate90.Set(m.Bounds().Max.Y-x-1, y, m.At(y, x))
		}
	}
	return rotate90
}

func mirror(m image.Image) image.Image {
	mirror := image.NewRGBA(image.Rect(0, 0, m.Bounds().Dx(), m.Bounds().Dy()))
	for x := m.Bounds().Min.X; x < m.Bounds().Max.X; x++ {
		for y := m.Bounds().Min.Y; y < m.Bounds().Max.Y; y++ {
			mirror.Set(x, y, m.At(m.Bounds().Max.X-x-1, y))
		}
	}
	return mirror
}

func rotate180(m image.Image) image.Image {
	rotate180 := image.NewRGBA(image.Rect(0, 0, m.Bounds().Dx(), m.Bounds().Dy()))
	for x := m.Bounds().Min.X; x < m.Bounds().Max.X; x++ {
		for y := m.Bounds().Min.Y; y < m.Bounds().Max.Y; y++ {
			rotate180.Set(m.Bounds().Max.X-x-1, m.Bounds().Max.Y-y-1, m.At(x, y))
		}
	}
	return rotate180
}

func narrow(m image.Image, scale float32) image.Image {
	return resize.Resize(uint(float32(m.Bounds().Dx())*scale), uint(float32(m.Bounds().Dy())*scale), m, resize.Lanczos3)
}

type moderate struct {
	Predictions struct {
		Teen     float64 `json:"teen"`
		Everyone float64 `json:"everyone"`
		Adult    float64 `json:"adult"`
	} `json:"predictions"`
}

func Gwkk(id string) ([]byte, error) {
	ok, err := db.Redis.SetNX("look:adult:"+id, "1", time.Hour*24*30).Result()
	if err != nil {
		return nil, errors.New("文件不存在")
	}
	if ok {
		return ioutil.ReadFile("./data/adult/" + id + ".png")
	}
	return nil, errors.New("已经有人鉴定了")
}

func IsAdult(args iotqq.Data, img *model.PicInfo) (int, bool, error) {
	//图片鉴黄
	url := "https://api.moderatecontent.com/moderate/?key=" + config.AppConfig.ModerateKey
	method := "POST"
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	part1,
		errFile1 := writer.CreateFormFile("file", filepath.Base("ht.jpg"))
	r := bytes.NewBuffer(img.Byte)
	_, errFile1 = io.Copy(part1, r)
	if errFile1 != nil {
		return 1, false, errFile1
	}
	err := writer.Close()
	if err != nil {
		return 1, false, err
	}
	client := &http.Client{
		Transport: &http.Transport{
			DialContext: proxy,
		},
		Timeout: time.Second * 60,
	}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		println(err.Error())
		return 1, false, nil
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		println(err.Error())
		return 1, false, nil
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	m := &moderate{}
	if err := json.Unmarshal(body, m); err != nil {
		return 1, false, err
	}
	if m.Predictions.Adult > 50 {
		id := utils.RandStringRunes(32)
		db.Redis.Set("adult:pic:"+id, body, 0)
		_ = ioutil.WriteFile("./data/adult/"+id+".png", img.Byte, 0775)
		if m.Predictions.Adult > 60 {
			user := strconv.FormatInt(args.FromUserID, 10)
			if user == config.AppConfig.QQ {
				reg := regexp.MustCompile("pixiv:(\\d+)")
				m := reg.FindStringSubmatch(args.Content)
				if len(m) > 0 {
					if u := db.Redis.Get("pixiv:send:qq:" + m[1]).Val(); u != "" {
						user = u
					}
				}
			}
			key := "adult:ban:" + user
			if len, _ := db.Redis.LPush(key, time.Now().Unix()).Result(); len == 1 {
				db.Redis.Expire(key, time.Hour)
			}
			flag := false
			if t, _ := db.Redis.LIndex(key, 4).Int64(); time.Now().Unix()-t < 600 {
				BanUser(args, user)
				flag = true
			}
			if m.Predictions.Adult > 90 {
				return 3, flag, errors.New("你的图片带有不宜内容,请注意你的言辞,图片已撤回,证据已保留ID:" + id)
			}
			return 2, flag, errors.New("你的图片可能带有不宜内容,请注意你的言辞,证据已保留ID:" + id)
		}
		return 4, false, errors.New("有点涩涩,保存了ID:" + id)
	}
	return 1, false, nil
}

func BanUser(args iotqq.Data, user string) {
	u, _ := strconv.ParseInt(user, 10, 64)
	iotqq.SendMsg(args.FromGroupID, u, "加入黑名单")
	time.Sleep(1 * time.Second)
	iotqq.ShutUp(args.FromGroupID, u, 86400)
	BlackList(user, "1", "86400")
}
