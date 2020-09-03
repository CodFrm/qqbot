package jdunion

import (
	"testing"
)

func TestNewJdUnion(t *testing.T) {
	jd := NewJdUnion(JdConfig{
		AppKey:    "b90ba1d5e3d0bbd982ab336a8eb47426",
		AppSecret: "1f1cb77e89ff4899ac1f0dd2b517bc75",
		SiteId:    "4000342575",
	})
	jd.GetPromotionLink("https://item.jd.com/10618612688.html")
}
