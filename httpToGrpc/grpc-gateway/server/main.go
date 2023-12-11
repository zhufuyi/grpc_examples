package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/zhufuyi/grpc_examples/httpToGrpc/grpc-gateway/proto/pb"
	"github.com/zhufuyi/grpc_examples/include"
	"github.com/zhufuyi/grpc_examples/swagger"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	webAddr  = "127.0.0.1:8080"
	grpcAddr = "127.0.0.1:8282"
)

var (
	autoCount int64 = 0
	sm              = new(sync.Map)

	_ pb.UserServiceServer = (*userServiceServer)(nil)
)

type userServiceServer struct {
	pb.UnimplementedUserServiceServer
}

// NewUserServiceServer  实现接口
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

// grpc服务
func grpcServer() {
	list, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer()
	pb.RegisterUserServiceServer(server, NewUserServiceServer())

	fmt.Println("start up grpc server ", grpcAddr)
	err = server.Serve(list)
	if err != nil {
		panic(err)
	}
}

func rpcGateway() (*runtime.ServeMux, error) {
	ctx := context.Background()

	mux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(runtime.DefaultHeaderMatcher),
		runtime.WithMarshalerOption(
			runtime.MIMEWildcard,
			&runtime.JSONPb{
				MarshalOptions: protojson.MarshalOptions{
					UseEnumNumbers:  false,
					EmitUnpopulated: true,
					UseProtoNames:   true,
				},
			},
		),
	)

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())} // 这里的option和rpc的client正常调用一致

	// 根据实际情况填写
	err := pb.RegisterUserServiceHandlerFromEndpoint(ctx, mux, grpcAddr, opts)
	if err != nil {
		return nil, err
	}

	return mux, nil
}

func webServer() {
	mux, err := rpcGateway()
	if err != nil {
		panic(err)
	}

	r := gin.Default()

	r.POST("/api/v1/*any", gin.WrapH(mux))
	r.DELETE("/api/v1/*any", gin.WrapH(mux))
	r.PUT("/api/v1/*any", gin.WrapH(mux))
	r.GET("/api/v1/*any", gin.WrapH(mux))
	swagger.Handler(r, include.Path("../httpToGrpc/grpc-gateway/proto/pb/user.swagger.json"))

	fmt.Println("http service is running", webAddr)
	server := &http.Server{
		Addr:           webAddr,
		Handler:        r,
		MaxHeaderBytes: 1 << 20,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(fmt.Errorf("listen server error: %v", err))
	}
}

func main() {
	go grpcServer()

	webServer()
}
