package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"sync"
	
	"grpc_examples/pkg/middleware"
	"grpc_examples/pkg/gtls"
	"grpc_examples/pkg/gtls/certfile"
	"grpc_examples/pkg/jwt"
	"grpc_examples/pkg/swagger"
	jsonfile "grpc_examples/security/jwt_token/proto"
	pb "grpc_examples/security/jwt_token/proto/accountpb"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/zhufuyi/pkg/snowFlake"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

const (
	//authScheme = "bearer"
	webAddr  = "127.0.0.1:8080"
	grpcAddr = "127.0.0.1:9090"
)

var isUseTLS bool // 是否开启TLS加密

type Account struct {
	pb.UnimplementedAccountServer
	middleware.SkipAuthMethod // 实现AuthFuncOverride方法，跳过认证路由
	m                         *sync.Map
}

type userInfo struct {
	ID       int64
	Name     string
	Password string
	Email    string
	Token    string
}

func (a *Account) getUserFromID(id int64) *userInfo {
	if userInfoVal, ok := a.m.Load(id); ok {
		return userInfoVal.(*userInfo)
	}
	return nil
}

func (a *Account) getUserFromName(name string) *userInfo {
	if id, ok := a.m.Load(name); ok {
		if userInfoVal, ok := a.m.Load(id); ok {
			return userInfoVal.(*userInfo)
		}
	}
	return nil
}

func (a *Account) saveUser(u *userInfo) {
	a.m.LoadOrStore(u.ID, u)
	a.m.LoadOrStore(u.Name, u.ID)
}

func (a *Account) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterReply, error) {
	uinfo := a.getUserFromName(req.Name)
	if uinfo != nil {
		return &pb.RegisterReply{Id: uinfo.ID, Token: uinfo.Token}, nil
	}

	// 生成id
	id := snowFlake.NewID()

	// 生成token
	token, err := jwt.GenerateTokenWithCustom(strconv.FormatInt(id, 10))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "generate token error %v", err)
	}
	token = middleware.GetAuthScheme() + " " + token
	//fmt.Printf("create user uid:%s %v token:%s \n", uid, req, token)

	a.saveUser(&userInfo{id, req.Name, req.Password, req.Name + "@126.com", token})

	fmt.Printf("save data: %+v\n", a.getUserFromID(id))
	return &pb.RegisterReply{Id: id, Token: token}, nil
}

// 需要鉴权
func (a *Account) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserReply, error) {
	//uid := metautils.ExtractIncoming(ctx).Get("uid")            // 这是从header获取
	tokenInfo, ok := ctx.Value("tokenInfo").(*jwt.CustomClaims) // 从拦截器设置值读取
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "not found field 'tokenInfo'")
	}

	uid := strconv.FormatInt(req.Id, 10)
	if tokenInfo.Uid != uid {
		return nil, status.Errorf(codes.InvalidArgument, "the token uid and the parameter uid do not match")
	}

	// 从数据库或缓存中获取
	uinfo := a.getUserFromID(req.Id)
	if uinfo == nil {
		return nil, status.Errorf(codes.NotFound, "not found user")
	}

	return &pb.GetUserReply{
		Id:    req.Id,
		Name:  uinfo.Name,
		Email: uinfo.Email,
	}, nil
}

func getServerOptions(isUseTLS bool) []grpc.ServerOption {
	var options []grpc.ServerOption

	// tls加密
	if isUseTLS {
		credentials, err := gtls.GetServerTLSCredentials(certfile.Path("one-way/server.crt"), certfile.Path("one-way/server.key"))
		if err != nil {
			panic(err)
		}
		options = append(options, grpc.Creds(credentials))
	}

	// token鉴权
	options = append(options, grpc.UnaryInterceptor(middleware.JwtAuth()))

	return options
}

func grpcServer() {
	fmt.Println("start rpc server", grpcAddr)
	middleware.AddSkipMethods("/proto.Account/Register") // 添加忽略token验证的方法，从pb文件的fullMethodName

	listen, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer(getServerOptions(isUseTLS)...)

	pb.RegisterAccountServer(server, &Account{m: new(sync.Map)})

	err = server.Serve(listen)
	if err != nil {
		panic(err)
	}
}

func getDialOptions(isUseTLS bool) []grpc.DialOption {
	var options []grpc.DialOption

	if !isUseTLS {
		// 不使用tls
		options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		// 使用tls加密，需要服务端同时使用tls
		credentials, err := gtls.GetClientTLSCredentials("localhost", certfile.Path("one-way/server.crt"))
		if err != nil {
			panic(err)
		}
		options = append(options, grpc.WithTransportCredentials(credentials))
	}

	return options
}

// web服务
func webServer() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	gwMux := runtime.NewServeMux()
	options := getDialOptions(isUseTLS)
	err := pb.RegisterAccountHandlerFromEndpoint(ctx, gwMux, grpcAddr, options)
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	// 注册rpc服务的api接口路由
	mux.Handle("/", gwMux)
	// 注册swagger-ui和swagger json文件路由
	swagger.RegisterRoute(mux, jsonfile.Path("/accountpb"))

	fmt.Println("start web server ", webAddr)
	if !isUseTLS {
		err = http.ListenAndServe(webAddr, mux)
	} else {
		err = http.ListenAndServeTLS(webAddr, certfile.Path("one-way/server.crt"), certfile.Path("one-way/server.key"), mux) // 浏览器不信任证书，报错 http: TLS handshake error from 127.0.0.1:2358: remote error: tls: unknown certificate
	}
	if err != nil {
		panic(err)
	}
}

func main() {
	// 设置是否开启TLS
	isUseTLS = true

	// 设置用户id生成器
	snowFlake.InitSnowFlake(1)

	go grpcServer()

	webServer()
}

/*
使用：
(1) 启动服务
go run main.go

(2) 因为isUseTLS = true，设置使用了TLS加密，访问都需要https
在浏览器访问swagger UI https://127.0.0.1:8080/swagger-ui/
输入 https://127.0.0.1:8080/swagger/account.swagger.json

(3) 测试http到grpc接口
*/
