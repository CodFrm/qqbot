package utils

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func ImageToBase64(img image.Image) (string, error) {
	buffer := bytes.NewBuffer(nil)
	if err := jpeg.Encode(buffer, img, &jpeg.Options{Quality: 100}); err != nil {
		return "", err
	}
	ret := base64.StdEncoding.EncodeToString(buffer.Bytes())
	return ret, nil
}

func FileBase64(path string) string {
	f, _ := ioutil.ReadFile(path)
	return base64.StdEncoding.EncodeToString(f)
}

func SaveImage(path string, img image.Image) (err error) {
	// 需要保存的文件
	imgfile, err := os.Create(path)
	defer imgfile.Close()
	err = png.Encode(imgfile, img)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func ImageCompression(img image.Image) (image.Image, error) {
	buffer := bytes.NewBuffer(nil)
	if err := jpeg.Encode(buffer, img, &jpeg.Options{Quality: 100}); err != nil {
		return nil, err
	}
	return jpeg.Decode(buffer)
}

func HttpGet(url string, header map[string]string, proxy func(ctx context.Context, network, addr string) (net.Conn, error)) ([]byte, error) {
	c := http.Client{
		Transport: &http.Transport{
			DialContext: proxy,
		},
		Timeout: time.Second * 60,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []byte{}, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.61 Safari/537.36")
	for k, v := range header {
		req.Header.Set(k, v)
	}
	resp, err := c.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	return buf.Bytes(), nil
}

var cfCache = make(map[string]string)

func CloudflareResolve(hostname string) (string, error) {
	if v, ok := cfCache[hostname]; ok {
		return v, nil
	}
	resp, err := http.Get("https://cloudflare-dns.com/dns-query?name=" + hostname + "&ct=application/dns-json&type=A&do=false&cd=false")
	if err != nil {
		return "", nil
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	m := make(map[string]interface{})
	if err := json.Unmarshal(body, &m); err != nil {
		return "", err
	}
	tmp, ok := m["Answer"].([]interface{})
	if !ok {
		return "", errors.New("error asnwer")
	}
	if len(tmp) == 0 {
		return "", errors.New("0 answer")
	}
	cfCache[hostname] = tmp[0].(map[string]interface{})["data"].(string)
	return cfCache[hostname], nil
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func StringToInt(i string) int {
	ret, _ := strconv.Atoi(i)
	return ret
}

func StringToInt64(i string) int64 {
	ret, _ := strconv.ParseInt(i, 10, 64)
	return ret
}

func Nmsl() string {
	f := rand.Intn(2)
	var err error
	var resp []byte
	if f == 1 {
		resp, err = HttpGet("https://nmsl.shadiao.app/api.php?level=min&lang=zh_cn", nil, nil)
	} else {
		resp, err = HttpGet("https://nmsl.shadiao.app/api.php?lang=zh_cn", nil, nil)
	}
	if err != nil {
		return "骂不出来了,劳资举报一手"
	}
	ret := string(resp)
	ret = strings.ReplaceAll(ret, "妈", "🐴")
	return ret
}
