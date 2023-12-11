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

type accountServer struct {
	pb.UnimplementedAccountServer
}

// Login 登录
func (g *accountServer) Login(ctx context.Context, r *pb.LoginRequest) (*pb.LoginReply, error) {
	err := r.Validate()
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	token := fmt.Sprintf("%s_%s", time.Now().Format("2006-01-02T15:04:05.000"), r.Email)
	return &pb.LoginReply{Token: token}, nil
}

func getServerOptions() []grpc.ServerOption {
	var options []grpc.ServerOption

	// 拦截器

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
	pb.RegisterAccountServer(server, &accountServer{})

	// start the server
	err = server.Serve(list)
	if err != nil {
		panic(err)
	}
}
