package utils

import (
	"bj39/web/proto/user"
	"github.com/asim/go-micro/plugins/registry/consul/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/registry"
)

func InitMicro() user.UserService {
	// Register consul
	reg := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{"127.0.0.1:8500"}
	})
	service := micro.NewService(
		micro.Registry(reg),
		micro.Name("User"),
		micro.Version("latest"),
	)

	// 初始化客户端
	microClient := user.NewUserService("User", service.Client())
	return microClient
}
