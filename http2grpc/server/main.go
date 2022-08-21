package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/zhufuyi/grpc_examples/http2grpc/proto/accountpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	webAddr  = "127.0.0.1:9090"
	grpcAddr = "127.0.0.1:8080"
)

type GreeterServer struct {
	pb.UnimplementedAccountServer
	m *sync.Map
}

func (g *GreeterServer) AddUser(ctx context.Context, user *pb.User) (*pb.ID, error) {
	g.m.LoadOrStore(user.Id, user)
	fmt.Printf("add user %v success\n", user)
	return &pb.ID{Id: user.Id}, nil
}

func (g *GreeterServer) GetUser(ctx context.Context, id *pb.ID) (*pb.User, error) {
	value, ok := g.m.Load(id.Id)
	if !ok {
		return nil, errors.New("not found user")
	}
	return value.(*pb.User), nil
}

// grpc服务
func grpcServer() {
	list, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer()
	pb.RegisterAccountServer(server, &GreeterServer{m: new(sync.Map)})

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
	err := pb.RegisterAccountHandlerFromEndpoint(ctx, gwMux, grpcAddr, opts)
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	// 注册rpc服务的api接口路由
	mux.Handle("/", gwMux)

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
