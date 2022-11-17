package main

import (
	"context"
	"fmt"
	"net"
	"time"

	pb "github.com/zhufuyi/grpc_examples/interceptor/proto/hellopb"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type greeterServer struct {
	pb.UnimplementedGreeterServer
}

func (g *greeterServer) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	if req.Name == "" {
		return nil, status.Errorf(codes.InvalidArgument, "name is empty")
	}
	return &pb.HelloReply{Message: "hello " + req.Name}, nil
}

// AccessLog 请求日志
func AccessLog(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	beginTime := time.Now().Local().Unix()
	fmt.Printf("grpc request, method=%s, beginTime=%d, request=%v\n", info.FullMethod, beginTime, req)

	resp, err := handler(ctx, req)

	endTime := time.Now().Local().Unix()
	fmt.Printf("grpc response, method=%s, beginTime=%d, end_time=%d, response=%v", info.FullMethod, beginTime, endTime, resp)
	return resp, err
}

// ErrorLog 错误日志
func ErrorLog(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	resp, err := handler(ctx, req)
	if err != nil {
		fmt.Printf("grpc error, method=%s, err=%s", info.FullMethod, err.Error())
	}

	return resp, err
}

func getServerOptions() []grpc.ServerOption {
	var options []grpc.ServerOption

	// 拦截器
	option := grpc.UnaryInterceptor(
		grpc_middleware.ChainUnaryServer(
			AccessLog,
			ErrorLog,
		),
	)
	options = append(options, option)

	return options
}

func main() {
	addr := ":8282"
	fmt.Println("start rpc server", addr)

	list, err := net.Listen("tcp", addr)

	server := grpc.NewServer(getServerOptions()...)
	defer server.Stop()

	pb.RegisterGreeterServer(server, &greeterServer{})
	if err != nil {
		panic(err)
	}

	err = server.Serve(list)
	if err != nil {
		panic(err)
	}
}
