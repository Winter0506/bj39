package main

import (
	"bj39/web/controller"
	"github.com/gin-gonic/gin"
)

// 添加gin框架开发3步骤
func main() {
	// 初始化路由
	router := gin.Default()

	// 路由匹配
	/*	router.GET("/", func(context *gin.Context) {
		context.Writer.WriteString("项目开始了....")
	})*/
	router.Static("/home", "view")

	router.GET("/api/v1.0/session", controller.GetSession)
	router.GET("/api/v1.0/imagecode/:uuid", controller.GetImageCd)

	// 启动运行
	router.Run(":8080")
}
