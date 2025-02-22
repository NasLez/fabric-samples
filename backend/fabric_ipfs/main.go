package main

import (
	"context"
	"fabric_ipfs/common/reader"
	"fabric_ipfs/common/render"
	"fabric_ipfs/common/wxl_cluster"
	"fabric_ipfs/repo"
	"fabric_ipfs/sal/config"
	"fabric_ipfs/sal/dao"
	"fabric_ipfs/service"
	"github.com/bytedance/gopkg/util/logger"
	fabric_ipfs "github.com/wxl-server/idl_gen/kitex_gen/fabric_ipfs/fabricipfs"
	"go.uber.org/dig"
)

var (
	initCtx   = context.Background()
	container = dig.New()
)

func main() {
	initContainer()

	wxl_cluster.NewServer(fabric_ipfs.NewServer, handler, "fabric_ipfs", 8091)
}

func initContainer() {
	// context
	{
		mustProvide(func() context.Context { return initCtx })
	}

	// config
	{
		mustProvide(reader.InitAppConfig[config.AppConfig])
	}

	// db
	{
		mustInvoke(dao.InitDB)
	}

	// repo
	{
		mustProvide(repo.NewFabricIpfsRepo)
	}

	// service
	{
		mustProvide(service.NewFabricIpfsService)
	}

	// handler
	{
		mustInvoke(NewHandler)
	}
}

func mustProvide(constructor interface{}, opts ...dig.ProvideOption) {
	if err := container.Provide(constructor, opts...); err != nil {
		logger.Errorf("container provide failed, err = %v, constructor = %v", err, render.Render(constructor))
		panic(err)
	}
}

func mustInvoke(function interface{}, opts ...dig.InvokeOption) {
	if err := container.Invoke(function, opts...); err != nil {
		logger.Errorf("container invoke failed, err = %v, function = %v", err, render.Render(function))
		panic(err)
	}
}
