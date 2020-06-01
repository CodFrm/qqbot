package command

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/CodFrm/iotqq-plugins/config"
	"github.com/CodFrm/iotqq-plugins/model"
	"github.com/nfnt/resize"
	"image"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

func RotatePic(command []string, pic *model.PicInfo) ([]image.Image, error) {
	if len(command) > 4 {
		return nil, errors.New("命令过多")
	}
	resp, err := http.Get(pic.Url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, err
	}
	if img.Bounds().Dx() > 1024 || img.Bounds().Dy() > 1024 {
		return nil, errors.New("图片过大(max:1024*1024)")
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
			default:
				continue
			}
		}
		if tmpimg.Bounds().Dx() > 1024 || tmpimg.Bounds().Dy() > 1024 {
			return nil, errors.New("图片过大(max:1024*1024)")
		}
		retImage = append(retImage, tmpimg)
	}
	return retImage, nil
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
		Timeout: time.Second * 10,
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
			rotate90.Set(m.Bounds().Max.Y-x, y, m.At(y, x))
		}
	}
	return rotate90
}

func mirror(m image.Image) image.Image {
	mirror := image.NewRGBA(image.Rect(0, 0, m.Bounds().Dx(), m.Bounds().Dy()))
	for x := m.Bounds().Min.X; x < m.Bounds().Max.X; x++ {
		for y := m.Bounds().Min.Y; y < m.Bounds().Max.Y; y++ {
			mirror.Set(x, y, m.At(m.Bounds().Max.X-x, y))
		}
	}
	return mirror
}

func rotate180(m image.Image) image.Image {
	rotate180 := image.NewRGBA(image.Rect(0, 0, m.Bounds().Dx(), m.Bounds().Dy()))
	for x := m.Bounds().Min.X; x < m.Bounds().Max.X; x++ {
		for y := m.Bounds().Min.Y; y < m.Bounds().Max.Y; y++ {
			rotate180.Set(m.Bounds().Max.X-x, m.Bounds().Max.Y-y, m.At(x, y))
		}
	}
	return rotate180
}

func narrow(m image.Image, scale float32) image.Image {
	return resize.Resize(uint(float32(m.Bounds().Dx())*scale), uint(float32(m.Bounds().Dy())*scale), m, resize.Lanczos3)
}
