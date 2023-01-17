package main

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"time"

	pb "github.com/zhufuyi/grpc_examples/loadbalance/client_loadbalance/proto/hellopb"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/zhufuyi/pkg/krand"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type greeterServer struct {
	pb.UnimplementedGreeterServer
}

func (g *greeterServer) SayHello(ctx context.Context, r *pb.HelloRequest) (*pb.HelloReply, error) {
	// return之后创建新的metadata
	defer func() {
		trailer := metadata.Pairs("foo", "bar2")
		err := grpc.SetTrailer(ctx, trailer)
		if err != nil {
			fmt.Println(err)
		}
	}()

	// 从client读取metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.DataLoss, "SayHello: failed to get metadata")
	}
	printMetadata(md, "token_str")
	printMetadata(md, contextRequestIDKey)

	// 创建新的metadata
	header := metadata.New(map[string]string{
		"foo": "bar",
	})
	_ = grpc.SendHeader(ctx, header)

	// 随机延时0~1000微秒
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Microsecond)

	return &pb.HelloReply{Message: "hello " + r.Name}, nil
}

func printMetadata(md metadata.MD, key string) {
	if ts, ok := md[key]; ok {
		fmt.Printf("metadata: %s=%v\n", key, ts)
	} else {
		fmt.Printf("not found '%s' in metadata\n", key)
	}
}

var contextRequestIDKey = "request_id"

func unaryServerRequestID() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		requestID := serverCtxRequestID(ctx)
		if requestID == "" {
			md, _ := metadata.FromIncomingContext(ctx)
			requestID = krand.String(krand.R_All, 10)
			if md == nil {
				md = metadata.Pairs(contextRequestIDKey, requestID)
			} else {
				md[contextRequestIDKey] = []string{requestID}
			}
			ctx = metadata.NewOutgoingContext(ctx, md)
		}

		return handler(ctx, req)
	}
}

func serverCtxRequestID(ctx context.Context) string {
	return metautils.ExtractIncoming(ctx).Get(contextRequestIDKey)
}

func getServerOptions() []grpc.ServerOption {
	var options []grpc.ServerOption

	options = append(options, grpc_middleware.WithUnaryServerChain(
		unaryServerRequestID(),
	))

	return options
}

func main() {
	addr := ":8282"
	fmt.Println("start rpc server", addr)

	list, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer(getServerOptions()...)

	pb.RegisterGreeterServer(server, &greeterServer{})

	err = server.Serve(list)
	if err != nil {
		panic(err)
	}
}
