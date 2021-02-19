package user

import "github.com/gin-gonic/gin"

func GetSomeNotices(ctx *gin.Context) {
	// 从列表中拉取最近10个通知。

}

func GetUnreadAmount(ctx *gin.Context) {
	// 获取未读取的数量
}

func GetAllNotes(ctx *gin.Context) {
	// 用户点击查看所有通知来查看所有的通知。
}
