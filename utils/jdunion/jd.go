package jdunion

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/CodFrm/iotqq-plugins/utils"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

type JdUnion struct {
	AppKey    string `yaml:"appKey"`
	AppSecret string `yaml:"appSecret"`
	SiteId    string `yaml:"siteId"`
	XlAppKey  string `yaml:"xlAppKey"`
	JdId      string `yaml:"jdId"`
}

func NewJdUnion(config JdUnion) *JdUnion {
	return &config
}

type GetPromotionLinkRespond struct {
	JdUnionOpenPromotionCommonGetResponse struct {
		Result string `json:"result"`
		Code   string `json:"code"`
	} `json:"jd_union_open_promotion_common_get_response"`
}

type PromotionLink struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"requestId"`
}

func (jd *JdUnion) GetPromotionLink(materialId string) (*PromotionLink, error) {
	ret, err := jd.PublicFunc("jd.union.open.promotion.common.get", &Kv{
		Key:   "param_json",
		Value: fmt.Sprintf(`{"promotionCodeReq":{"siteId":"%s","materialId":"%s"}}`, jd.SiteId, materialId),
	})
	if err != nil {
		return nil, err
	}
	respond := &GetPromotionLinkRespond{}
	if err := json.Unmarshal([]byte(ret), respond); err != nil {
		return nil, err
	}
	retJson := &PromotionLink{}
	if err := json.Unmarshal([]byte(respond.JdUnionOpenPromotionCommonGetResponse.Result), retJson); err != nil {
		return nil, err
	}
	if retJson.Code != 0 {
		return nil, errors.New(retJson.Message)
	}
	return retJson, nil
}

func (jd *JdUnion) PromotionGoodsInfo(skuIds string) (*PromotionLink, error) {
	ret, err := jd.PublicFunc("jd.union.open.goods.promotiongoodsinfo.query", &Kv{
		Key:   "param_json",
		Value: fmt.Sprintf(`{"skuIds":"%s"}`, skuIds),
	})
	if err != nil {
		return nil, err
	}
	respond := &GetPromotionLinkRespond{}
	if err := json.Unmarshal([]byte(ret), respond); err != nil {
		return nil, err
	}
	retJson := &PromotionLink{}
	if err := json.Unmarshal([]byte(respond.JdUnionOpenPromotionCommonGetResponse.Result), retJson); err != nil {
		return nil, err
	}
	if retJson.Code != 0 {
		return nil, errors.New(retJson.Message)
	}
	return retJson, nil
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
	ret = jd.AppSecret + ret + jd.AppSecret
	ret = fmt.Sprintf("%x", md5.Sum([]byte(ret)))
	return strings.ToUpper(ret)
}

func (jd *JdUnion) request(param []*Kv) (string, error) {
	data := ""
	for _, v := range param {
		data += v.Key + "=" + url.QueryEscape(v.Value) + "&"
	}
	data = data[:len(data)-1]
	resp, err := http.Post("https://router.jd.com/api", "application/x-www-form-urlencoded;charset=utf-8", bytes.NewBuffer([]byte(data)))
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

type Link struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		ClickURL string `json:"clickURL"`
		ShortURL string `json:"shortURL"`
	} `json:"data"`
}

func (jd *JdUnion) ConversionLink(url string) (*Link, error) {
	data, err := utils.HttpGet("https://openapi.linkstars.com/api/changeurl_jd?apikey="+jd.XlAppKey+"&v=1.0.0&unionid="+jd.JdId+
		"&materialId="+url, nil, nil)
	if err != nil {
		return nil, err
	}
	ret := &Link{}
	if err := json.Unmarshal(data, ret); err != nil {
		return nil, err
	}
	return ret, nil
}
