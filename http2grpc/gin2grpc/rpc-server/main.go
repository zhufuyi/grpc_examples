package main

import (
	"fmt"
	"net"

	"github.com/zhufuyi/grpc_examples/http2grpc/gin2grpc/rpc-server/api/user/v1/pb"
	"github.com/zhufuyi/grpc_examples/http2grpc/gin2grpc/rpc-server/internal/service"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/zhufuyi/pkg/grpc/interceptor"
	"github.com/zhufuyi/pkg/logger"
	"google.golang.org/grpc"
)

const grpcAddr = ":8282"

func getServerOptions() []grpc.ServerOption {
	var options []grpc.ServerOption

	options = append(options, grpc_middleware.WithUnaryServerChain(
		interceptor.UnaryServerRecovery(),
		interceptor.UnaryServerCtxTags(),
		interceptor.UnaryServerLog(logger.Get()),
	))

	return options
}

func main() {
	list, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer(getServerOptions()...)
	pb.RegisterUserServiceServer(server, service.NewUserServiceServer())

	fmt.Println("start up grpc server ", grpcAddr)
	err = server.Serve(list)
	if err != nil {
		panic(err)
	}
}
