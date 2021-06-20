package main

import (
	"getCaptcha/handler"
	pb "getCaptcha/proto"
	"github.com/asim/go-micro/plugins/registry/consul/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/logger"
)

func main() {
	// Register consul
	reg := consul.NewRegistry()
	service := micro.NewService(
		micro.Registry(reg),
		micro.Name("GetCaptcha"),
		micro.Version("latest"),
	)

	// Register Handler
	if err := pb.RegisterGetCaptchaHandler(service.Server(), new(handler.GetCaptcha)); err != nil {
		logger.Fatal(err)
	}

	// Run service
	if err := service.Run(); err != nil {
		logger.Fatal(err)
	}
}
