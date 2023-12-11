package service

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/zhufuyi/grpc_examples/httpToGrpc/ginToGrpc/grpc-server/api/user/v1/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	autoCount int64 = 0
	sm              = new(sync.Map)

	_ pb.UserServiceServer = (*userServiceServer)(nil)
)

type userServiceServer struct {
	pb.UnimplementedUserServiceServer
}

// NewUserServiceServer 实现接口
func NewUserServiceServer() pb.UserServiceServer {
	return &userServiceServer{}
}

// CreateUser 创建用户
func (s userServiceServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserReply, error) {
	id := atomic.AddInt64(&autoCount, 1)

	sm.LoadOrStore(id, &pb.User{
		Id:    id,
		Name:  req.Name,
		Email: req.Email,
	})

	return &pb.CreateUserReply{Id: id}, nil
}

// GetUser 获取用户详情
func (s userServiceServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserReply, error) {
	value, ok := sm.Load(req.Id)
	if !ok {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("id %v not found", req.Id))
	}
	user, ok := value.(*pb.User)
	if !ok {
		return nil, status.Error(codes.Internal, "type error")
	}

	return &pb.GetUserReply{User: user}, nil
}
