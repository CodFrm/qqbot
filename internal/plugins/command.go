package plugins

import (
	"context"

	"github.com/CodFrm/qqbot/internal/plugins/taoke"
)

func InitPlugins(ctx context.Context) error {
	taoke, err := taoke.NewTaoKe()
	if err != nil {
		return err
	}
	if err := taoke.Init(ctx); err != nil {
		return err
	}

	return nil
}
