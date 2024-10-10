package main

import (
	"context"
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"mic-study/account_srv/proto/pb"
	"mic-study/internal"
)

func main() {
	addr := fmt.Sprintf("%s:%d", internal.AppConf.ConsulConfig.Host, internal.AppConf.ConsulConfig.Port)
	dialAddr := fmt.Sprintf("consul://%s/%s?wait=14", addr, internal.AppConf.AccountSrvConfig.SrvName)
	conn, err := grpc.Dial(dialAddr, grpc.WithInsecure(), grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`))
	if err != nil {
		zap.S().Fatal(err)
	}
	defer conn.Close()
	client := pb.NewAccountServiceClient(conn)
	res, err := client.GetAccountList(context.Background(), &pb.PagingRequest{
		PageNo:   1,
		PageSize: 3,
	})
	if err != nil {
		zap.S().Fatal(err)
	}
	for idx, item := range res.AccountList {
		fmt.Println(fmt.Sprintf("%d---%v", idx, item))
	}
}
