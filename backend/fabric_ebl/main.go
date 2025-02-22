package main

import (
	"context"
	"fabric_ebl/common/reader"
	"fabric_ebl/repo"
	"fabric_ebl/sal/config"
	"fabric_ebl/sal/dao"
	"fabric_ebl/service"
	"github.com/wxl-server/idl_gen/kitex_gen/fabric_ebl/fabricebl"

	"fabric_ebl/common/render"
	"fabric_ebl/common/wxl_cluster"
	"github.com/bytedance/gopkg/util/logger"
	"go.uber.org/dig"
)

var (
	initCtx   = context.Background()
	container = dig.New()
)

func main() {
	initContainer()

	wxl_cluster.NewServer(fabricebl.NewServer, handler, "fabric_ebl", 8090)
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
		mustProvide(repo.NewFabricEblRepo)
	}

	// service
	{
		mustProvide(service.NewUserService)
		mustProvide(service.NewConnectService)
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
