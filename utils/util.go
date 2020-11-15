package utils

import (
	"fmt"
	"gin_vue_project/dto"
	"github.com/cloopen/go-sms-sdk/cloopen"
	"github.com/spf13/viper"
	"log"
	"math/rand"
	"time"
)

func RandomString(n int) string {
	var letters = []byte("asdfghjklzxcvbnmqwertyuiopASDFGHJKLZXCVBNMQWERTYUIOP")
	result := make([]byte, n)

	rand.Seed(time.Now().Unix())
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

func InterfaceToUserDto(value interface{}) dto.UserDto {
	var u dto.UserDto
	switch v := value.(type) {
	case dto.UserDto:
		op, _ := value.(dto.UserDto)
		return op
	default:
		fmt.Println(v)
	}
	return u
}

func SendMessage(telephone string) string {

	cfg := cloopen.DefaultConfig().
		WithAPIAccount(viper.GetString("messageInfo.APIAccount")).
		// 主账号令牌 TOKEN,登陆云通讯网站后,可在控制台首页看到开发者主账号ACCOUNT SID和主账号令牌AUTH TOKEN
		WithAPIToken(viper.GetString("messageInfo.APIToken"))
	sms := cloopen.NewJsonClient(cfg).SMS()

	// 随机4位字符
	captchaCode := RandomString(4)

	// 发送短信
	input := &cloopen.SendRequest{
		// 应用的APPID
		AppId: viper.GetString("messageInfo.AppID"),
		// 手机号码
		To: telephone,
		// 模版ID
		TemplateId: viper.GetString("messageInfo.TemplateID"),
		// 模版变量内容 非必
		Datas: []string{captchaCode, "5"},
	}
	// 下发
	resp, err := sms.Send(input)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	log.Printf("Response MsgId: %s \n", resp.TemplateSMS.SmsMessageSid)
	return captchaCode

}

func GetLoggerWithTimeAndLine() *log.Logger {
	var logger log.Logger
	logger.SetFlags(log.Ldate | log.Lshortfile)
	return &logger
}
