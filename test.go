package main

import (
    "fmt"
    "github.com/gin-gonic/gin"
    "log"
    "math/rand"
    "net/http"
    "time"
    "github.com/jinzhu/gorm"
)

func InitDB() *gorm.DB {
    driverName := "mysql"
    host := "localhost"
    port := "3306"
    database := "try"
    username := "root"
    password := "Zz7722559"
    charset := "utf8"
    args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
        username,
        password,
        host,
        port,
        database,
        charset,
    )

    db, err := gorm.Open(driverName, args)
    if err != nil {
        panic("链接失败，错误:" + err.Error())
    }
    return db

}

type User struct {
    gorm.Model
    Name      string `gorm:"type:varchar(20);not null"`
    Telephone string `gorm:"varchar(110);not null;unique"`
    Password  string `gorm:"size:255;not null"`
}

func main() {
    db := InitDB()
    defer db.Close()
    r := gin.Default()
    r.GET("/api/auth/register", func(ctx *gin.Context) {
        name := ctx.PostForm("name")
        telephone := ctx.PostForm("telephone")
        password := ctx.PostForm("password")
        // 检查手机号是否小于11位
        if len(telephone) != 11 {
            ctx.JSON(http.StatusUnprocessableEntity, gin.H{
                "code": 422,
                "msg":  "手机号必须为11位",
            })
        }
        // 检查密码是否少于6
        if len(password) < 6 {
            ctx.JSON(http.StatusUnprocessableEntity, gin.H{
                "code": 422,
                "msg":  "密码不能少于6位",
            })
        }
        // 如果没有设置名称。则自动设置名称
        if len(name) == 0 {
            name = RandomString(10)
        }

        if isTelephoneExist() {

        }

        log.Println(name, password, telephone)
    })
}

func isTelephoneExist(db *gorm.DB, telephone string) bool {
    // 初始化定义一个User
    var user User
    // 将结果传给user变量
    db.Where("telephone = ?", telephone).First(&user)
    if user.ID != 0 { // 0是初始值，判断是否为0则判断是否存在查询结果
        return true
    }
    return false
}

func RandomString(n int) string {
    var letters = []byte("asdfghjklzxcvbnmqwertyuiopASDFGHJKLZXCVBNMQWERTYUIOP")
    result := make([]byte, n)

    rand.Seed(time.Now().Unix())
    for i := range result {
        result[i] = letters[rand.Intn(len(letters))]
    }
    return result
}
