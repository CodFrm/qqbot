package jdunion

import (
	"crypto/hmac"
	"crypto/md5"
	"encoding/hex"
	"sort"
	"strings"
	"time"
)

type JdUnion struct {
	AppKey    string
	AppSecret string
}

func NewJdUnion() *JdUnion {
	return &JdUnion{}
}

func (jd *JdUnion) GetPromotionLink() {

}

func (jd *JdUnion) PublicFunc(method string, param ...*Kv) (string, error) {
	if param == nil {
		param = make([]*Kv, 0)
	}
	param = append(param, GenKv("method", method))
	param = append(param, GenKv("app_key", jd.AppKey))
	param = append(param, GenKv("timestamp", time.Now().Format("2006-01-02 15:04:05")))
	param = append(param, GenKv("format", "json"))
	param = append(param, GenKv("v", "1.0"))
	param = append(param, GenKv("sign_method", "md5"))

	sign := jd.Sign(param)
	param = append(param, GenKv("sign", sign))
	return jd.request(param)
}

func (jd *JdUnion) Sign(param []*Kv) string {
	ret := ""
	sort.Sort(Kvs(param))
	for _, v := range param {
		ret += v.Key + v.Value
	}
	h := hmac.New(md5.New, []byte(jd.AppSecret))
	h.Write([]byte(ret))
	sha := hex.EncodeToString(h.Sum(nil))
	return strings.ToUpper(sha)

}

func (jd *JdUnion) request(param []*Kv) (string, error) {
	return "", nil
}
