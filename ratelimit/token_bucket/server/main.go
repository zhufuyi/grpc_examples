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
	return &pb.HelloReply{Message: time.Now().Format("2006-01-02T15:04:05.000") + " hello " + r.Name}, nil
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
