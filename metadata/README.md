## Example of GRPC setting and reading metadata

Client request: metadata is passed from the client's context to the server's context, and the context's key-value is transferred to the http2 request header key-value.

Server reply: metadata is passed from the server-side context to the client-side context, again via the http2 request header.

When client-->server or server-->client passes the context, the key-value of the context is read。


```go
key = "authorization"

// Use metadata packages
md, ok := metadata.FromIncomingContext(ctx)
if !ok {
    return status.Errorf(codes.DataLoss, "failed to get metadata")
}
if authorization, ok := md[key]; ok {
    fmt.Printf("metadata: %s=%v\n", key, authorization)
} else {
    fmt.Printf("not found '%s' in metadata\n", key)
}

// Use of encapsulated third party packages
authorization := metautils.ExtractIncoming(ctx).Get(key)
fmt.Println(authorization)
```

### Running

```bash
# run grpc server
cd server && go run main.go

# run grpc client
cd client && go run main.go
```

