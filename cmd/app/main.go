package main

import (
	"context"
	"log"

	"github.com/CodFrm/qqbot/internal/bot"
	"github.com/codfrm/cago"
	"github.com/codfrm/cago/configs"
	"github.com/codfrm/cago/database/redis"
	"github.com/codfrm/cago/pkg/logger"
)

func main() {
	ctx := context.Background()
	cfg, err := configs.NewConfig("qqbot")
	if err != nil {
		log.Fatalf("load config err: %v", err)
	}

	err = cago.New(ctx, cfg).
		Registry(cago.FuncComponent(logger.Logger)).
		Registry(cago.FuncComponent(redis.Redis)).
		RegistryCancel(bot.Bot()).
		Start()
	if err != nil {
		log.Fatalf("start err: %v", err)
		return
	}

}
