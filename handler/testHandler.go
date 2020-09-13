package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"time"
)

func Test(ctx *gin.Context) {
	rand.Seed(time.Now().UnixNano())
	_ = ctx.Query("file")

	pre := []string{"test1", "test2", "test3"}
	puf := []string{".txt", ".exe", ".py", ".php"}
	fileAmount := rand.Intn(3) // 最多10个
	folderAmount := rand.Intn(3)
	fileList := make([]string, fileAmount)
	folderList := make([]string, folderAmount)
	for i := 0; i < fileAmount; i++ {
		randomPre := pre[rand.Intn(len(pre))]
		randomPuf := puf[rand.Intn(len(pre))]
		fileList[i] = randomPre + randomPuf
	}
	for i := 0; i < folderAmount; i++ {
		randomPre := pre[rand.Intn(len(pre))]
		folderList[i] = randomPre
	}

	fmt.Println(1)
	ctx.JSON(200, gin.H{
		"folderList": folderList,
		"fileList":   fileList,
	})

}

func TestWithPost(ctx *gin.Context) {
	rand.Seed(time.Now().UnixNano())

	path := ctx.PostForm("path")
	fmt.Println(path)

	isFile := ctx.PostForm("is_file")
	if isFile == "1" {
		length := rand.Intn(100)
		str := "0123456789abcdefghijklmnopqrstuvwxyz"
		bytes := []byte(str)
		result := []byte{}
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		for i := 0; i < length; i++ {
			result = append(result, bytes[r.Intn(len(bytes))])
		}
		ctx.JSON(200, gin.H{
			"content": result,
		})
	} else {
		pre := []string{"test1", "test2", "test3"}
		puf := []string{".txt", ".exe", ".py", ".php"}
		fileAmount := rand.Intn(3) // 最多10个
		folderAmount := rand.Intn(3)
		fileList := make([]string, fileAmount)
		folderList := make([]string, folderAmount)
		for i := 0; i < fileAmount; i++ {
			randomPre := pre[rand.Intn(len(pre))]
			randomPuf := puf[rand.Intn(len(pre))]
			fileList[i] = randomPre + randomPuf
		}
		for i := 0; i < folderAmount; i++ {
			randomPre := pre[rand.Intn(len(pre))]
			folderList[i] = randomPre
		}

		fmt.Println(1)
		ctx.JSON(200, gin.H{
			"folderList": folderList,
			"fileList":   fileList,
		})
	}
}
