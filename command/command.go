package command

import (
	"context"
	"github.com/CodFrm/iotqq-plugins/config"
	"github.com/mzz2017/shadowsocksR/client"
	pxy "github.com/nadoo/glider/proxy"
	"net"
)

var proxy func(ctx context.Context, network, addr string) (conn net.Conn, err error)

func Init() error {
	dia, err := client.NewSSRDialer(config.AppConfig.Ssr, pxy.Default)
	if err != nil {
		return err
	}
	proxy = func(ctx context.Context, network, addr string) (conn net.Conn, err error) {
		return dia.Dial(network, addr)
	}
	return nil
}
