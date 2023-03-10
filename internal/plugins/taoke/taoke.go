package taoke

import (
	"context"

	"github.com/CodFrm/qqbot/utils/jdunion"
	"github.com/CodFrm/qqbot/utils/taobaoopen"
	"github.com/codfrm/cago/configs"
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
	return &TaoKe{
		repo: &repo{},
		tb:   taobaoopen.NewTaobao(cfg.Taobao),
		jd:   jdunion.NewJdUnion(cfg.Jd),
	}, nil
}

func (t *TaoKe) Init(ctx context.Context) error {
	t.forward()
	// admin 管理
	t.admin()
	return nil
}
