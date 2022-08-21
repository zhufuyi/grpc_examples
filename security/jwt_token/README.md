## Example of jwt authentication for GRPC

### Using

(1) Start the grpc server `go run main.go`

(2) test grpc client `go run main.go`, output id value and token value

(3) test http request GET `https://127.0.0.1:9090/v1/getUser?id=<id> -h Authorization="Bearer <token>"`
