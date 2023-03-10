package bot

import (
	"context"

	"github.com/CodFrm/qqbot/internal/plugins"
	"github.com/codfrm/cago"
	"github.com/codfrm/cago/configs"
	"github.com/codfrm/cago/pkg/logger"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/driver"
	"go.uber.org/zap"
)

type Config struct {
	Uin       uint   `yaml:"uin"`
	CqAddress string `yaml:"cq-address"`
	Admin     int64  `yaml:"admin"`
}

type bot struct {
}

func Bot() cago.ComponentCancel {
	return &bot{}
}

func (b *bot) StartCancel(ctx context.Context, cancel context.CancelFunc, cfg *configs.Config) error {
	config := &Config{}
	if err := cfg.Scan("bot", config); err != nil {
		return err
	}
	err := plugins.InitPlugins(ctx)
	if err != nil {
		logger.Ctx(ctx).Error("init plugins error", zap.Error(err))
		return err
	}
	go zero.RunAndBlock(&zero.Config{
		NickName:      []string{"bot"},
		CommandPrefix: "/",
		SuperUsers:    []int64{config.Admin},
		Driver: []zero.Driver{
			// 正向 WS
			driver.NewWebSocketClient("ws://"+config.CqAddress, ""),
		},
	}, func() {
	})
	return nil
}

func (b *bot) Start(ctx context.Context, cfg *configs.Config) error {
	return nil
}

func (b *bot) CloseHandle() {
}
