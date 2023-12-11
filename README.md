## grpc_examples

### Installation

(1) Copy the protobuf file dependency folder [include](include) to the `$GOPATH/bin` directory.

(2) Store all the downloaded plugins in `$GOPATH/bin` directory.

```bash
# install protoc in linux
mkdir -p protocDir \
  && curl -L -o protocDir/protoc.zip https://github.com/protocolbuffers/protobuf/releases/download/v3.20.1/protoc-3.20.1-linux-x86_64.zip \
  && unzip protocDir/protoc.zip -d protocDir\
  && mv protocDir/bin/protoc protocDir/include/ $GOROOT/bin/ \
  && rm -rf protocDir

# install plugin
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.0
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.10.0
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.10.0
go install github.com/envoyproxy/protoc-gen-validate@v0.6.7
go install github.com/mohuishou/protoc-gen-go-gin@v0.1.0
go install github.com/srikrsna/protoc-gen-gotag@v0.6.2
go install github.com/gogo/protobuf/protoc-gen-gogofast@v1.3.2
go install github.com/gogo/protobuf/protoc-gen-gogofaster@v1.3.2

# Various plug-in versions
# protoc                    v3.20.1      command
# protoc-gen-go             v1.28.0      plugin, generate *.pb.go file based on proto files, which are populated, serialized and retrieved message type code.
# protoc-gen-go-grpc        v1.2.0       plugin, generate *_grpc.pb.go file based on proto files, which are client-side and server-side method and interface code.
# protoc-gen-grpc-gateway   v2.10.0      plugin, generate *.pb.gw.go file based on proto file, which is the api code for web.
# protoc-gen-openapiv2      v2.10.0      plugin, generate *.swagger.json file based on proto file, which is swagger-ui interface documentation.
# protoc-gen-validate       v0.6.7       plugin, generate *.pb.validate.go file according to proto file, is the check field code
# protoc-gen-go-gin             v0.1.0      plugin, generate *gin.pb.go file based on proto files, which is gin handler.
# protoc-gen-gogofast     v1.28.0      plugin, generate *.pb.go file based on proto files, replaces protoc-gen-go plugin for faster encoding and decoding, custom tags are also supported.
```

<br>

### List of examples

- [Serialization and deserialization of protobuf](protobuf)
- [4 ways of calling in helloworld demo](helloworld)
- [interceptor](interceptor)
- [recovery](recovery)
- [logging](logging)
- [keepalive](keepalive)
- [metadata set and read](metadata)
- [timeout](timeout)
- [swagger](swagger/example)
- [tag](tag)
- [validate](validate)
- [waitForReady](waitForReady)
- [http to grpc](httpToGrpc)
  - [call grpc in gin](httpToGrpc/ginToGrpc)
  - [grpc gateway](httpToGrpc/grpc-gateway)
- [security](security)
  - [TLS server-side authentication and two-way authentication](security/tls)
  - [kv token authentication](security/kv_token)
  - [jwt token authentication](security/jwt_token)
- [register and discovery](registerDiscovery)
  - [consul](registerDiscovery/consul)
  - [etcd](registerDiscovery/etcd)
  - [nacos](registerDiscovery/nacos)
- [load-balance](loadbalance)
  - [ip loadbalance](loadbalance/ipAddr)
  - [consul_loadbalance](loadbalance/consul)
  - [etcd_loadbalance](loadbalance/etcd)
  - [nacos_loadbalance](loadbalance/nacos)
- [ratelimit](ratelimit)
- [retry](retry)
- [breaker](breaker)
- [metrics](metrics)
  - [default grpc metrics](metrics/defaultMetrics)
  - [customized grpc metrics](metrics/customMetrics)
- [tracing](tracing)
  - [http-->grpc tracing](tracing/http2rpc)
  - [grpc-->grpc tracing](tracing/rpc2rpc)

<br>
