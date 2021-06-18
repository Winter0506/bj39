package main

import (
	"getCaptcha/handler"
	pb "getCaptcha/proto"
	"github.com/asim/go-micro/plugins/registry/consul/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/registry"
	"github.com/asim/go-micro/v3/util/log"
)

func main() {
	// Register consul
	reg := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{"127.0.0.1:8500"}
	})
	service := micro.NewService(
		micro.Registry(reg),
	)

	// Register Handler
	pb.RegisterGetCaptchaHandler(service.Server(), new(handler.GetCaptcha))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
