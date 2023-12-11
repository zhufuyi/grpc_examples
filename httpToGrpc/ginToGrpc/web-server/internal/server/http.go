package server

import (
	"net/http"
	"time"

	"github.com/zhufuyi/grpc_examples/httpToGrpc/ginToGrpc/web-server/api/user/v1/pb"
	"github.com/zhufuyi/grpc_examples/httpToGrpc/ginToGrpc/web-server/internal/routers"
	"github.com/zhufuyi/grpc_examples/httpToGrpc/ginToGrpc/web-server/internal/service"
)

// NewHTTPServer creates a new web server
func NewHTTPServer(addr string) *http.Server {
	router := routers.NewRouter()
	server := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  time.Second * 3,
		WriteTimeout: time.Second * 3,
	}

	pb.RegisterUserServiceHTTPServer(router, service.NewUserServiceServer(NewUserClient()))

	return server
}
