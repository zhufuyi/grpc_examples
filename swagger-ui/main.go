package main

import (
	"fmt"
	"net/http"

	"github.com/zhufuyi/grpc_examples/pkg/swagger"
	jsonfile "github.com/zhufuyi/grpc_examples/swagger-ui/proto"
)

var swaggerAddr = "127.0.0.1:8080"

func main() {
	mux := http.NewServeMux()

	// 注册swagger-ui和swagger json文件路由
	swagger.RegisterRoute(mux, jsonfile.Path("/hellopb"))

	fmt.Println("start up swagger server ", swaggerAddr)
	err := http.ListenAndServe(swaggerAddr, mux)
	if err != nil {
		panic(err)
	}
}

/*
启动之后，在浏览器打开 http://127.0.0.1:8080/swagger-ui
然后输入查看本地json文件 http://127.0.0.1:8080/swagger/hello.swagger.json
*/
