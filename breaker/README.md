## GRPC breaker example

### Running

```bash
# run grpc server
cd server && go run main.go

# run grpc client
cd client && go run main.go
```

<br>

### Breaker Code

```go
func getDialOptions() []grpc.DialOption {
	var options []grpc.DialOption

	// use insecure transfer
	options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))

	// circuit breaker
	option := grpc.WithUnaryInterceptor(
		grpc_middleware.ChainUnaryClient(
			interceptor.UnaryClientCircuitBreaker(),
		),
	)
	options = append(options, option)

	return options
}
```
