package main

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// 初始化容器
	store, _ := redis.NewStore(10, "tcp", "127.0.0.1:6379", "", []byte("bj39"))

	// 设置临时session
	/*	store.Options(sessions.Options{
		MaxAge:0,
	})*/

	// 使用容器
	router.Use(sessions.Sessions("mysession", store)) // mysession 是 cookie的名字

	router.GET("/test", func(context *gin.Context) {
		// 调用session, 设置session数据
		s := sessions.Default(context)
		// 设置session
		// s.Set("wang", "qichao")  // 这是 session 的 key和value
		// 修改session时候 需要Save函数配合 否则不生效
		// s.Save()

		v := s.Get("wang")
		fmt.Println("获取Session:", v.(string))

		context.Writer.WriteString("测试Session...")
	})
	router.Run(":8089")
}
