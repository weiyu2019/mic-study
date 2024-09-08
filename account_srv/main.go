package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"mic-study/account_srv/biz"
	"mic-study/account_srv/internal"
	"mic-study/account_srv/proto/pb"
	"net"
)

func init() {
	internal.InitDB()
}

func main() {
	ip := flag.String("ip", "127.0.0.1", "输入IP")
	port := flag.Int("ip", 8080, "输入端口")
	flag.Parse()
	addr := fmt.Sprintf("%s:%d", *ip, *port)

	server := grpc.NewServer()
	pb.RegisterAccountServiceServer(server, &biz.AccountServer{})
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	err = server.Serve(listen)
	if err != nil {
		panic(err)
	}
}
