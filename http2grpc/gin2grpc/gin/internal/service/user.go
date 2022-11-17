package service

import (
	"context"

	"github.com/zhufuyi/grpc_examples/http2grpc/gin2grpc/gin/api/user/v1/pb"
	userPb "github.com/zhufuyi/grpc_examples/http2grpc/gin2grpc/rpc-server/api/user/v1/pb"
)

var (
	_ pb.UserServiceHTTPServer = (*userServiceServer)(nil)
)

type userServiceServer struct {
	userRPCCli userPb.UserServiceClient
}

// NewUserServiceServer 实现接口
func NewUserServiceServer(cli userPb.UserServiceClient) pb.UserServiceHTTPServer {
	return &userServiceServer{userRPCCli: cli}
}

// CreateUser 创建用户
func (s userServiceServer) CreateUser(ctx context.Context, request *pb.CreateUserRequest) (*pb.CreateUserReply, error) {
	in := &userPb.CreateUserRequest{
		Name:  request.Name,
		Email: request.Email,
	}

	out, err := s.userRPCCli.CreateUser(ctx, in)
	if err != nil {
		return nil, err
	}

	return &pb.CreateUserReply{Id: out.Id}, nil
}

// GetUser 获取用户详情
func (s userServiceServer) GetUser(ctx context.Context, request *pb.GetUserRequest) (*pb.GetUserReply, error) {
	in := &userPb.GetUserRequest{
		Id: request.Id,
	}

	out, err := s.userRPCCli.GetUser(ctx, in)
	if err != nil {
		return nil, err
	}

	user := &pb.User{
		Id:    out.User.Id,
		Name:  out.User.Name,
		Email: out.User.Email,
	}

	return &pb.GetUserReply{User: user}, nil
}
