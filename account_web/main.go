package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"mic-study/account_web/handler"
)

func main() {
	ip := flag.String("ip", "127.0.0.1", "输入IP")
	port := flag.Int("port", 8080, "输入端口")
	flag.Parse()
	addr := fmt.Sprintf("%s:%d", *ip, *port)
	r := gin.Default()
	accountGroup := r.Group("/v1/account")
	{
		accountGroup.GET("/list", handler.AccountListHandler)
		accountGroup.POST("/login", handler.LoginByPasswordHandler)
		accountGroup.POST("/captcha", handler.CaptchaHandler)
	}
	r.Run(addr)
}
