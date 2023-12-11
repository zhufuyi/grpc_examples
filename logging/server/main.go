package main

import (
	"context"
	"fmt"
	"net"

	"github.com/zhufuyi/sponge/pkg/grpc/interceptor"

	pb "github.com/zhufuyi/grpc_examples/logging/proto/hellopb"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var logger *zap.Logger

type greeterServer struct {
	pb.UnimplementedGreeterServer
}

func (g *greeterServer) SayHello(ctx context.Context, r *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "hello " + r.Name}, nil
}

func getServerOptions() []grpc.ServerOption {
	var options []grpc.ServerOption

	// 日志设置，默认打印客户端断开连接信息，示例 https://pkg.go.dev/github.com/grpc-ecosystem/go-grpc-middleware/logging/zap
	//middleware.AddLoggingFields(map[string]interface{}{"hello": "world"}) // 添加打印自定义字段
	//middleware.AddSkipLoggingMethods("/proto.Greeter/SayHello") // 跳过打印调用的方法
	options = append(options, grpc_middleware.WithUnaryServerChain(
		interceptor.UnaryServerLog(logger),
	))

	return options
}

func main() {
	logger, _ = zap.NewProduction()

	addr := ":8282"
	fmt.Println("grpc service is running", addr)

	// listening on TCP port
	list, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	// create a grpc server object where interceptors can be injected
	server := grpc.NewServer(getServerOptions()...)

	// register greeterServer to the server
	pb.RegisterGreeterServer(server, &greeterServer{})

	// start the server
	err = server.Serve(list)
	if err != nil {
		panic(err)
	}
}
