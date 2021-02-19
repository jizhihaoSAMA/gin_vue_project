package rpc

import (
	"context"
	"fmt"
	"gin_vue_project/response"
	"gin_vue_project/rpcService/translate"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"time"
)

func TranslateHandler(ctx *gin.Context) {
	fmt.Println(ctx.PostForm("selected_text"))
	conn, err := grpc.Dial("localhost:"+viper.GetString("port.grpcServer"), grpc.WithInsecure(), grpc.WithInsecure())
	if err != nil {
		//log.Fatal("can not connect", err)
		response.Fail(ctx, nil, "RPC异常")
	}
	defer conn.Close()

	client := translate.NewTestClient(conn)

	context, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	//logger := utils.GetLoggerWithTimeAndLine()

	request, err := client.Translate(context, &translate.TranslateRequest{
		OriginText: ctx.PostForm("selected_text"),
	})

	if err != nil {
		//logger.Fatalf("err is :%v", err)
		response.Fail(ctx, nil, "RPC异常")
	}

	ctx.JSON(200, gin.H{
		"MsgCode": request.GetMsgCode(),
		"Result":  request.GetTranslateResult(),
	})
}
