package main

import (
	"fmt"

	"github.com/zhufuyi/grpc_examples/httpToGrpc/ginToGrpc/web-server/internal/server"
)

const addr = ":8080"

func main() {
	httpServer := server.NewHTTPServer(addr)

	fmt.Println("http service listening on " + addr)
	if err := httpServer.ListenAndServe(); err != nil {
		panic(err)
	}
}
