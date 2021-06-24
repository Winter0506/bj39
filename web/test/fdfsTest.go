package main

import (
	"fmt"
	"github.com/tedcy/fdfs_client"
)

func main() {
	// 初始化客户端 --- 配置文件
	// 我这是windows没办法弄了
	clt, err := fdfs_client.NewClientWithConfig("/etc/fdfs/client.conf")
	if err != nil {
		fmt.Println("初始化客户端错误, err:", err)
		return
	}

	// 上传文件 --- 尝试文件名上传 传入到storage
	resp, err := clt.UploadByFilename("头像1.jpg")
	fmt.Println(resp, err)
}
