package taobaoopen

import (
	"bytes"
	"net/http"
	"strings"
	"time"
)

type Kv struct {
	Key   string
	Value string
}

type TaobaoConfig struct {
	AppKey    string `yaml:"appKey"`
	AppSecret string `yaml:"appSecret"`
	Entrance  string
	AdzoneId  string `yaml:"adzoneId"`
	ZtkAppKey string `yaml:"ztkAppKey"`
	Sid       string `yaml:"sid"`
	Pid       string `yaml:"pid"`
}

type Kvs []*Kv

func GenKv(key, val string) *Kv {
	return &Kv{
		Key:   key,
		Value: val,
	}
}

func (k Kvs) Len() int {
	return len(k)
}

func (k Kvs) Less(i, j int) bool {
	return strings.Compare(k[i].Key, k[j].Key) < 0
}

func (k Kvs) Swap(i, j int) {
	k[i], k[j] = k[j], k[i]
}

func HttpGet(url string, header map[string]string) ([]byte, error) {
	c := http.Client{
		Transport: &http.Transport{},
		Timeout:   time.Second * 20,
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
