## http request forwarding rpc example

gRPC-Gateway is a plugin for the Google protocol buffer compiler protoc. It reads the protobuf service definition and generates a reverse proxy server that translates the RESTful HTTP API into gRPC. this server is generated based on the google.api.http annotation in your service definition.

![flowchart](grpc-gateway.png)

<br>

### Using

(1) Start the grpc server `go run main.go`

(2) test grpc client `go run main.go`

(3) test http request GET `http://127.0.0.1:9090/v1/getUser?id=1`
