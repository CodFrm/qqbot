package taobaoopen

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewTaobao(t *testing.T) {
	key := os.Getenv("TAOBAO_APPKEY")
	secret := os.Getenv("TAOBAO_APPSECRET")
	tb := NewTaobao(key, secret, "https://eco.taobao.com/router/rest")
	resp, err := tb.PublicFunc("taobao.tbk.order.details.get",
		GenKv("end_time", "2020-07-21 12:28:22"), GenKv("start_time", "2020-07-19 12:28:22"),
		GenKv("force_sensitive_param_fuzzy", "true"),
		GenKv("partner_id", "top-apitools"))
	assert.Nil(t, err)
	assert.NotNil(t, resp)
}
