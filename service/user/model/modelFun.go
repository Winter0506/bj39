package model

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
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

// 校验短信验证码 --redis
func CheckSmsCode(phone, code string) error {
	// 链接redis
	conn := RedisPool.Get()
	// 从redis中, 根据key获取Value
	smsCode, err := redis.String(conn.Do("get", phone+"_code"))
	if err != nil {
		fmt.Println("redis get phone_code err:", err)
		return err
	}
	// 验证码匹配
	if smsCode != code {
		return errors.New("验证码匹配失败!")
	}
	// 匹配成功
	return nil
}

// 注册用户信息 写mysql数据库
func RegisterUser(mobile, pwd string) error {
	var user User
	user.Name = mobile // 默认使用手机号作为用户名
	// 使用md5对pwd加密
	m5 := md5.New()                             // 初始md5对象
	m5.Write([]byte(pwd))                       // 将pwd写入缓冲区
	pwd_hash := hex.EncodeToString(m5.Sum(nil)) // 不使用额外的密钥

	user.Password_hash = pwd_hash
	// 插入数据到mysql
	fmt.Println(user)
	return GlobalConn.Create(&user).Error
}
