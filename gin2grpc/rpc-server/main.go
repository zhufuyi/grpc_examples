package main

import (
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/zhufuyi/pkg/grpc/middleware"
	"github.com/zhufuyi/pkg/logger"
	"net"

	"github.com/zhufuyi/grpc_examples/gin2grpc/rpc-server/api/user/v1/pb"
	"github.com/zhufuyi/grpc_examples/gin2grpc/rpc-server/internal/service"
	"google.golang.org/grpc"
)

const grpcAddr = ":9090"

func getServerOptions() []grpc.ServerOption {
	var options []grpc.ServerOption

	options = append(options, grpc_middleware.WithUnaryServerChain(
		middleware.UnaryServerRecovery(),
		middleware.UnaryServerCtxTags(),
		middleware.UnaryServerZapLogging(logger.Get()),
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
