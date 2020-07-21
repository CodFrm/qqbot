package iotqq

import (
	"math/rand"
	"time"
)

type GroupMsg struct {
	Qqgroup int
	At      int64
	Content string
}

var groupQueue chan *GroupMsg

func init() {
	groupQueue = make(chan *GroupMsg, 100)
	go sendQueueMsg()
}

func sendQueueMsg() {
	for {
		m := <-groupQueue
		SendMsg(m.Qqgroup, m.At, m.Content)
		time.Sleep(time.Second * time.Duration(rand.Intn(3)+1))
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
