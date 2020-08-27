package command

import "testing"

func TestGenWordCloud(t *testing.T) {
	GenWordCloud("data/group/614202391_2020_08_27.txt")
	cronGenWordCloud()
}
