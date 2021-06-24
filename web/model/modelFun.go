package model

import (
	"crypto/md5"
	"encoding/hex"
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

// 处理登录业务 根据手机号/密码 获取用户名
func Login(mobile, pwd string) (string, error) {
	var user User
	// 对参数pwd做md5 hash
	m5 := md5.New()
	m5.Write([]byte(pwd))
	pwdHash := hex.EncodeToString(m5.Sum(nil))

	err := GlobalConn.Select("name").Where("mobile = ?", mobile).
		Where("password_hash = ?", pwdHash).Find(&user).Error

	return user.Name, err
}

// 获取用户信息
func GetUserInfo(userName string) (User, error) {
	// 实现SQL: select * from user where name = userName;
	var user User
	err := GlobalConn.First(&user).Where("name = ?", userName).Error
	return user, err
}

// 更新用户名
func UpdateUserName(newName, oldName string) error {
	// update user set name = 'itcast' where name = 旧用户名
	return GlobalConn.Model(new(User)).Where("name = ?", oldName).
		Update("name", newName).Error
}
