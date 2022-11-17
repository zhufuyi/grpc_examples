package main

import (
	"github.com/zhufuyi/pkg/logger"

	"github.com/zhufuyi/grpc_examples/http2grpc/gin2grpc/gin/internal/server"
)

const addr = ":8080"

func main() {
	httpServer := server.NewHTTPServer(addr)

	logger.Info("http service listening on " + addr)
	if err := httpServer.ListenAndServe(); err != nil {
		panic(err)
	}
}
