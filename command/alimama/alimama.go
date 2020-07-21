package alimama

import (
	"github.com/CodFrm/iotqq-plugins/utils/iotqq"
	"github.com/robfig/cron/v3"
)

func Init() error {
	c := cron.New(cron.WithSeconds())
	c.AddFunc("0 30 11 * * ?", notice)
	c.AddFunc("0 30 17 * * ?", notice)
	c.Start()
	return nil
}

func notice() {
	list, err := iotqq.GetGroupList()
	if err != nil {
		return
	}
	for _, v := range list {
		iotqq.QueueSendMsg(v.GroupId, 0, "快到饭点了,来一份外卖吧~\nhttps://m.tb.cn/h.VsVRwwj\n复制这条信息，$nH3n1zNqDip$，到【手机淘宝】即可查看")
	}
}
