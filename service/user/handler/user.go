package handler

import (
	"context"
	"fmt"
	"math/rand"
	"time"
	"user/model"
	"user/utils"

	user "user/proto"
)

type User struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *User) SendSms(ctx context.Context, req *user.Request, rsp *user.Response) error {

	// 校验图片验证码 是否正确
	result := model.CheckImgCode(req.Uuid, req.ImgCode)
	// 创建容器 存储返回的数据信息
	// resp := make(map[string]string)
	if result {
		// 模拟发送短息
		// 生成一个随机6位数 做验证码
		rand.Seed(time.Now().UnixNano())
		smsCode := fmt.Sprintf("%06d", rand.Int31n(1000000))
		fmt.Printf("验证码: %s\n", smsCode)

		// 发送短信验证码成功
		rsp.Errno = utils.RECODE_OK
		rsp.Errmsg = utils.RecodeText(utils.RECODE_OK)

		// 将 电话号:短信验证码 存入到redis数据库中
		err := model.SaveSmsCode(req.Phone, smsCode)
		if err != nil {
			fmt.Println("存储短信验证码到redis失败:", err)
			rsp.Errno = utils.RECODE_DBERR
			rsp.Errmsg = utils.RecodeText(utils.RECODE_DBERR)
		}
	} else {
		// 校验失败 发送错误信息
		rsp.Errno = utils.RECODE_DATAERR
		rsp.Errmsg = utils.RecodeText(utils.RECODE_DATAERR)
	}
	return nil
}

func (e *User) Register(ctx context.Context, req *user.RegReq, rsp *user.Response) error {
	// 校验短信验证码是否正确，redis中存储短信校验码
	err := model.CheckSmsCode(req.Mobile, req.SmsCode)
	if err == nil {
		// 如果校验正确,注册用户 将数据写入到mysql数据库
		err = model.RegisterUser(req.Mobile, req.Password)
		if err != nil {
			rsp.Errno = utils.RECODE_DBERR
			rsp.Errmsg = utils.RecodeText(utils.RECODE_DBERR)
		} else {
			rsp.Errno = utils.RECODE_OK
			rsp.Errmsg = utils.RecodeText(utils.RECODE_OK)
		}

	} else { // 短信验证码错误
		rsp.Errno = utils.RECODE_DATAERR
		rsp.Errmsg = utils.RecodeText(utils.RECODE_DATAERR)
	}
	return nil
}
