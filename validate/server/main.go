package main

import (
	"context"
	"fmt"
	"net"
	"time"

	pb "github.com/zhufuyi/grpc_examples/validate/proto/accountpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AccountServer struct {
	pb.UnimplementedAccountServer
}

func (g *AccountServer) Login(ctx context.Context, r *pb.LoginRequest) (*pb.LoginReply, error) {
	err := r.Validate()
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	token := fmt.Sprintf("%s_%s", time.Now().Format("2006-01-02T15:04:05.000000"), r.Email)
	return &pb.LoginReply{Token: token}, nil
}

func getServerOptions() []grpc.ServerOption {
	var options []grpc.ServerOption

	// 拦截器

	return options
}

func main() {
	addr := ":8080"
	fmt.Println("start rpc server", addr)

	// 监听TCP端口
	list, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	// 创建grpc server对象，拦截器可以在这里注入
	server := grpc.NewServer(getServerOptions()...)

	// grpc的server内部服务和路由
	pb.RegisterAccountServer(server, &AccountServer{})

	// 调用服务器执行阻塞等待客户端
	err = server.Serve(list)
	if err != nil {
		panic(err)
	}
}
