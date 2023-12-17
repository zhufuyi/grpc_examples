## Example of GRPC validation

### Using

install plugin `protoc-gen-validate`

> go install github.com/envoyproxy/protoc-gen-validate@v0.6.7

Automatic generation of the validation file `*.validate.go` file

> protoc --validate_opt=paths=source_relative --validate_out=lang=go:. *.proto

Use in server-side methods

```go
err := req.Validate()
if err != nil {
    return nil, status.Errorf(codes.InvalidArgument, "%v", err)
}
```

<br>

Click for more [constraint rules](proto/constraint_rules.md).

<br>

### Running

```bash
# run grpc server
cd server && go run main.go

# run grpc client
cd client && go run main.go
```