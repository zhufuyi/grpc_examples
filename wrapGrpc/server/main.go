package main

import (
	"context"
	"fmt"

	pb "github.com/zhufuyi/grpc_examples/wrapGrpc/proto/hellopb"

	"github.com/zhufuyi/sponge/pkg/grpc/gtls"
	"github.com/zhufuyi/sponge/pkg/grpc/gtls/certfile"
	"github.com/zhufuyi/sponge/pkg/grpc/interceptor"
	"github.com/zhufuyi/sponge/pkg/grpc/server"
	"github.com/zhufuyi/sponge/pkg/logger"
	"github.com/zhufuyi/sponge/pkg/servicerd/registry/etcd"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type greeterServer struct {
	pb.UnimplementedGreeterServer
}

func (g *greeterServer) SayHello(ctx context.Context, r *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "hello " + r.Name}, nil
}

var newUnaryInterceptors = []grpc.UnaryServerInterceptor{
	interceptor.UnaryServerRecovery(),
	interceptor.UnaryServerLog(logger.Get()),
}

var newServiceRegister = func() {
	instanceEndpoint := fmt.Sprintf("grpc://%s:%d", "127.0.0.1", 8282)
	iRegistry, instance, err := etcd.NewRegistry(
		[]string{"192.168.3.37:2379"},
		"test-id",
		"hello-demo",
		[]string{instanceEndpoint},
	)
	if err != nil {
		panic(err)
	}
	err = iRegistry.Register(context.Background(), instance)
	if err != nil {
		panic(err)
	}
}

func newCredential() credentials.TransportCredentials {
	var (
		credential credentials.TransportCredentials
		err        error
	)

	// server-side authentication only
	credential, err = gtls.GetServerTLSCredentials(
		certfile.Path("one-way/server.crt"),
		certfile.Path("one-way/server.key"),
	)

	// two-way authentication between client and server
	//credential, err = gtls.GetServerTLSCredentialsByCA(
	//	certfile.Path("two-way/ca.pem"),
	//	certfile.Path("two-way/server/server.pem"),
	//	certfile.Path("two-way/server/server.key"),
	//)
	if err != nil {
		panic(err)
	}
	return credential
}

func main() {
	port := 8282
	fn := func(s *grpc.Server) {
		pb.RegisterGreeterServer(s, &greeterServer{})
	}

	server.Run(port, []server.RegisterFn{fn},
		server.WithSecure(newCredential()),
		server.WithUnaryInterceptor(newUnaryInterceptors...),
		// WithStreamInterceptor(streamInterceptors...),
		//server.WithServiceRegister(newServiceRegister),
	)
	logger.Info("grpc server is running", logger.Int("port", port))
	select {}
}
