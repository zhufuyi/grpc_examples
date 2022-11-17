package server

import (
	"github.com/zhufuyi/grpc_examples/http2grpc/gin2grpc/rpc-server/api/user/v1/pb"

	"github.com/zhufuyi/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const userServerAddr = "127.0.0.1:8282"

// NewUserClient user客户端
func NewUserClient() pb.UserServiceClient {
	conn, err := grpc.Dial(userServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	logger.Info("connect to 'user' rpc server successfully.", logger.String("addr", userServerAddr), logger.String("status", conn.GetState().String()))

	return pb.NewUserServiceClient(conn)
}
