## Example of a token bucket rate limit for GRPC

Token bucket algorithm: tokens are added to the bucket at a fixed rate, and no more tokens are added when the bucket is full. When the service receives a request, it attempts to retrieve a token from the barrel, and if it gets a token, it continues to execute the subsequent business logic; if it does not get a token, it directly returns a rhetorical frequency overrun error code.

```bash
# run grpc server
cd server && go run main.go

# run grpc client
cd client && go run main.go
```
