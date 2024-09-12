package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"mic-study/account_web/handler"
	"mic-study/internal"
)

func init() {
	err := internal.RegisterService(internal.AppConf.AccountWebConfig.Host,
		internal.AppConf.AccountWebConfig.SrvName,
		internal.AppConf.AccountWebConfig.SrvName,
		internal.AppConf.AccountWebConfig.Port,
		internal.AppConf.AccountWebConfig.Tags)
	if err != nil {
		panic(err)
	}
}

func main() {
	ip := flag.String("ip", "127.0.0.1", "输入IP")
	port := flag.Int("port", 8081, "输入端口")
	flag.Parse()
	addr := fmt.Sprintf("%s:%d", *ip, *port)
	r := gin.Default()
	accountGroup := r.Group("/v1/account")
	{
		accountGroup.GET("/list", handler.AccountListHandler)
		accountGroup.POST("/login", handler.LoginByPasswordHandler)
		accountGroup.POST("/captcha", handler.CaptchaHandler)
	}
	r.GET("/health", handler.HealthHandler)
	r.Run(addr)
}
