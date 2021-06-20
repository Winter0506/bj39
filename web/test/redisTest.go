package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

func main() {
	// 1.链接redis
	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("redis Dial err:", err)
		return
	}
	defer conn.Close()

	// 2.操作数据库
	reply, err := conn.Do("set", "xiaowang", "heihei")
	// 3.回复助手类 确定具体的数据类型
	r, err := redis.String(reply, err)

	fmt.Println(r, err)
}
