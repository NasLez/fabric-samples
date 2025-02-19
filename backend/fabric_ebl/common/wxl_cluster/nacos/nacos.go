package nacos

import (
	"fabric_ebl/common/choose"
	"fabric_ebl/common/env"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

func NewNacosClient() (naming_client.INamingClient, error) {
	if env.IsProd() {
		return clients.NewNamingClient(
			vo.NacosClientParam{
				ClientConfig: &constant.ClientConfig{
					NamespaceId: choose.If(env.IsProd(), "public", "boe"),
					Username:    "nacos",
					Password:    "wxl5211314",
				},
				ServerConfigs: []constant.ServerConfig{
					*constant.NewServerConfig("wxl475.cn", 30898),
				},
			},
		)
	} else {
		return clients.NewNamingClient(
			vo.NacosClientParam{
				ClientConfig: &constant.ClientConfig{
					NamespaceId: "public",
					Username:    "nacos",
					Password:    "wxl5211314",
				},
				ServerConfigs: []constant.ServerConfig{
					*constant.NewServerConfig("127.0.0.1", 8848),
				},
			},
		)
	}
}
