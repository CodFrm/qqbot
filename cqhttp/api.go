package cqhttp

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/CodFrm/qqbot/config"
)

func Post(api string, m map[string]interface{}) ([]byte, error) {
	t, _ := json.Marshal(m)
	resp, err := http.Post("http://"+config.AppConfig.Url+"/"+api, "application/json", bytes.NewBuffer(t))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, nil
}

func SendGuildChannelMsg(guild_id, channel_id int64, message string) error {
	_, err := Post("send_guild_channel_msg", map[string]interface{}{
		"guild_id":   guild_id,
		"channel_id": channel_id,
		"message":    message,
	})
	return err
}

func SendPrivateMsg(user_id, group_id int64, message string) error {
	_, err := Post("send_private_msg", map[string]interface{}{
		"user_id":  user_id,
		"group_id": group_id,
		"message":  message,
	})
	return err
}
