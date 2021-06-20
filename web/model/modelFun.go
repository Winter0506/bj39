package model

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

// 创建全局redis连接池句柄
var RedisPool redis.Pool

// 创建函数 初始化Redis连接池
func InitRedis() {
	RedisPool = redis.Pool{
		MaxIdle:         20,
		MaxActive:       50,
		MaxConnLifetime: 60 * 5,
		IdleTimeout:     60,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "127.0.0.1:6379")
		},
	}
}

// 校验图片验证码
func CheckImgCode(uuid, imgCode string) bool {
	// 链接redis
	/*conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("redis.Dial err:", err)
		return false
	}*/
	conn := RedisPool.Get()
	defer conn.Close()

	// 查询redis数据
	code, err := redis.String(conn.Do("get", uuid))
	if err != nil {
		fmt.Println("查询错误 err:", err)
		return false
	}
	// 返回校验结果
	return code == imgCode
}

// 存储短信验证码
func SaveSmsCode(phone, code string) error {
	conn := RedisPool.Get()
	defer conn.Close()

	// 存储短信验证码到redis中
	_, err := conn.Do("setex", phone+"_code", 60*3, code)
	return err
}