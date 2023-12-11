## Example of jwt authentication for GRPC

### Using

Start the grpc server `go run main.go`

1. test grpc client `go run main.go`, output id value and token value

2. test http request. 

- if set isUseTLS=true, visit `https://127.0.0.1:8080/token/swagger/index.html`
- if set isUseTLS=false, visit `http://127.0.0.1:8080/token/swagger/index.html`

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

