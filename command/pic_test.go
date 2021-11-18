package command

import (
	"github.com/CodFrm/qqbot/model"
	"testing"
)

func TestRotatePic(t *testing.T) {
	RotatePic(nil, &model.PicInfo{Url: "http://gchat.qpic.cn/gchatpic_new/958139621/340882274-2534335053-3D004B15539B48286542A49B681AA9DE/0?vuin=3623637397&term=255&pictype=0"})
}

func TestBanUser(t *testing.T) {
	BlackList("501546035", "1", "86400")
}
