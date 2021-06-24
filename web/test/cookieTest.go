package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/test", func(context *gin.Context) {
		// 设置Cookie
		// context.SetCookie("testCookie", "xiaowang", 0, "", "", false, true)
		// 获取Cookie
		cookieVal, _ := context.Cookie("testCookie")
		fmt.Println("获取到的Cookie为:", cookieVal)
		context.Writer.WriteString("测试Cookie...")
	})
	router.Run(":8089")
}
