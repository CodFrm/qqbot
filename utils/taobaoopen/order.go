package taobaoopen

import "time"

type TaobaoOrderRebate struct {
	tb     *Taobao
	config RebateConfig
	store  RebateStore
}

type RebateConfig interface {
	LastTime() time.Time
}

type RebateStore interface {
}

func NewTaobaoOrderFl(tb *Taobao, config RebateConfig, store RebateStore) *TaobaoOrderRebate {
	ret := &TaobaoOrderRebate{tb: tb, config: config, store: store}
	go ret.watch()
	return ret
}

func (t *TaobaoOrderRebate) watch() {
	for {
		t.tb.QueryOrder(func(options *OrderOptions) {

		})
		time.Sleep(time.Minute * 15)
	}
}
