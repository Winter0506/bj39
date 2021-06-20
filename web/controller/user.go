package controller

import (
	"bj39/web/model"
	getCaptcha "bj39/web/proto"
	"bj39/web/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/afocus/captcha"
	"github.com/asim/go-micro/plugins/registry/consul/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/registry"
	"github.com/gin-gonic/gin"
	"image/png"
	"math/rand"
	"net/http"
	"time"
)

func GetSession(ctx *gin.Context) {
	// 初始化错误返回的 map
	resp := make(map[string]string)

	resp["errno"] = utils.RECODE_SESSIONERR
	resp["errmsg"] = utils.RecodeText(utils.RECODE_SESSIONERR)

	ctx.JSON(http.StatusOK, resp)
}

func GetImageCd(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	// Register consul
	reg := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{"127.0.0.1:8500"}
	})
	service := micro.NewService(
		micro.Registry(reg),
		micro.Name("GetCaptcha"),
		micro.Version("latest"),
	)

	// 初始化客户端
	microClient := getCaptcha.NewGetCaptchaService("GetCaptcha", service.Client())

	// 调用远程函数
	resp, err := microClient.Call(context.TODO(), &getCaptcha.Request{Uuid: uuid})
	if err != nil {
		fmt.Println("未找到远程服务...")
		return
	}

	// 将得到的数据,反序列化,得到图片数据
	var img captcha.Image
	json.Unmarshal(resp.Img, &img)

	// 将图片写出到 浏览器.
	png.Encode(ctx.Writer, img)

	fmt.Println("uuid = ", uuid)
}

// 获取短信验证码
func GetSmscd(ctx *gin.Context) {
	// 获取短信验证码
	phone := ctx.Param("phone")
	// 拆分GET请求中的URL格式
	imgCode := ctx.Query("text")
	uuid := ctx.Query("id")

	fmt.Println(phone, imgCode, uuid)

	// 校验图片验证码 是否正确
	result := model.CheckImgCode(uuid, imgCode)
	// 创建容器 存储返回的数据信息
	resp := make(map[string]string)
	if result {
		// 模拟发送短息
		// 生成一个随机6位数 做验证码
		rand.Seed(time.Now().UnixNano())
		smsCode := fmt.Sprintf("%06d", rand.Int31n(1000000))
		fmt.Printf("验证码: %s\n", smsCode)

		// 发送短信验证码成功
		resp["errno"] = utils.RECODE_OK
		resp["errmsg"] = utils.RecodeText(utils.RECODE_OK)

		// 将 电话号:短信验证码 存入到redis数据库中
		err := model.SaveSmsCode(phone, smsCode)
		if err != nil {
			fmt.Println("存储短信验证码到redis失败:", err)
			resp["errno"] = utils.RECODE_DBERR
			resp["errmsg"] = utils.RecodeText(utils.RECODE_DBERR)
		}
	} else {
		// 校验失败 发送错误信息
		resp["errno"] = utils.RECODE_DATAERR
		resp["errmsg"] = utils.RecodeText(utils.RECODE_DATAERR)
	}
	// 发送成功/失败 结果
	ctx.JSON(http.StatusOK, resp)
}
