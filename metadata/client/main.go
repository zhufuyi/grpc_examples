package main

import (
	"context"
	"fmt"

	pb "github.com/zhufuyi/grpc_examples/metadata/proto/hellopb"

	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/zhufuyi/sponge/pkg/krand"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func sayHello(client pb.GreeterClient) error {
	// 创建 metadata 和 context
	tokenStr := krand.String(krand.R_All)
	ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs("token_str", tokenStr))

	// 使用带有元数据的上下文进行RPC调用。
	var header, trailer, md metadata.MD
	options := []grpc.CallOption{
		grpc.Header(&header),
		grpc.Trailer(&trailer),
	}

	resp, err := client.SayHello(ctx, &pb.HelloRequest{Name: "grpc"}, options...)
	if err != nil {
		return err
	}
	fmt.Println(resp.Message)

	printMetadata(header, "notfound")
	printMetadata(header, "foo")

	mergerMD := metadata.Join(md, header, trailer)
	printMetadata(mergerMD, "foo")

	return nil
}

func printMetadata(md metadata.MD, key string) {
	if ts, ok := md[key]; ok {
		fmt.Printf("metadata: %s=%v\n", key, ts)
	} else {
		fmt.Printf("not found '%s' in metadata\n", key)
	}
}

var contextRequestIDKey = "request_id"

func unaryClientRequestID() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		requestID := clientCtxRequestID(ctx)
		if requestID == "" {
			md, _ := metadata.FromOutgoingContext(ctx)
			requestID = krand.String(krand.R_All, 10)
			if md == nil {
				md = metadata.Pairs(contextRequestIDKey, requestID)
			} else {
				md[contextRequestIDKey] = []string{requestID}
			}
			ctx = metadata.NewOutgoingContext(ctx, md)
		}

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func clientCtxRequestID(ctx context.Context) string {
	return metautils.ExtractOutgoing(ctx).Get(contextRequestIDKey)
}

func getDialOptions() []grpc.DialOption {
	var options []grpc.DialOption

	// 使用不安全传输
	options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))

	// request id拦截器
	option := grpc.WithUnaryInterceptor(
		unaryClientRequestID(),
	)
	options = append(options, option)

	return options
}

func main() {
	conn, err := grpc.Dial("127.0.0.1:8282", getDialOptions()...)
	if err != nil {
		panic(err)
	}
	defer conn.Close() //nolint

	client := pb.NewGreeterClient(conn)

	err = sayHello(client)
	if err != nil {
		panic(err)
	}
}
