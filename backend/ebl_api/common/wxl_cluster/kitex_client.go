package wxl_cluster

import (
	"ebl_api/common/wxl_cluster/kitex_middleware"
	"ebl_api/common/wxl_cluster/nacos"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/cloudwego/kitex/client"
	"github.com/kitex-contrib/registry-nacos/resolver"
)

type newClientFunc[T any] func(destService string, opts ...client.Option) (T, error)

func NewClient[T any](newClient newClientFunc[T], destService string) T {
	cli, err := nacos.NewNacosClient()
	if err != nil {
		logger.Errorf("[Init] clients.NewNamingClient failed, err = %v", err)
		panic(err)
	}
	destClient, err := newClient(
		destService,
		client.WithResolver(resolver.NewNacosResolver(cli)),
		client.WithMiddleware(kitex_middleware.ClientLogMiddleware),
	)
	if err != nil {
		logger.Errorf("[Init] minercore.NewClient failed, err = %v", err)
		panic(err)
	}
	return destClient
}
