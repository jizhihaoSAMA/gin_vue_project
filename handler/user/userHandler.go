package user

import (
	"fmt"
	"gin_vue_project/common"
	"gin_vue_project/dto"
	"gin_vue_project/model"
	"gin_vue_project/response"
	"gin_vue_project/utils"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

func GetInfoHandler(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	ctx.JSON(200, gin.H{
		"code": 200,
		"data": gin.H{
			"user": user,
		},
	})
}

func LoginHandler(ctx *gin.Context) {
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

func RegisterHandler(ctx *gin.Context) {
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

func UpdateTokenHandler(ctx *gin.Context) { // token正确的情况下，返回新token
	u, _ := ctx.Get("user")
	userDTO, _ := (u).(dto.UserDto)
	mysqlDB := common.InitMySQL()
	defer mysqlDB.Close()

	var user model.User
	mysqlDB.Where("id = ?", userDTO.ID).First(&user)
	token, err := common.GetToken(user)
	if err != nil {
		log.Println("发放token失败")
		response.ServerError(ctx, nil, "系统错误")
	}
	response.Success(ctx, gin.H{"token": token}, "更新成功")
}

func UpdateInfoHandler(ctx *gin.Context) {
	u, _ := ctx.Get("user")
	p, _ := (u).(dto.UserDto)

	updateDetail := ctx.PostForm("updatedDetail")
	updateUsername := ctx.PostForm("updatedUsername")

	mysqlDB := common.InitMySQL()
	defer mysqlDB.Close()

	var user model.User
	mysqlDB.Where("id = ?", p.ID).First(&user)

	if user.ID == 0 {
		response.ServerError(ctx, nil, "服务器错误")
	} else {
		user.Detail = updateDetail
		user.Username = updateUsername
		mysqlDB.Model(&user).Update(model.User{Detail: updateDetail, Username: updateUsername})
		response.Success(ctx, gin.H{"user": dto.ToUserDto(user)}, "更新成功")
	}
}

func UploadIconHandler(ctx *gin.Context) {
	u, _ := ctx.Get("user")
	user, ok := (u).(dto.UserDto)
	if ok {
		file, err := ctx.FormFile("updateIcon") // 图片上传
		if err != nil {
			response.Fail(ctx, nil, "错误")
			return
		}
		src, _ := file.Open()
		defer src.Close()

		data := make([]byte, file.Size)
		if _, err := src.Read(data); err != nil && err != io.EOF {
			fmt.Println(err.Error())
		}
		uploadedHeader := http.DetectContentType(data)
		//log.Println(uploadedHeader)
		if suffix := utils.FormatMapper(utils.UploadUserIcon, uploadedHeader); suffix == "" { // 上传的文件找不到对应的映射
			response.Fail(ctx, nil, "传输格式有误")
			return
		} else {
			path := "static/userInfo/userIcon/"
			err := os.MkdirAll(path, os.ModePerm)
			if err != nil {
				log.Println("Error occurs when create file :", err.Error())
			}

			res, err := utils.ToPNG(data, suffix)
			if err != nil {
				log.Println("Error occurs when translate image :", err.Error())
			}
			if err := ioutil.WriteFile(path+strconv.Itoa(int(user.ID))+".png", res, 0644); err != nil {
				log.Println("Error occurs when saving image :", err.Error())
				response.ServerError(ctx, nil, "服务器错误")
				return
			} else {
				response.Success(ctx, nil, "上传成功")
			}

			//if err := utils.SaveUploadIcon(file, path+"userID_"+strconv.Itoa(int(user.ID))+suffix); err != nil{
			//    log.Println(err.Error())
			//    response.ServerError(ctx, nil, "系统错误")
			//}else{
			//    response.Success(ctx, nil, "上传成功")
			//}

		}

	} else {
		response.Fail(ctx, nil, "错误")
	}

}

func getUnreadMessage(ctx *gin.Context) {

}
