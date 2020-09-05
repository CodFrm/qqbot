package alimama

import (
	"github.com/CodFrm/iotqq-plugins/db"
	"github.com/CodFrm/iotqq-plugins/utils/iotqq"
	"strconv"
	"strings"
)

type publisher struct {
	topic string
	tag   string
	param string
}

type subscribe struct {
	handler func(info string, keyword *publisher)
	param   string
}

type broker struct {
	m map[string]map[string]*subscribe
}

func NewBroker() *broker {
	return &broker{m: make(map[string]map[string]*subscribe)}
}

func (t *broker) publisher(content string) {
	for key, v := range t.m {
		if strings.Index(content, key) != -1 {
			for k, s := range v {
				if s == nil {
					continue
				}
				s.handler(content, &publisher{
					topic: key,
					tag:   k,
					param: s.param,
				})
			}
		}
	}
}

func (t *broker) subscribe(topic string, tag string, handler *subscribe) {
	_, ok := t.m[topic]
	if !ok {
		t.m[topic] = make(map[string]*subscribe)
	}
	t.m[topic][tag] = handler
}

func (t *broker) unsubscribe(topic string, tag string) {
	_, ok := t.m[topic]
	if !ok {
		t.m[topic] = make(map[string]*subscribe)
	}
	t.m[topic][tag] = nil
}

func Subscribe(group int, qq int64, topic string) error {
	sgroup := strconv.Itoa(group)
	sqq := strconv.FormatInt(qq, 10)
	_, err := db.Redis.HSet("alimama:subscribe:list", sqq, sgroup).Result()
	if err != nil {
		return err
	}
	_, err = db.Redis.SAdd("alimama:subscribe:topic:"+sqq, topic).Result()
	if err != nil {
		return err
	}
	mq.subscribe(topic, sqq, &subscribe{
		handler: func(info string, keyword *publisher) {
			group, _ := strconv.ParseInt(keyword.param, 10, 64)
			qq, _ := strconv.ParseInt(keyword.tag, 10, 64)
			iotqq.QueueSendPrivateMsg(int(group), qq, "您关注的'"+topic+"'有新消息\n"+info+"\n回复'退订"+keyword.topic+"'可退订指定消息.不加关键字退订全部")
		},
		param: sgroup,
	})
	return nil
}

func UnSubscribe(qq int64, topic string) error {
	sqq := strconv.FormatInt(qq, 10)
	var err error
	var list []string
	if topic == "" {
		list, err = db.Redis.SMembers("alimama:subscribe:topic:" + sqq).Result()
		if err != nil {
			return err
		}
	} else {
		list = append(list, topic)
	}
	for _, v := range list {
		mq.unsubscribe(v, sqq)
		db.Redis.SRem("alimama:subscribe:topic:"+sqq, v)
	}
	return nil
}
