package main

import (
	"flag"
	"fmt"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"mic-study/account_srv/biz"
	"mic-study/account_srv/proto/pb"
	"mic-study/internal"
	"net"
)

func init() {
	internal.InitDB()
}

func main() {
	ip := flag.String("ip", "127.0.0.1", "输入IP")
	port := flag.Int("port", 8080, "输入端口")
	flag.Parse()
	addr := fmt.Sprintf("%s:%d", *ip, *port)

	server := grpc.NewServer()
	pb.RegisterAccountServiceServer(server, &biz.AccountServer{})
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		zap.S().Error("account_srv启动异常:", zap.Error(err))
		panic(err)
	}
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())
	defaultConfig := api.DefaultConfig()
	defaultConfig.Address = fmt.Sprintf("&s:%d", internal.ViperConf.ConsulConfig.Host, internal.ViperConf.ConsulConfig.Port)
	client, err := api.NewClient(defaultConfig)
	if err != nil {
		panic(err)
	}
	checkAddr := fmt.Sprintf("%s:%d", internal.ViperConf.AccountSrvConfig.Host, internal.ViperConf.AccountSrvConfig.Port)
	check := &api.AgentServiceCheck{
		GRPC:                           checkAddr,
		Timeout:                        "3s",
		Interval:                       "1s",
		DeregisterCriticalServiceAfter: "5s",
	}
	req := api.AgentServiceRegistration{
		Name:  internal.ViperConf.AccountSrvConfig.SrvName,
		ID:    internal.ViperConf.AccountSrvConfig.SrvName,
		Port:  internal.ViperConf.AccountSrvConfig.Port,
		Tags:  internal.ViperConf.AccountSrvConfig.Tags,
		Check: check,
	}
	err = client.Agent().ServiceRegister(&req)
	if err != nil {
		panic(err)
	}
	err = server.Serve(listen)
	if err != nil {
		panic(err)
	}
}
