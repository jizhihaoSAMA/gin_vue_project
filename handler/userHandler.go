package handler

import (
	"fmt"
	"gin_vue_project/common"
	"gin_vue_project/model"
	"gin_vue_project/response"
	"gin_vue_project/utils"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"log"
)

func GetInfo(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	ctx.JSON(200, gin.H{
		"code": 200,
		"data": gin.H{
			"user": user,
		},
	})
}

func Login(ctx *gin.Context) {
	db := common.InitMySQL()
	defer db.Close()
	telephone := ctx.PostForm("telephone")
	password := ctx.PostForm("password")
	var user model.User

	// 判断手机号是否存在, 并将数据保存到user变量种
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID == 0 {
		response.Response(ctx, 422, 422, nil, "用户不存在")
		return
	}

	// 判断密码是否正确，用CompareHashAndPassword将数据库中的密码和用户提交的密码对比，如果有错误，则出现err
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		response.Fail(ctx, nil, "用户信息错误")
		return
	}

	// 发放token
	token, err := common.GetToken(user)
	if err != nil {
		response.Response(ctx, 500, 500, nil, "系统异常")
		log.Println("token generate error:" + err.Error())
		return
	}

	// 返回结果

	response.Success(ctx, gin.H{"token": token}, "登录成功")
}

func Register(ctx *gin.Context) {
	db := common.InitMySQL()
	defer db.Close()
	var u model.User
	// 获取参数
	ctx.Bind(&u)
	username := u.Username
	telephone := u.Telephone
	password := u.Password

	fmt.Println(username, telephone, password)
	// 检查手机号是否小于11位
	if len(telephone) != 11 {
		response.Response(ctx, 422, 422, nil, "手机号不得小于11位")
		return
	}
	// 检查密码是否少于6
	if len(password) < 6 {
		response.Response(ctx, 422, 422, nil, "密码不得小于6位")
		return
	}
	// 如果没有设置名称。则自动设置名称
	if len(username) == 0 {
		username = utils.RandomString(10)
	}
	// 判断用户是否存在
	if isTelephoneExist(db, telephone) {
		response.Response(ctx, 422, 422, nil, "用户已存在")
		return
	}
	// 不明文保存密码，加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		response.Response(ctx, 500, 500, nil, "系统加密错误")
		return
	}

	newUser := model.User{
		Username:  username,
		Telephone: telephone,
		Password:  string(hashedPassword),
	}

	// 验证码检测
	captcha := ctx.PostForm("captcha")
	rdb := common.InitRedis()
	defer rdb.Close()
	captchaInRedis, err := rdb.Get(telephone + "_reg").Result()

	if captchaInRedis == captcha { // 相同代表正确
		// 创建用户
		db.Create(&newUser)

		token, err := common.GetToken(newUser)

		if err != nil {
			response.Response(ctx, 500, 500, nil, "系统异常")
			log.Println("token generate error:" + err.Error())
			return
		}

		response.Success(ctx, gin.H{"token": token}, "注册成功")
		log.Println(username, password, telephone)
	} else {
		ctx.JSON(403, gin.H{
			"code": "403",
			"msg":  "验证码不正确",
		})
		return
	}
}

func isTelephoneExist(db *gorm.DB, telephone string) bool {
	// 初始化定义一个User
	var user model.User
	// 将结果传给user变量
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID != 0 { // 0是初始值，判断是否为0则判断是否存在查询结果
		return true
	}
	return false
}
