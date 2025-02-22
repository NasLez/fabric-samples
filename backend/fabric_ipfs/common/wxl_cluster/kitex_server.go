package wxl_cluster

import (
	"fabric_ipfs/common/env"
	"fabric_ipfs/common/wxl_cluster/kitex_middleware"
	"fabric_ipfs/common/wxl_cluster/nacos"
	"net"
	"strconv"

	"github.com/cloudwego/kitex/pkg/rpcinfo"
	server "github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/registry-nacos/registry"
)

type newServerFunc[T any] func(handler T, opts ...server.Option) server.Server

func NewServer[T any](newServer newServerFunc[T], handler T, serverName string, boePort int) {
	options := make([]server.Option, 0)
	// boe 环境下指定服务地址
	if env.IsBoe() {
		host := ""
		port := boePort
		addr, err := net.ResolveTCPAddr("tcp", buildHostPort(host, port))
		if err != nil {
			panic(err)
		}
		options = append(options, server.WithServiceAddr(addr))
	}

	// 注册服务到nacos
	cli, err := nacos.NewNacosClient()
	if err != nil {
		panic(err)
	}
	options = append(options, server.WithRegistry(registry.NewNacosRegistry(cli)))
	options = append(options, server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: serverName}))

	// 日志中间件
	options = append(options, server.WithMiddleware(kitex_middleware.ServerLogMiddleware))

	// 启动服务
	svr := newServer(handler, options...)
	if err = svr.Run(); err != nil {
		panic(err)
	}
}

func buildHostPort(host string, port int) string {
	return host + ":" + strconv.FormatInt(int64(port), 10)
}
