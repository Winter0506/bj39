package main

import (
	"github.com/asim/go-micro/plugins/registry/consul/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/logger"
	"user/handler"
	"user/model"
	pb "user/proto"
)

func main() {
	// 初始化 MySQL 连接池
	model.InitDB()
	// 初始化连接池
	model.InitRedis()
	// Register consul
	reg := consul.NewRegistry()
	srv := micro.NewService(
		micro.Registry(reg),
		micro.Name("User"),
		micro.Version("latest"),
	)

	// Register handler
	if err := pb.RegisterUserHandler(srv.Server(), new(handler.User)); err != nil {
		logger.Fatal(err)
	}

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
