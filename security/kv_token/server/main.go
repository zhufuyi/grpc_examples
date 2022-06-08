package main

import (
	"context"
	"fmt"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/zhufuyi/grpc_examples/pkg/gtls"
	"github.com/zhufuyi/grpc_examples/pkg/gtls/certfile"
	pb "github.com/zhufuyi/grpc_examples/security/kv_token/proto/hellopb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type Geek struct {
	pb.UnimplementedGreeterServer
}

func (g *Geek) SayHello(ctx context.Context, request *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "hello " + request.Name}, nil
}

// CheckToken 检查token
func CheckToken(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	appID := metautils.ExtractIncoming(ctx).Get("app_id")
	appKey := metautils.ExtractIncoming(ctx).Get("app_key")
	if appID != "grpc" || appKey != "123456" {
		return nil, fmt.Errorf("%v token invalide: appid=%s, appkey=%s", codes.Unauthenticated, appID, appKey)
	}

	return handler(ctx, req)
}

func getServerOptions() []grpc.ServerOption {
	var options []grpc.ServerOption

	// tls加密
	credentials, err := gtls.GetServerTLSCredentials(certfile.Path("one-way/server.crt"), certfile.Path("one-way/server.key"))
	if err != nil {
		panic(err)
	}
	options = append(options, grpc.Creds(credentials))

	// token鉴权
	options = append(options, grpc.UnaryInterceptor(CheckToken))

	return options
}

func main() {
	addr := ":8080"
	fmt.Println("start rpc server", addr)

	list, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer(getServerOptions()...)

	pb.RegisterGreeterServer(s, &Geek{})

	err = s.Serve(list)
	if err != nil {
		panic(err)
	}
}
