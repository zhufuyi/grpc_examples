package middleware

import (
	"context"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/zhufuyi/grpc_examples/pkg/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ---------------------------------- server interceptor ----------------------------------

var (
	// auth Scheme
	defaultAuthScheme = "Bearer"

	// 跳过认证方法集合
	skipMethods = map[string]struct{}{}
)

// GetAuthScheme 获取Scheme
func GetAuthScheme() string {
	return defaultAuthScheme
}

// AddSkipMethods 添加跳过认证方法，在服务初始化时设置
func AddSkipMethods(routers ...string) {
	for _, router := range routers {
		skipMethods[router] = struct{}{}
	}
}

// JwtVerify 从context获取token验证是否合法
func JwtVerify(ctx context.Context) (context.Context, error) {
	token, err := grpc_auth.AuthFromMD(ctx, defaultAuthScheme)
	if err != nil {
		return nil, err
	}

	cc, err := jwt.VerifyTokenCustom(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "%v", err)
	}

	newCtx := context.WithValue(ctx, "tokenInfo", cc) // 后面方法可以通过ctx.Value("tokenInfo").(*jwt.CustomClaims)

	return newCtx, nil
}

// UnaryServerJwtAuth jwt认证拦截器
func UnaryServerJwtAuth() grpc.UnaryServerInterceptor {
	return grpc_auth.UnaryServerInterceptor(JwtVerify)
}

// SkipAuthMethod 跳过认证，嵌入服务
type SkipAuthMethod struct{}

// AuthFuncOverride 重写认证方法
func (s *SkipAuthMethod) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	if _, ok := skipMethods[fullMethodName]; ok {
		return ctx, nil
	}

	return JwtVerify(ctx)
}
