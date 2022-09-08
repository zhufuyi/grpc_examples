## gin2grpc

The handlers for gin are automatically generated via the proto file. the request http is followed by a call to the rpc server via the client.

<br>

### install plugin

> go install github.com/mohuishou/protoc-gen-go-gin@v0.1.0

<br>

### Using

(1) start the rpc-server `go run main.go`

(2) start gin server `go run main.go`

(3) open a browser and type in the address `http://127.0.0.1:8080/swagger/index.html`, test the request.
