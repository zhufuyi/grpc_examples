## http request forwarding grpc example

gRPC-Gateway is a plugin for the Google protocol buffer compiler protoc. It reads the protobuf service definition and generates a reverse proxy server that translates the RESTful HTTP API into gRPC. this server is generated based on the google.api.http annotation in your service definition.

![flowchart](grpc-gateway.png)

<br>

### Using

Start the grpc server `go run main.go`

1. test grpc client `go run main.go`

2. test http request, open a browser and type in the address `http://127.0.0.1:8080/swagger/index.html`
