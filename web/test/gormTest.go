package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql" //"_" 代码不直接使用包, 底层链接要使用!
	"github.com/jinzhu/gorm"
)

// 创建全局结构体
type Student struct {
	Id   int
	Name string
	Age  int
}

func main() {

	conn, err := gorm.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/go_test")
	if err != nil {
		fmt.Println("gorm.Open err:", err)
		return
	}
	defer conn.Close()
	// 不要复数表名
	conn.SingularTable(true)
	// 借助gorm创建数据库表
	fmt.Println(conn.AutoMigrate(new(Student)).Error)
}
