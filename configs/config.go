package configs

import (
	"strconv"

	"github.com/codfrm/cago/configs"
)

func AdminQQ() int64 {
	qq, _ := strconv.ParseInt(configs.Default().String("admin.qq"), 10, 64)
	return qq
}
