package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"sync"

	"github.com/zhufuyi/grpc_examples/include"
	pb "github.com/zhufuyi/grpc_examples/security/jwt_token/proto/accountpb"
	"github.com/zhufuyi/grpc_examples/swagger"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/zhufuyi/sponge/pkg/grpc/gtls"
	"github.com/zhufuyi/sponge/pkg/grpc/gtls/certfile"
	"github.com/zhufuyi/sponge/pkg/grpc/interceptor"
	"github.com/zhufuyi/sponge/pkg/jwt"
	"github.com/zhufuyi/sponge/pkg/krand"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

const (
	webAddr  = ":8080"
	grpcAddr = ":8282"
)

var isUseTLS bool // 是否开启TLS加密

type accountServer struct {
	pb.UnimplementedAccountServer
	m *sync.Map
}

type userInfo struct {
	ID            int64
	Name          string
	Password      string
	Email         string
	Authorization string
}

func (a *accountServer) getUserFromID(id int64) *userInfo {
	if userInfoVal, ok := a.m.Load(id); ok {
		return userInfoVal.(*userInfo)
	}
	return nil
}

func (a *accountServer) getUserFromName(name string) *userInfo {
	if id, ok := a.m.Load(name); ok {
		if userInfoVal, ok := a.m.Load(id); ok {
			return userInfoVal.(*userInfo)
		}
	}
	return nil
}

func (a *accountServer) saveUser(u *userInfo) {
	a.m.LoadOrStore(u.ID, u)
	a.m.LoadOrStore(u.Name, u.ID)
}

// Register 注册
func (a *accountServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterReply, error) {
	uinfo := a.getUserFromName(req.Name)
	if uinfo != nil {
		return &pb.RegisterReply{Id: uinfo.ID, Authorization: uinfo.Authorization}, nil
	}

	// 生成id
	id := int64(krand.Int())

	// 生成token
	token, err := jwt.GenerateToken(strconv.FormatInt(id, 10))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "generate token error %v", err)
	}
	authorization := interceptor.GetAuthorization(token)
	//fmt.Printf("create user uid:%s %v token:%s \n", uid, req, token)

	a.saveUser(&userInfo{id, req.Name, req.Password, req.Name + "@bar.com", token})

	fmt.Printf("save data: %+v\n", a.getUserFromID(id))
	return &pb.RegisterReply{Id: id, Authorization: authorization}, nil
}

// GetUser 获取用详情，需要鉴权
func (a *accountServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserReply, error) {
	//uid := metautils.ExtractIncoming(ctx).Get("uid")            // 这是从header获取
	tokenInfo, ok := ctx.Value("tokenInfo").(*jwt.Claims) // 从拦截器ctx读取
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "not found field 'tokenInfo'")
	}

	uid := strconv.FormatInt(req.Id, 10)
	if tokenInfo.UID != uid {
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
	options = append(options, grpc.UnaryInterceptor(interceptor.UnaryServerJwtAuth(
		interceptor.WithAuthIgnoreMethods("/proto.Account/Register"), // 添加忽略token验证的方法
	)))

	return options
}

func grpcServer() {
	fmt.Println("grpc service is running", grpcAddr)

	listen, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer(getServerOptions(isUseTLS)...)

	pb.RegisterAccountServer(server, &accountServer{m: new(sync.Map)})

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

	// 注册swagger路由
	prefixPath := "/token/"
	router := swagger.Router(prefixPath, include.Path("../security/jwt_token/proto/accountpb/account.swagger.json"))
	mux.Handle(prefixPath, router) // 必须以/结尾的路径

	fmt.Println("http service is running", webAddr)
	if !isUseTLS {
		err = http.ListenAndServe(webAddr, mux)
	} else {
		// 浏览器不信任证书，报错 http: TLS handshake error from 127.0.0.1:2358: remote error: tls: unknown certificate
		err = http.ListenAndServeTLS(webAddr, certfile.Path("one-way/server.crt"), certfile.Path("one-way/server.key"), mux)
	}
	if err != nil {
		panic(err)
	}
}

func main() {
	// 设置是否开启TLS
	isUseTLS = true

	// 初始jwt
	jwt.Init()

	go grpcServer()

	webServer()
}
