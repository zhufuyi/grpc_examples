package main

import (
	"context"
	"fmt"
	"net"

	"github.com/zhufuyi/grpc_examples/pkg/gtls"
	"github.com/zhufuyi/grpc_examples/pkg/gtls/certfile"
	pb "github.com/zhufuyi/grpc_examples/security/tls/proto/hellopb"
	"google.golang.org/grpc"
)

type GreeterServer struct {
	pb.UnimplementedGreeterServer
}

func (g *GreeterServer) SayHello(ctx context.Context, r *pb.HelloRequest) (*pb.HelloReply, error) {
	fmt.Println("\nSayHello receive req: " + r.Name)
	return &pb.HelloReply{Message: "hello " + r.Name}, nil
}

func main() {
	// 单向认证(服务端认证)
	//credentials, err := gtls.GetServerTLSCredentials(certfile.Path("/one-way/server.crt"), certfile.Path("/one-way/server.key"))
	// 双向认证
	credentials, err := gtls.GetServerTLSCredentialsByCA(
		certfile.Path("two-way/ca.pem"),
		certfile.Path("two-way/server/server.pem"),
		certfile.Path("two-way/server/server.key"),
	)
	if err != nil {
		panic(err)
	}

	// 拦截器
	opts := []grpc.ServerOption{
		grpc.Creds(credentials),
	}

	// 监听TCP端口
	addr := ":8080"
	fmt.Println("start rpc server", addr)
	list, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	// 创建grpc server对象，拦截器可以在这里注入
	server := grpc.NewServer(opts...)

	// grpc的server内部服务和路由
	pb.RegisterGreeterServer(server, &GreeterServer{})

	// 调用服务器执行阻塞等待客户端
	err = server.Serve(list)
	if err != nil {
		panic(err)
	}
}
