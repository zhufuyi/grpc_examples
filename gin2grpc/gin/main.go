package main

import "github.com/zhufuyi/grpc_examples/gin2grpc/gin/internal/server"

const addr = ":8080"

func main() {
	httpServer := server.NewHTTPServer(addr)

	if err := httpServer.ListenAndServe(); err != nil {
		panic(err)
	}
}
