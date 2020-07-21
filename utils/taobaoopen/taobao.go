package taobaoopen

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

type Tabao struct {
	AppKey    string
	AppSecret string
	Entrance  string
}

type Kv struct {
	Key   string
	Value string
}

type Kvs []*Kv

func NewTaobao(AppKey, AppSecret, Entrance string) *Tabao {
	return &Tabao{
		AppKey:    AppKey,
		AppSecret: AppSecret,
		Entrance:  Entrance,
	}
}

func GenKv(key, val string) *Kv {
	return &Kv{
		Key:   key,
		Value: val,
	}
}

func (t *Tabao) PublicFunc(method string, param ...*Kv) (string, error) {
	if param == nil {
		param = make([]*Kv, 0)
	}
	param = append(param, GenKv("method", method))
	param = append(param, GenKv("app_key", t.AppKey))
	param = append(param, GenKv("timestamp", time.Now().Format("2006-01-02 15:04:05")))
	//param = append(param, GenKv("timestamp", "2020-07-21 10:36:28"))
	param = append(param, GenKv("format", "json"))
	param = append(param, GenKv("v", "2.0"))
	param = append(param, GenKv("sign_method", "hmac-sha256"))

	sign := t.Sign(param)
	param = append(param, GenKv("sign", sign))
	return t.request(param)
}

func (t *Tabao) request(param []*Kv) (string, error) {
	data := ""
	for _, v := range param {
		data += v.Key + "=" + url.QueryEscape(v.Value) + "&"
	}
	data = data[:len(data)-1]
	resp, err := http.Post(t.Entrance, "application/x-www-form-urlencoded;charset=utf-8", bytes.NewBuffer([]byte(data)))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (t *Tabao) Sign(param []*Kv) string {
	ret := ""
	sort.Sort(Kvs(param))
	for _, v := range param {
		ret += v.Key + v.Value
	}
	h := hmac.New(sha256.New, []byte(t.AppSecret))
	h.Write([]byte(ret))
	sha := hex.EncodeToString(h.Sum(nil))
	return strings.ToUpper(sha)
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
