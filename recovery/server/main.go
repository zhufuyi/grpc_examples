package main

import (
	"context"
	"fmt"
	"net"
	"time"

	pb "github.com/zhufuyi/grpc_examples/recovery/proto/hellopb"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UnaryServerRecovery recovery unary interceptor
func UnaryServerRecovery() grpc.UnaryServerInterceptor {
	// https://pkg.go.dev/github.com/grpc-ecosystem/go-grpc-middleware/recovery
	customFunc := func(p interface{}) (err error) {
		return status.Errorf(codes.Internal, "triggered panic: %v", p)
	}
	opts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(customFunc),
	}

	return grpc_recovery.UnaryServerInterceptor(opts...)
}

type greeterServer struct {
	pb.UnimplementedGreeterServer
}

func (g *greeterServer) SayHello(ctx context.Context, r *pb.HelloRequest) (*pb.HelloReply, error) {
	var data []int
	fmt.Println(data[5]) // 下标越界，触发panic
	return &pb.HelloReply{Message: time.Now().Format("2006-01-02T15:04:05.000") + " hello " + r.Name}, nil
}

func getServerOptions() []grpc.ServerOption {
	var options []grpc.ServerOption

	recoveryOption := grpc_middleware.WithUnaryServerChain(
		UnaryServerRecovery(),
	)
	options = append(options, recoveryOption)

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
