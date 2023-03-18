package taoke

import (
	"context"

	"github.com/CodFrm/qqbot/command/alimama"
	"github.com/CodFrm/qqbot/utils/jdunion"
	"github.com/CodFrm/qqbot/utils/taobaoopen"
	"github.com/codfrm/cago/configs"
	"github.com/codfrm/cago/pkg/logger"
	zero "github.com/wdvxdr1123/ZeroBot"
	"go.uber.org/zap"
)

type TaoKe struct {
	repo *repo
	tb   *taobaoopen.Taobao
	jd   *jdunion.JdUnion
}

type Config struct {
	Taobao taobaoopen.TaobaoConfig `yaml:"taobao"`
	Jd     jdunion.JdUnion         `yaml:"jd"`
}

func NewTaoKe() (*TaoKe, error) {
	cfg := &Config{}
	if err := configs.Default().Scan("taoke", cfg); err != nil {
		return nil, err
	}
	alimama.Tb = taobaoopen.NewTaobao(cfg.Taobao)
	return &TaoKe{
		repo: &repo{},
		tb:   alimama.Tb,
		jd:   jdunion.NewJdUnion(cfg.Jd),
	}, nil
}

func (t *TaoKe) Init(ctx context.Context) error {
	t.forward()
	// admin 管理
	t.admin()
	zero.OnCommand("帮助", zero.OnlyPrivate).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		ctx.Send("1.优惠购物,触发指令:'有无[物品名]',可获取商品列表和内部优惠券,选择你心爱的物品下单吧\n" +
			"2.优惠查券,触发指令:'[口令或者链接]',可查询优惠券和优惠口令\n" +
			"3.订阅商品,触发指令:'订阅[关键字]',订阅指定的商品,如果有活动会直接私聊推送给你哦\n" +
			"饿了么每日红包:$5YiUccAeTlY$")
	})
	zero.OnCommand("有无", zero.OnlyPrivate).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		str, err := alimama.Search(ctx.State["args"].(string))
		if err != nil {
			logger.Default().Error("查询失败", zap.Error(err), zap.String("message", ctx.MessageString()))
			str = "查询失败: " + err.Error()
		}
		ctx.Send(str + "\n另可以发送'订阅" + ctx.State["args"].(string) + "'来关注本类商品哦")
	})
	zero.OnMessage(zero.OnlyPrivate).SetPriority(100).Handle(func(ctx *zero.Ctx) {
		ctx.Send("(私聊功能不稳定,请加好友发送'帮助'查看命令)")
	})
	return nil
}
