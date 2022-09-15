package main

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	pb "github.com/zhufuyi/grpc_examples/retry/proto/hellopb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const timestampFormat = "2006-01-02T15:04:05.000000"

type greeterServer struct {
	pb.UnimplementedGreeterServer
	failingBuilder
}

func newGreeterServer(reqPassSize uint) *greeterServer {
	gs := &greeterServer{}
	gs.reqPassSize = reqPassSize
	return gs
}

func (g *greeterServer) SayHello(ctx context.Context, r *pb.HelloRequest) (*pb.HelloReply, error) {
	err := g.maybeFailRequest()
	if err != nil {
		fmt.Printf("[%s] request failed, couter=%d\n", time.Now().Format(timestampFormat), g.reqCounter)
		return nil, err
	}
	fmt.Printf("[%s] successfully, counter=%d\n\n", time.Now().Format(timestampFormat), g.reqCounter)

	return &pb.HelloReply{Message: "hello " + r.Name}, nil
}

type failingBuilder struct {
	mu          sync.Mutex
	reqCounter  uint // 请求次数
	reqPassSize uint // 当请求次数达到值时，正常放行
}

func (s *failingBuilder) maybeFailRequest() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.reqCounter++
	if (s.reqPassSize > 0) && (s.reqCounter%s.reqPassSize == 0) {
		return nil
	}

	return status.Errorf(codes.Unavailable, "maybeFailRequest: failing it")
}

func main() {
	addr := ":8080"
	fmt.Println("start rpc server", addr)

	list, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer()

	pb.RegisterGreeterServer(server, newGreeterServer(3))

	err = server.Serve(list)
	if err != nil {
		panic(err)
	}
}
