// Code generated by github.com/mohuishou/protoc-gen-go-gin. DO NOT EDIT.

package pb

import (
	context "context"
	errors "errors"
	gin "github.com/gin-gonic/gin"
	metadata "google.golang.org/grpc/metadata"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the mohuishou/protoc-gen-go-gin package it is being compiled against.
// context.metadata.
//gin.errors.

type UserServiceHTTPServer interface {
	CreateUser(context.Context, *CreateUserRequest) (*CreateUserReply, error)

	GetUser(context.Context, *GetUserRequest) (*GetUserReply, error)
}

func RegisterUserServiceHTTPServer(r gin.IRouter, srv UserServiceHTTPServer) {
	s := UserService{
		server: srv,
		router: r,
		resp:   defaultUserServiceResp{},
	}
	s.RegisterService()
}

type UserService struct {
	server UserServiceHTTPServer
	router gin.IRouter
	resp   interface {
		Error(ctx *gin.Context, err error)
		ParamsError(ctx *gin.Context, err error)
		Success(ctx *gin.Context, data interface{})
	}
}

// Resp 返回值
type defaultUserServiceResp struct{}

func (resp defaultUserServiceResp) response(ctx *gin.Context, status, code int, msg string, data interface{}) {
	ctx.JSON(status, map[string]interface{}{
		"code": code,
		"msg":  msg,
		"data": data,
	})
}

// Error 返回错误信息
func (resp defaultUserServiceResp) Error(ctx *gin.Context, err error) {
	code := -1
	status := 500
	msg := "未知错误"

	if err == nil {
		msg += ", err is nil"
		resp.response(ctx, status, code, msg, nil)
		return
	}

	type iCode interface {
		HTTPCode() int
		Message() string
		Code() int
	}

	var c iCode
	if errors.As(err, &c) {
		status = c.HTTPCode()
		code = c.Code()
		msg = c.Message()
	}

	_ = ctx.Error(err)

	resp.response(ctx, status, code, msg, nil)
}

// ParamsError 参数错误
func (resp defaultUserServiceResp) ParamsError(ctx *gin.Context, err error) {
	_ = ctx.Error(err)
	resp.response(ctx, 400, 400, "参数错误", nil)
}

// Success 返回成功信息
func (resp defaultUserServiceResp) Success(ctx *gin.Context, data interface{}) {
	resp.response(ctx, 200, 0, "成功", data)
}

func (s *UserService) CreateUser_0(ctx *gin.Context) {
	var in CreateUserRequest

	if err := ctx.ShouldBindJSON(&in); err != nil {
		s.resp.ParamsError(ctx, err)
		return
	}

	md := metadata.New(nil)
	for k, v := range ctx.Request.Header {
		md.Set(k, v...)
	}
	newCtx := metadata.NewIncomingContext(ctx, md)
	out, err := s.server.(UserServiceHTTPServer).CreateUser(newCtx, &in)
	if err != nil {
		s.resp.Error(ctx, err)
		return
	}

	s.resp.Success(ctx, out)
}

func (s *UserService) GetUser_0(ctx *gin.Context) {
	var in GetUserRequest

	if err := ctx.ShouldBindUri(&in); err != nil {
		s.resp.ParamsError(ctx, err)
		return
	}

	if err := ctx.ShouldBindQuery(&in); err != nil {
		s.resp.ParamsError(ctx, err)
		return
	}

	md := metadata.New(nil)
	for k, v := range ctx.Request.Header {
		md.Set(k, v...)
	}
	newCtx := metadata.NewIncomingContext(ctx, md)
	out, err := s.server.(UserServiceHTTPServer).GetUser(newCtx, &in)
	if err != nil {
		s.resp.Error(ctx, err)
		return
	}

	s.resp.Success(ctx, out)
}

func (s *UserService) RegisterService() {

	s.router.Handle("POST", "/v1/user", s.CreateUser_0)

	s.router.Handle("GET", "/v1/user/:id", s.GetUser_0)

}
