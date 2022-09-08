package service

import (
	"context"
	"github.com/zhufuyi/grpc_examples/gin2grpc/gin/api/user/v1/pb"
	userPB "github.com/zhufuyi/grpc_examples/gin2grpc/rpc-server/api/user/v1/pb"
)

var (
	_ pb.UserServiceHTTPServer = (*UserServiceServer)(nil)
)

type UserServiceServer struct {
	userRPCCli userPB.UserServiceClient
}

func NewUserServiceServer(cli userPB.UserServiceClient) pb.UserServiceHTTPServer {
	return &UserServiceServer{userRPCCli: cli}
}

func (s UserServiceServer) CreateUser(ctx context.Context, request *pb.CreateUserRequest) (*pb.CreateUserReply, error) {
	in := &userPB.CreateUserRequest{
		Name:  request.Name,
		Email: request.Email,
	}

	out, err := s.userRPCCli.CreateUser(ctx, in)
	if err != nil {
		return nil, err
	}

	return &pb.CreateUserReply{Id: out.Id}, nil
}

func (s UserServiceServer) GetUser(ctx context.Context, request *pb.GetUserRequest) (*pb.GetUserReply, error) {
	in := &userPB.GetUserRequest{
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
