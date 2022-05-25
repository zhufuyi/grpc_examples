## GRPC tag example

Example of adding a custom tag to the message field in a proto

```protobuf
message HelloRequest {
    string name = 1 [(gogoproto.moretags) = 'gorm:"name"'];
}
```

The plugin `protoc-gen-gogofaster` needs to be installed to replace `protoc-gen-go` and download file `https://github.com/gogo/protobuf/blob/master/gogoproto/gogo.proto`  to be placed in the `$GOPATH /bin/include/gogo/protobuf/gogoproto/`.

