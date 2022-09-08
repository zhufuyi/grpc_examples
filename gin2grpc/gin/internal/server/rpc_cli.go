package server

import (
	"github.com/zhufuyi/pkg/logger"

	userPB "github.com/zhufuyi/grpc_examples/gin2grpc/rpc-server/api/user/v1/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const userServerAddr = "127.0.0.1:9090"

func NewUserClient() userPB.UserServiceClient {
	conn, err := grpc.Dial(userServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	logger.Info("connect to 'user' rpc server successfully.", logger.String("addr", userServerAddr), logger.String("status", conn.GetState().String()))
	conn.GetState().String()
	return userPB.NewUserServiceClient(conn)
}
