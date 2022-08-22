## Example of jwt authentication for GRPC

### Using

(1) Start the grpc server `go run main.go`

(2) test grpc client `go run main.go`, output id value and token value

(3) open a browser and type in the address `https://127.0.0.1:9090/token/swagger/index.html`, if set isUseTLS=false, use `http`.

request `/v1/user/register`in swagger, response result:

```json
{
  "id": "824351411239129088",
  "token": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiI4MjQzNTE0MTEyMzkxMjkwODgiLCJyb2xlIjoiIiwiZXhwIjoxNjYxMTYxODQ3fQ.wdxs41PSKgvLBMrdCYCrjTIua6ffQY7oBWSpdh4HIfY"
}
```

add token to swagger `Authorize`, and request `/v1/getUser`, response result:

```json
{
  "id": "824351411239129088",
  "name": "foo",
  "email": "foo@bar.com"
}
```

