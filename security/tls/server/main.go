package main

import (
	"context"
	"fmt"
	"net"

	pb "github.com/zhufuyi/grpc_examples/security/tls/proto/hellopb"

	"github.com/zhufuyi/sponge/pkg/grpc/gtls"
	"github.com/zhufuyi/sponge/pkg/grpc/gtls/certfile"
	"google.golang.org/grpc"
)

type greeterServer struct {
	pb.UnimplementedGreeterServer
}

func (g *greeterServer) SayHello(ctx context.Context, r *pb.HelloRequest) (*pb.HelloReply, error) {
	fmt.Println("\nSayHello receive req: " + r.Name)
	return &pb.HelloReply{Message: "hello " + r.Name}, nil
}

func main() {
	// 单向认证(服务端认证)
	//credentials, err := gtls.GetServerTLSCredentials(certfile.Path("one-way/server.crt"), certfile.Path("one-way/server.key"))
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

	// listening on TCP port
	addr := ":8282"
	fmt.Println("grpc service is running", addr)
	list, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	// create a grpc server object where interceptors can be injected
	server := grpc.NewServer(opts...)

	// register greeterServer to the server
	pb.RegisterGreeterServer(server, &greeterServer{})

	// start the server
	err = server.Serve(list)
	if err != nil {
		panic(err)
	}
}
