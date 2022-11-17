package main

import (
	"context"
	"fmt"
	"net"
	"time"

	pb "github.com/zhufuyi/grpc_examples/ratelimit/token_bucket/proto/hellopb"

	"github.com/grpc-ecosystem/go-grpc-middleware/ratelimit"
	"github.com/reugn/equalizer"
	"google.golang.org/grpc"
)

type greeterServer struct {
	pb.UnimplementedGreeterServer
}

func (g *greeterServer) SayHello(ctx context.Context, r *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: time.Now().Format("2006-01-02T15:04:05.000000") + " hello " + r.Name}, nil
}

type myLimiter struct {
	TB *equalizer.TokenBucket // 令牌桶
}

func (m *myLimiter) Limit() bool {
	ok := m.TB.Ask()
	if ok {
		return false
	}

	fmt.Printf("rate limit triggered\n")
	return true
}

func getServerOptions() []grpc.ServerOption {
	var options []grpc.ServerOption

	limiter := &myLimiter{equalizer.NewTokenBucket(5, time.Millisecond*200)} // 限制5次/秒
	limitOption := grpc.ChainUnaryInterceptor(
		ratelimit.UnaryServerInterceptor(limiter),
	)
	options = append(options, limitOption)

	return options
}

func main() {
	addr := ":8282"
	fmt.Println("start rpc server", addr)

	// 监听TCP端口
	list, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	// 创建grpc server对象，拦截器可以在这里注入
	server := grpc.NewServer(getServerOptions()...)

	// grpc的server内部服务和路由
	pb.RegisterGreeterServer(server, &greeterServer{})

	// 调用服务器执行阻塞等待客户端
	err = server.Serve(list)
	if err != nil {
		panic(err)
	}
}
