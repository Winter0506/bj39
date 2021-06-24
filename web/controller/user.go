package controller

import (
	"bj39/web/model"
	"bj39/web/proto/getCaptcha"
	"bj39/web/proto/user"
	"bj39/web/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/afocus/captcha"
	"github.com/asim/go-micro/plugins/registry/consul/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/registry"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"image/png"
	"net/http"
)

/*func GetSession(ctx *gin.Context) {
	// 初始化错误返回的 map
	resp := make(map[string]string)

	resp["errno"] = utils.RECODE_SESSIONERR
	resp["errmsg"] = utils.RecodeText(utils.RECODE_SESSIONERR)

	ctx.JSON(http.StatusOK, resp)
}*/

// 获取session信息
func GetSession(ctx *gin.Context) {
	resp := make(map[string]interface{})

	// 获取Session数据
	s := sessions.Default(ctx) // 初始化Session对象
	userName := s.Get("userName")

	// 用户没有目录 ---没存在MySQL中, 也没存在Session中
	if userName == nil {
		resp["errno"] = utils.RECODE_SESSIONERR
		resp["errmsg"] = utils.RecodeText(utils.RECODE_SESSIONERR)
	} else {
		resp["errno"] = utils.RECODE_OK
		resp["errmsg"] = utils.RecodeText(utils.RECODE_OK)

		// 键是name 值是 姓名
		var nameData struct {
			Name string `json:"name"`
		}
		nameData.Name = userName.(string) // 类型断言
		resp["data"] = nameData
	}
	ctx.JSON(http.StatusOK, resp)
}

func GetImageCd(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	fmt.Println(uuid)
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
	// fmt.Println(resp)
	if err != nil {
		fmt.Println("未找到getCaptcha远程服务...")
		return
	}

	// 将得到的数据,反序列化,得到图片数据
	var img captcha.Image
	json.Unmarshal(resp.Img, &img)

	// 将图片写出到 浏览器.
	png.Encode(ctx.Writer, img)

	fmt.Println("uuid = ", uuid)
}

var microClient user.UserService

// 获取短信验证码
func GetSmscd(ctx *gin.Context) {
	// 获取短信验证码
	phone := ctx.Param("phone")
	// 拆分GET请求中的URL格式
	imgCode := ctx.Query("text")
	uuid := ctx.Query("id")

	fmt.Println(phone, imgCode, uuid)

	microClient = utils.InitMicro()
	// 调用远程函数
	resp, err := microClient.SendSms(context.TODO(), &user.Request{Phone: phone, ImgCode: imgCode, Uuid: uuid})
	if err != nil {
		fmt.Println("未找到User远程服务...")
		return
	}

	// 发送成功/失败 结果
	ctx.JSON(http.StatusOK, resp)
}

/*func GetSmscd(ctx *gin.Context) {
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
}*/

// 发送注册信息
func PostRet(ctx *gin.Context) {
	// 获取数据
	var regData struct {
		Mobile   string `json:"mobile"`
		PassWord string `json:"password"`
		SmsCode  string `json:"sms_code"`
	}
	ctx.Bind(&regData)
	// fmt.Println("获取到的数据为:", regData)

	// 初始化consul
	// Register consul
	/*reg := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{"127.0.0.1:8500"}
	})
	service := micro.NewService(
		micro.Registry(reg),
		micro.Name("User"),
		micro.Version("latest"),
	)

	// 初始化客户端
	microClient := user.NewUserService("User", service.Client())*/

	// 调用远程函数
	resp, err := microClient.Register(context.TODO(), &user.RegReq{
		Mobile:   regData.Mobile,
		SmsCode:  regData.SmsCode,
		Password: regData.PassWord,
	})
	if err != nil {
		fmt.Println("注册用户,找不到远程服务!", err)
		return
	}
	// 写给浏览器
	ctx.JSON(http.StatusOK, resp)
}

// 获取地域信息
func GetArea(ctx *gin.Context) {
	// 先从mysql中获取链接
	var areas []model.Area

	// 从缓存redis中, 获取数据
	conn := model.RedisPool.Get()
	// 当初使用“字节切片”存入 现在使用切片类型接收
	areaData, _ := redis.Bytes(conn.Do("get", "areaData"))
	// 没有从redis中获取到数据
	if len(areaData) == 0 {
		fmt.Println("从MySQL中获取数据...")
		model.GlobalConn.Find(&areas)
		// 把数据写入到redis中, 存储结构体序列化吼的json串
		areaBuf, _ := json.Marshal(areas)
		conn.Do("set", "areaData", areaBuf)
	} else {
		fmt.Println("从redis中获取数据...")
		// redis中有数据
		json.Unmarshal(areaData, &areas)
	}

	/*model.GlobalConn.Find(&areas)

	// 再把数据写入到redis中
	conn := model.RedisPool.Get()  // 获取链接
	conn.Do("set", "areaData", areas)*/

	resp := make(map[string]interface{})

	resp["errno"] = "0"
	resp["errmsg"] = utils.RecodeText(utils.RECODE_OK)
	resp["data"] = areas

	ctx.JSON(http.StatusOK, resp)
}

// 处理登录业务
func PostLogin(ctx *gin.Context) {
	// 获取前端数据
	var loginData struct {
		Mobile   string `json:"mobile"`
		PassWord string `json:"password"`
	}
	ctx.Bind(&loginData)

	resp := make(map[string]interface{})

	// 获取 数据库数据 查询是否和数据匹配
	userName, err := model.Login(loginData.Mobile, loginData.PassWord)
	if err == nil {
		// 登录成功!
		resp["errno"] = utils.RECODE_OK
		resp["errmsg"] = utils.RecodeText(utils.RECODE_OK)

		// 将登陆状态 保存到Session中
		s := sessions.Default(ctx)  // 初始化session
		s.Set("userName", userName) // 将用户名设置到session中
		s.Save()

	} else {
		// 登录失败!
		resp["errno"] = utils.RECODE_LOGINERR
		resp["errmsg"] = utils.RecodeText(utils.RECODE_LOGINERR)
	}

	ctx.JSON(http.StatusOK, resp)
}

func DeleteSession(ctx *gin.Context) {
	resp := make(map[string]interface{})

	// 初始化 Session对象
	s := sessions.Default(ctx)
	// 删除Session数据
	s.Delete("userName") // 没有返回值
	// 必须使用Save保存
	err := s.Save() // 有返回值

	if err != nil {
		resp["errno"] = utils.RECODE_IOERR // 没有合适错误,使用 IO 错误!
		resp["errmsg"] = utils.RecodeText(utils.RECODE_IOERR)

	} else {
		resp["errno"] = utils.RECODE_OK
		resp["errmsg"] = utils.RecodeText(utils.RECODE_OK)
	}

	ctx.JSON(http.StatusOK, resp)
}

// 获取用户基本信息
func GetUserInfo(ctx *gin.Context) {
	resp := make(map[string]interface{})

	defer ctx.JSON(http.StatusOK, resp)

	// 获取Session 得到当前用户信息
	s := sessions.Default(ctx)
	userName := s.Get("userName")

	// 判断用户名是否存在
	if userName == nil { // 用户没有登陆 但进入了该页面 恶意进入
		resp["errno"] = utils.RECODE_SESSIONERR
		resp["errmsg"] = utils.RecodeText(utils.RECODE_SESSIONERR)
		return // 如果出错, 报错, 退出
	}

	// 根据用户名 获取用户信息  --- 查MySQL数据库user表
	user, err := model.GetUserInfo(userName.(string))
	if err != nil {
		resp["errno"] = utils.RECODE_DBERR
		resp["errmsg"] = utils.RecodeText(utils.RECODE_DBERR)
		return // 如果出错, 报错, 退出
	}

	resp["errno"] = utils.RECODE_OK
	resp["errmsg"] = utils.RecodeText(utils.RECODE_OK)

	temp := make(map[string]interface{})
	temp["user_id"] = user.ID
	temp["name"] = user.Name
	temp["mobile"] = user.Mobile
	temp["real_name"] = user.Real_name
	temp["id_card"] = user.Id_card
	temp["avatar_url"] = user.Avatar_url

	resp["data"] = temp
}

// 更新用户名
func PutUserInfo(ctx *gin.Context) {
	// 获取当前用户名
	s := sessions.Default(ctx) // 初始化Session对象
	userName := s.Get("userName")

	// 获取新用户名  -- 处理Request Payload 类型数据 Bind()
	var nameData struct {
		Name string `json:"name"`
	}
	ctx.Bind(&nameData)

	// 更新用户名
	resp := make(map[string]interface{})
	defer ctx.JSON(http.StatusOK, resp)

	// 更新用户名
	err := model.UpdateUserName(nameData.Name, userName.(string))
	if err != nil {
		resp["errno"] = utils.RECODE_DBERR
		resp["errmsg"] = utils.RecodeText(utils.RECODE_DBERR)
		return
	}

	// 更新Session数据
	s.Set("userName", nameData.Name)
	err = s.Save() // 必须保存
	if err != nil {
		resp["errno"] = utils.RECODE_SESSIONERR
		resp["errmsg"] = utils.RecodeText(utils.RECODE_SESSIONERR)
		return
	}
	resp["errno"] = utils.RECODE_OK
	resp["errmsg"] = utils.RecodeText(utils.RECODE_OK)
	resp["data"] = nameData
}
