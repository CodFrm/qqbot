package iotqq

import (
	"math/rand"
	"time"
)

type GroupMsg struct {
	Qqgroup int
	At      int64
	Content string
	Url     string
	Private bool
}

var groupQueue chan *GroupMsg

func init() {
	groupQueue = make(chan *GroupMsg, 100)
	go sendQueueMsg()
}

func sendQueueMsg() {
	for {
		m := <-groupQueue
		if m.Private {
			if m.Qqgroup <= 0 {
				SendFriendMsg(m.At, m.Content)
			} else {
				SendPrivateMsg(m.Qqgroup, m.At, m.Content)
			}
		} else {
			if m.Url != "" {
				SendPicByUrl(m.Qqgroup, m.At, m.Content, m.Url)
			} else {
				SendMsg(m.Qqgroup, m.At, m.Content)
			}
		}
		time.Sleep(time.Second * time.Duration(rand.Intn(4)+1))
	}
}

func QueueSendMsg(qqgroup int, At int64, Content string) error {
	groupQueue <- &GroupMsg{
		Qqgroup: qqgroup,
		At:      At,
		Content: Content,
	}
	return nil
}

func QueueSendPicMsg(qqgroup int, At int64, Content string, Url string) error {
	groupQueue <- &GroupMsg{
		Qqgroup: qqgroup,
		At:      At,
		Content: Content,
		Url:     Url,
	}
	return nil
}

func QueueSendPrivateMsg(qqgroup int, qq int64, Content string) error {
	groupQueue <- &GroupMsg{
		Qqgroup: qqgroup,
		At:      qq,
		Content: Content,
		Private: true,
	}
	return nil
}
