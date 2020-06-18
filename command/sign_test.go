package command

import (
	"testing"
	"time"
)

func Test_rewardKick(t *testing.T) {
	a := time.Now()
	b, _ := time.Parse("2006-01-02", "2020-06-12")
	d := a.Sub(b)
	println(d.Hours() / 24)
}
