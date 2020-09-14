package taobaoopen

import (
	"log"
	"time"
)

type TaobaoOrderRebate struct {
	tb      *Taobao
	config  RebateConfig
	store   RebateStore
	queuecl chan []*OrderItem
}

type RebateConfig interface {
	LastTime() (time.Time, error)
	UpdateTime(t time.Time) error
}

type RebateStore interface {
	Store([]*OrderItem) error
}

func NewTaobaoOrderFl(tb *Taobao, config RebateConfig, store RebateStore) *TaobaoOrderRebate {
	ret := &TaobaoOrderRebate{tb: tb,
		config: config, store: store,
		queuecl: make(chan []*OrderItem),
	}
	go ret.watch()
	go ret.storeQueue()
	return ret
}

func (t *TaobaoOrderRebate) watch() {
	for {
		var tm time.Time
		var err error
		n := 1
		for {
			tm, err = t.config.LastTime()
			if err != nil {
				log.Println("TaobaoOrderRebate watch LastTime ", err.Error())
				time.Sleep(time.Second * time.Duration(n))
				n *= 2
				continue
			}
			break
		}
		//扫描15天前的订单和今天的订单
		allOrder := t.ScanOrder(tm, time.Minute*15)
		t.queuecl <- allOrder
		oldAllOrder := t.ScanOrder(tm, time.Minute*15)
		t.queuecl <- oldAllOrder
	}
}

func (t *TaobaoOrderRebate) storeQueue() {
	for {
		orders := <-t.queuecl
		if err := t.store.Store(orders); err != nil {
			log.Println("TaobaoOrderRebate storeQueue Store ", err.Error())
			t.queuecl <- orders
			continue
		}
	}
}

func (t *TaobaoOrderRebate) ScanOrder(tm time.Time, d time.Duration) []*OrderItem {
	var order []*OrderItem
	var resp *OrderQueryRespond
	var err error
	allOrder := make([]*OrderItem, 0)
	pageNo := 1
	n := 1
	for {
		order, resp, err = t.tb.QueryOrder(func(options *OrderOptions) {
			options.StartTime = tm.Add(-d).Format("2006-01-02 15:04:05")
			options.EndTime = tm.Format("2006-01-02 15:04:05")
			options.PageNo = pageNo
		})
		if err != nil {
			log.Println("TaobaoOrderRebate watch QueryOrder ", tm.Format("2006-01-02 15:04:05"), " ", err.Error())
			time.Sleep(time.Second * time.Duration(n))
			n *= 2
			continue
		}
		allOrder = append(allOrder, order...)
		if !resp.TbkScOrderDetailsGetResponse.Data.HasNext {
			break
		}
		n = 1
		pageNo += 1
	}
	return allOrder
}
