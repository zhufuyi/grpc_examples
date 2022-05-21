package main

import (
	"context"
	"fmt"
	"net"

	pb "grpc_examples/logging/proto/hellopb"
	"grpc_examples/pkg/middleware"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var logger *zap.Logger

type GreeterServer struct {
	pb.UnimplementedGreeterServer
}

func (g *GreeterServer) SayHello(ctx context.Context, r *pb.HelloRequest) (*pb.HelloReply, error) {
	fmt.Println("\nSayHello receive req: " + r.Name)

	tag := grpc_ctxtags.NewTags().Set("fullMethodName", "/proto.Greeter/SayHello")
	grpc_ctxtags.SetInContext(ctx, tag)

	return &pb.HelloReply{Message: "hello " + r.Name}, nil
}

func getServerOptions() []grpc.ServerOption {
	var options []grpc.ServerOption

	// 日志设置，默认打印客户端断开连接信息，示例 https://pkg.go.dev/github.com/grpc-ecosystem/go-grpc-middleware/logging/zap
	//middleware.AddLoggingFields(map[string]interface{}{"hello": "world"}) // 添加打印自定义字段
	//middleware.AddSkipLoggingMethods("/proto.Greeter/SayHello") // 跳过打印调用的方法
	options = append(options, grpc_middleware.WithUnaryServerChain(
		middleware.CtxFieldExtractor(),
		middleware.ZapLogging(logger),
	))

	return options
}

func main() {
	logger, _ = zap.NewProduction()

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
	pb.RegisterGreeterServer(server, &GreeterServer{})

	// 调用服务器执行阻塞等待客户端
	err = server.Serve(list)
	if err != nil {
		panic(err)
	}
}
