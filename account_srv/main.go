package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"mic-study/account_srv/biz"
	"mic-study/account_srv/proto/pb"
	"mic-study/internal"
	"mic-study/util"
	"net"
)

func init() {
	internal.InitDB()
}

func main() {
	//ip := flag.String("ip", "127.0.0.1", "输入IP")
	//port := flag.Int("port", 8080, "输入端口")
	//flag.Parse()
	//addr := fmt.Sprintf("%s:%d", *ip, *port)
	port := util.GenRandomPort()

	addr := fmt.Sprintf("%s:%d", internal.AppConf.AccountSrvConfig.Host, port)
	server := grpc.NewServer()
	pb.RegisterAccountServiceServer(server, &biz.AccountServer{})
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		zap.S().Error("account_srv启动异常:", zap.Error(err))
		panic(err)
	}
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())
	defaultConfig := api.DefaultConfig()
	defaultConfig.Address = fmt.Sprintf("&s:%d", internal.AppConf.ConsulConfig.Host, port)
	client, err := api.NewClient(defaultConfig)
	if err != nil {
		panic(err)
	}
	checkAddr := fmt.Sprintf("%s:%d", internal.AppConf.AccountSrvConfig.Host, port)
	check := &api.AgentServiceCheck{
		GRPC:                           checkAddr,
		Timeout:                        "3s",
		Interval:                       "1s",
		DeregisterCriticalServiceAfter: "5s",
	}
	randUUID := uuid.New().String()
	req := api.AgentServiceRegistration{
		Name:  internal.AppConf.AccountSrvConfig.SrvName,
		ID:    randUUID,
		Port:  port,
		Tags:  internal.AppConf.AccountSrvConfig.Tags,
		Check: check,
	}
	err = client.Agent().ServiceRegister(&req)
	if err != nil {
		panic(err)
	}
	fmt.Println(fmt.Sprintf("%s启动在%d", randUUID, port))
	err = server.Serve(listen)
	if err != nil {
		panic(err)
	}
}
