package model

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

// 存储图片id 到redis数据库
func SaveImgCode(code, uuid string) error {
	// 1.链接redis
	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("redis Dial err:", err)
		return err
	}
	defer conn.Close()

	// 2.操作数据库
	_, err = conn.Do("setex", uuid, 60*5, code)
	return err
}
