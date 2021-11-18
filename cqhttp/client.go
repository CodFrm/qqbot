package cqhttp

import (
	"encoding/json"
	"fmt"

	"github.com/golang/glog"
	"github.com/gorilla/websocket"
)

type RecvMsg struct {
	Interval      int64  `json:"interval"`
	MetaEventType string `json:"meta_event_type"`
	PostType      string `json:"post_type"`
	SelfId        string `json:"self_id"`
}

type Client struct {
	msgCallback func(msg *MessageModel)
	addr        string
	conn        *websocket.Conn
}

func NewClient(addr string) *Client {
	return &Client{
		addr: addr,
	}
}

func (c *Client) OnMessage(callback func(msg *MessageModel)) {
	c.msgCallback = callback
}

func (c *Client) Start() error {
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s/", c.addr), nil)
	if err != nil {
		return err
	}
	c.conn = conn
	defer func() {
		err := c.conn.Close()
		if err != nil {
			glog.Errorf("close connect: %v", err)
		}
		c.conn = nil
	}()
	for {
		v := make(map[string]interface{})
		_, p, err := c.conn.ReadMessage()
		if err != nil {
			glog.Errorf("read message: %v", err)
			return err
		}
		if err := json.Unmarshal(p, &v); err != nil {
			glog.Errorf("unmarshal: %v", err)
			continue
		}
		postType, ok := v["post_type"].(string)
		if !ok {
			continue
		}
		if postType != "message" {
			continue
		}
		c.handlerMessage(p, v)
		glog.V(8).Infof("%v", v)
	}
	return nil
}

func (c *Client) handlerMessage(message []byte, m map[string]interface{}) {
	defer func() {
		if r := recover(); r != nil {
			glog.Errorf("handlerMessage: %v", r)
		}
	}()
	v, ok := m["message_type"].(string)
	if !ok {
		return
	}
	var msg Message
	switch v {
	case "guild":
		msg = &GuildMsg{}
	case "private":
		msg = &PrivateMsg{}
	default:
		return
	}
	if err := json.Unmarshal(message, msg); err != nil {
		glog.Errorf("unmarshal: %v", err)
		return
	}
	c.msgCallback(&MessageModel{Message: msg})
}
