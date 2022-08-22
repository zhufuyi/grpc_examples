package main

import (
	"fmt"
	"net/http"

	"github.com/zhufuyi/grpc_examples/swagger"
	"github.com/zhufuyi/grpc_examples/swagger/example/proto"
)

var swaggerAddr = ":9090"

func main() {
	router := swagger.Router("/", proto.Path("hellopb/hello.swagger.json"))
	server := &http.Server{
		Addr:    swaggerAddr,
		Handler: router,
	}

	fmt.Println("start up web server ", swaggerAddr)
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
