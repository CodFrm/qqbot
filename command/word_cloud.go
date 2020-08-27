package command

import (
	"bytes"
	"github.com/CodFrm/iotqq-plugins/utils"
	"github.com/CodFrm/iotqq-plugins/utils/iotqq"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

func GenWordCloud(file string) (string, error) {
	cmd := exec.Command("python3", "data/tmp/test.py", file, "data/tmp/tmp.png")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return utils.FileBase64("data/tmp/tmp.png"), nil
}

func cronGenWordCloud() {
	ls, err := ioutil.ReadDir("data/group")
	if err != nil {
		log.Print("文件生成失败")
		return
	}
	for _, v := range ls {
		if !v.IsDir() {
			s := strings.Split(v.Name(), "_")
			group := utils.StringToInt(s[0])
			img, err := GenWordCloud(v.Name())
			if err != nil {
				println("词云生成失败")
				continue
			}
			iotqq.SendPicByBase64(group, 0, "昨日词云", img)
			time.Sleep(time.Second * 10)
		}
	}
}
