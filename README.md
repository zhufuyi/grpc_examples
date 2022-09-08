## grpc_examples

### Installation

(1) Copy the protobuf file dependency folder [include](include) to the `$GOPATH/bin` directory.

(2) Store all the downloaded plugins in `$GOPATH/bin` directory.

```bash
# Various plug-in versions
# protoc                    v3.20.1      command
# protoc-gen-go             v1.28.0      plugin, generate *.pb.go file based on proto files, which are populated, serialized and retrieved message type code.
# protoc-gen-gogofast     v1.28.0      plugin, generate *.pb.go file based on proto files, replaces protoc-gen-go plugin for faster encoding and decoding, custom tags are also supported.
# protoc-gen-go-grpc        v1.2.0       plugin, generate *_grpc.pb.go file based on proto files, which are client-side and server-side method and interface code.
# protoc-gen-grpc-gateway   v2.10.0      plugin, generate *.pb.gw.go file based on proto file, which is the api code for web.
# protoc-gen-openapiv2      v2.10.0      plugin, generate *.swagger.json file based on proto file, which is swagger-ui interface documentation.
# protoc-gen-validate       v0.6.7       plugin, generate *.pb.validate.go file according to proto file, is the check field code
# protoc-gen-go             v0.1.0      plugin, generate *gin.pb.go file based on proto files, which is gin handler.

# download  protoc url
https://github.com/protocolbuffers/protobuf/releases/tag/v3.20.1

# install plugin protoc-gen-go, protoc-gen-go-grpc, protoc-gen-validate, protoc-gen-gogofast, protoc-gen-gogofaster
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install github.com/gogo/protobuf/protoc-gen-gogofaster@v1.3.2
go install github.com/gogo/protobuf/protoc-gen-gogofast@v1.3.2
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
go install github.com/envoyproxy/protoc-gen-validate@v0.6.7
go install github.com/mohuishou/protoc-gen-go-gin@v0.1.0

# download pugin protoc-gen-grpc-gateway,protoc-gen-openapiv2 url, total 2 files.
https://github.com/grpc-ecosystem/grpc-gateway/releases/tag/v2.10.1
```

<br>

### List of examples

- [Serialization and deserialization of protobuf](protobuf)
- [4 ways of calling in helloworld](helloworld)
- [interceptor](interceptor)
- [recovery](recovery)
- [logging](logging)
- [keepalive](keepalive)
- [metadata set and read](metadata)
- [timeout](timeout)
- [swagger](swagger)
- [tag](tag)
- [validate](validate)
- [waitForReady](waitForReady)
- [grpc-gateway](http2grpc)
- [gin call rpc](gin2grpc)
- [security](security)
  - [TLS server-side authentication and two-way authentication](security/tls)
  - [kv token authentication](security/kv_token)
  - [jwt token authentication](security/jwt_token)
- [etcd register and discovery](registerDiscovery)
- [load-balance](loadbalance)
  - [cliend side loadbalance](loadbalance/client_loadbalance)
  - [loadbalance use with etcd](loadbalance/etcd_loadbalance)
- [ratelimit](ratelimit)
- [retry](retry)
- [hystrix](hystrix)
  - [base use](hystrix/baseuse)
  - [with metrics](hystrix/withMetrics)
- [metrics](metrics)
  - [default grpc metrics](metrics/defaultMetrics)
  - [customized grpc metrics](metrics/customizedMetrics)
- [tracing](tracing)
  - [gin-->rpc tracing](tracing/http2rpc)
  - [rpc-->rpc tracing](tracing/rpc2rpc)

<br>
