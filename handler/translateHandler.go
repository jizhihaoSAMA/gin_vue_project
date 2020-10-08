package handler

import (
	"context"
	"fmt"
	"gin_vue_project/rpcService/translate"
	"gin_vue_project/utils"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"log"
	"time"
)

func Translate(ctx *gin.Context) {
	fmt.Println(ctx.PostForm("selected_text"))
	conn, err := grpc.Dial("localhost:"+viper.GetString("port.grpcServer"), grpc.WithInsecure(), grpc.WithInsecure())
	if err != nil {
		log.Fatal("can not connect", err)
	}
	defer conn.Close()

	client := translate.NewTestClient(conn)

	context, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	logger := utils.GetLoggerWithTimeAndLine()

	request, err := client.Translate(context, &translate.TranslateRequest{
		OriginText: ctx.PostForm("selected_text"),
	})

	if err != nil {
		logger.Fatalf("err is :%v", err)
	}

	ctx.JSON(200, gin.H{
		"MsgCode": request.GetMsgCode(),
		"Result":  request.GetTranslateResult(),
	})
}
