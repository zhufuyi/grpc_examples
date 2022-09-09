package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/zhufuyi/grpc_examples/http2grpc/proto/pb"
	"github.com/zhufuyi/grpc_examples/include"
	"github.com/zhufuyi/grpc_examples/swagger"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

const (
	webAddr  = "127.0.0.1:9090"
	grpcAddr = "127.0.0.1:8080"
)

var (
	autoCount int64 = 0
	sm              = new(sync.Map)

	_ pb.UserServiceServer = (*UserServiceServer)(nil)
)

type UserServiceServer struct {
	pb.UnimplementedUserServiceServer
}

func NewUserServiceServer() pb.UserServiceServer {
	return &UserServiceServer{}
}

func (s UserServiceServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserReply, error) {
	id := atomic.AddInt64(&autoCount, 1)

	sm.LoadOrStore(id, &pb.User{
		Id:    id,
		Name:  req.Name,
		Email: req.Email,
	})

	return &pb.CreateUserReply{Id: id}, nil
}

func (s UserServiceServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserReply, error) {
	value, ok := sm.Load(req.Id)
	if !ok {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("id %v not found", req.Id))
	}
	user, ok := value.(*pb.User)
	if !ok {
		return nil, status.Error(codes.Internal, "type error")
	}

	return &pb.GetUserReply{User: user}, nil
}

// grpc服务
func grpcServer() {
	list, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer()
	pb.RegisterUserServiceServer(server, NewUserServiceServer())

	fmt.Println("start up grpc server ", grpcAddr)
	err = server.Serve(list)
	if err != nil {
		panic(err)
	}
}

// web服务
func webServer() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	gwMux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())} // 这里的option和rpc的client正常调用一致
	err := pb.RegisterUserServiceHandlerFromEndpoint(ctx, gwMux, grpcAddr, opts)
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	// 注册rpc服务的api接口路由
	mux.Handle("/", gwMux)

	// 注册swagger路由
	prefixPath := "/http2grpc/"
	router := swagger.Router(prefixPath, include.Path("../http2grpc/proto/pb/user.swagger.json"))
	mux.Handle(prefixPath, router) // 必须以/结尾的路径

	fmt.Println("start up web server ", webAddr)
	err = http.ListenAndServe(webAddr, mux)
	if err != nil {
		panic(err)
	}
}

func main() {
	go grpcServer()

	webServer()
}
