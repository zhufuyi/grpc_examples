## English  | [简体中文](readme-cn.md)

### GRPC Examples

This is an example of commonly used knowledge in gRPC, very suitable for users who want to learn grpc in a comprehensive and in-depth way. These knowledge points have been applied to the development framework **sponge**.

[Sponge](https://github.com/zhufuyi/sponge) is a basic development framework that integrates `code auto generation`, `Gin and GRPC framework`. It is easy to build a complete project from development to deployment, just fill in the specific business logic code on the generated template code, the use of Go can also be "low-code development".

Github Repo: [https://github.com/zhufuyi/sponge](https://github.com/zhufuyi/sponge) .

<br>

### Quick Start

#### Go Configuration

This step can be skipped if go has already been set up.

```bash
# Linux or MacOS
export GOROOT="/opt/go"     # your go installation directory
export GOPATH=$HOME/go      # Set the directory where the "go get" command downloads third-party packages
export GOBIN=$GOPATH/bin    # Set the directory where the executable binaries are compiled by the "go install" command.
export PATH=$PATH:$GOBIN:$GOROOT/bin  # Add the $GOBIN directory to the system path.


# Windows
setx GOPATH "D:\your-directory"      # Set the directory where the "go get" command downloads third-party packages
setx GOBIN "D:\your-directory\bin"   # Set the directory where the executable binary files generated by the "go install" command are stored.
```

<br>

#### Installation of Protoc and Plugins

1. Copy the protobuf file dependency folder [include](include) to the `$GOBIN` directory.

2. Download protoc from: [https://github.com/protocolbuffers/protobuf/releases/tag/v25.2](https://github.com/protocolbuffers/protobuf/releases/tag/v25.2)

> Download the protoc binaries according to the system type, move the protoc binaries to `$GOBIN`.

3. Install Protoc Plugins

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
go install github.com/envoyproxy/protoc-gen-validate@latest
go install github.com/mohuishou/protoc-gen-go-gin@latest
go install github.com/srikrsna/protoc-gen-gotag@latest
```

<br>

### List of examples

- [Serialization and deserialization of protobuf](https://github.com/zhufuyi/grpc_examples/tree/main/protobuf)
- [4 ways of calling in helloworld demo](https://github.com/zhufuyi/grpc_examples/tree/main/helloworld)
- [interceptor](https://github.com/zhufuyi/grpc_examples/tree/main/interceptor)
- [recovery](https://github.com/zhufuyi/grpc_examples/tree/main/recovery)
- [logging](https://github.com/zhufuyi/grpc_examples/tree/main/logging)
- [keepalive](https://github.com/zhufuyi/grpc_examples/tree/main/keepalive)
- [metadata set and read](https://github.com/zhufuyi/grpc_examples/tree/main/metadata)
- [timeout](https://github.com/zhufuyi/grpc_examples/tree/main/timeout)
- [swagger](https://github.com/zhufuyi/grpc_examples/tree/main/swagger/example)
- [tag](https://github.com/zhufuyi/grpc_examples/tree/main/tag)
- [validate](https://github.com/zhufuyi/grpc_examples/tree/main/validate)
- [wait for ready](https://github.com/zhufuyi/grpc_examples/tree/main/waitForReady)
- [http to grpc](https://github.com/zhufuyi/grpc_examples/tree/main/httpToGrpc)
  - [call grpc in gin](https://github.com/zhufuyi/grpc_examples/tree/main/httpToGrpc/ginToGrpc)
  - [grpc gateway](https://github.com/zhufuyi/grpc_examples/tree/main/httpToGrpc/grpc-gateway)
- [security](https://github.com/zhufuyi/grpc_examples/tree/main/security)
  - [TLS server-side authentication and two-way authentication](https://github.com/zhufuyi/grpc_examples/tree/main/security/tls)
  - [kv token authentication](https://github.com/zhufuyi/grpc_examples/tree/main/security/kv_token)
  - [jwt token authentication](https://github.com/zhufuyi/grpc_examples/tree/main/security/jwt_token)
- [register and discovery](https://github.com/zhufuyi/grpc_examples/tree/main/registerDiscovery)
  - [consul](https://github.com/zhufuyi/grpc_examples/tree/main/registerDiscovery/consul)
  - [etcd](https://github.com/zhufuyi/grpc_examples/tree/main/registerDiscovery/etcd)
  - [nacos](https://github.com/zhufuyi/grpc_examples/tree/main/registerDiscovery/nacos)
- [load-balance](https://github.com/zhufuyi/grpc_examples/tree/main/loadbalance)
  - [ip loadbalance](https://github.com/zhufuyi/grpc_examples/tree/main/loadbalance/ipAddr)
  - [consul_loadbalance](https://github.com/zhufuyi/grpc_examples/tree/main/loadbalance/consul)
  - [etcd_loadbalance](https://github.com/zhufuyi/grpc_examples/tree/main/loadbalance/etcd)
  - [nacos_loadbalance](https://github.com/zhufuyi/grpc_examples/tree/main/loadbalance/nacos)
- [ratelimit](https://github.com/zhufuyi/grpc_examples/tree/main/ratelimit)
- [breaker](https://github.com/zhufuyi/grpc_examples/tree/main/breaker)
- [retry](https://github.com/zhufuyi/grpc_examples/tree/main/retry)
- [metrics](https://github.com/zhufuyi/grpc_examples/tree/main/metrics)
  - [default grpc metrics](https://github.com/zhufuyi/grpc_examples/tree/main/metrics/defaultMetrics)
  - [customized grpc metrics](https://github.com/zhufuyi/grpc_examples/tree/main/metrics/customMetrics)
- [tracing](https://github.com/zhufuyi/grpc_examples/tree/main/tracing)
  - [http-->grpc tracing](https://github.com/zhufuyi/grpc_examples/tree/main/tracing/http2rpc)
  - [grpc-->grpc tracing](https://github.com/zhufuyi/grpc_examples/tree/main/tracing/rpc2rpc)
- [practical project use](https://github.com/zhufuyi/grpc_examples/tree/main/usage)

<br>
