## Example of a GRPC implementation of a custom interceptor

Client-side interceptor implementation method

```go
func(ctx context.Context, method string, req, resp interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error
```

<br>

Server-side interceptor implementation method

```go
func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error)
```

