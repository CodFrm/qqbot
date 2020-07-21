package taobaoopen

import "strings"

type Kv struct {
	Key   string
	Value string
}

type TaobaoConfig struct {
	AppKey    string `yaml:"appKey"`
	AppSecret string `yaml:"appSecret"`
	Entrance  string
	AdzoneId  string `yaml:"adzoneId"`
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
