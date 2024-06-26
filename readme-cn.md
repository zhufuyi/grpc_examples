## [English](README.md)  | 简体中文

### 使用前说明

这是grpc常用的知识点示例，非常适合全面深入学习和使用grpc，这些grpc知识点是从go开发框架 **sponge** 提取出来的。

[sponge](https://github.com/zhufuyi/sponge) 是一个集成了`自动生成代码`、`gin和grpc框架`的强大的开发框架。从生成代码、开发、测试、api文档到部署一站式项目开发，大幅提高了开发效率和降低了开发难度，实现"低代码方式"进行开发项目。

github 地址: [https://github.com/zhufuyi/sponge](https://github.com/zhufuyi/sponge)

<br>

### 快速开始

#### Go设置

如果已经设置过go可以跳过此步骤。

```bash
# Linux 或 MacOS
export GOROOT="/opt/go"     # 你的go安装目录
export GOPATH=$HOME/go      # 设置 go get 命令下载第三方包的目录
export GOBIN=$GOPATH/bin    # 设置 go install 命令编译后生成可执行二进制文件的存放目录
export PATH=$PATH:$GOBIN:$GOROOT/bin   # 把$GOBIN目录添加到系统path


# Windows
setx GOPATH "D:\你的目录"     # 设置 go get 命令下载第三方包的目录
setx GOBIN "D:\你的目录\bin"   # 设置 go install 命令编译后生成可执行二进制文件的存放目录
```

<br>

#### 安装protoc和插件

1. 复制目录 [include](include) 到 `$GOBIN`.

2. 下载protoc地址: [https://github.com/protocolbuffers/protobuf/releases/tag/v25.2](https://github.com/protocolbuffers/protobuf/releases/tag/v25.2)

> 根据系统类型下载对应的protoc二进制文件，并把二进制文件移动到 `$GOBIN`.

3. 安装protoc插件

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

### 所有示例列表

- [protobuf 的序列化和反序列化](https://github.com/zhufuyi/grpc_examples/tree/main/protobuf)
- [protobuf tag](https://github.com/zhufuyi/grpc_examples/tree/main/tag)
- [grpc 四种调用方式](https://github.com/zhufuyi/grpc_examples/tree/main/helloworld)
- [grpc 拦截器](https://github.com/zhufuyi/grpc_examples/tree/main/interceptor)
- [grpc 恢复](https://github.com/zhufuyi/grpc_examples/tree/main/recovery)
- [grpc 日志](https://github.com/zhufuyi/grpc_examples/tree/main/logging)
- [grpc keepalive](https://github.com/zhufuyi/grpc_examples/tree/main/keepalive)
- [grpc 元数据传递与读写](https://github.com/zhufuyi/grpc_examples/tree/main/metadata)
- [grpc 超时](https://github.com/zhufuyi/grpc_examples/tree/main/timeout)
- [grpc 参数校验](https://github.com/zhufuyi/grpc_examples/tree/main/validate)
- [grpc 鉴权](https://github.com/zhufuyi/grpc_examples/tree/main/security)
  - [TLS 鉴权](https://github.com/zhufuyi/grpc_examples/tree/main/security/tls)
  - [kv token 鉴权](https://github.com/zhufuyi/grpc_examples/tree/main/security/kv_token)
  - [jwt 鉴权](https://github.com/zhufuyi/grpc_examples/tree/main/security/jwt_token)
- [grpc 注册与发现](https://github.com/zhufuyi/grpc_examples/tree/main/registerDiscovery)
  - [consul](https://github.com/zhufuyi/grpc_examples/tree/main/registerDiscovery/consul)
  - [etcd](https://github.com/zhufuyi/grpc_examples/tree/main/registerDiscovery/etcd)
  - [nacos](https://github.com/zhufuyi/grpc_examples/tree/main/registerDiscovery/nacos)
- [grpc 负载均衡](https://github.com/zhufuyi/grpc_examples/tree/main/loadbalance)
  - [ip](https://github.com/zhufuyi/grpc_examples/tree/main/loadbalance/ipAddr)
  - [consul](https://github.com/zhufuyi/grpc_examples/tree/main/loadbalance/consul)
  - [etcd](https://github.com/zhufuyi/grpc_examples/tree/main/loadbalance/etcd)
  - [nacos](https://github.com/zhufuyi/grpc_examples/tree/main/loadbalance/nacos)
- [grpc 限流](https://github.com/zhufuyi/grpc_examples/tree/main/ratelimit)
- [grpc 熔断](https://github.com/zhufuyi/grpc_examples/tree/main/breaker)
- [grpc 重试](https://github.com/zhufuyi/grpc_examples/tree/main/retry)
- [grpc 链路跟踪](https://github.com/zhufuyi/grpc_examples/tree/main/tracing)
  - [http --> grpc](https://github.com/zhufuyi/grpc_examples/tree/main/tracing/http2rpc)
  - [grpc --> grpc](https://github.com/zhufuyi/grpc_examples/tree/main/tracing/rpc2rpc)
- [grpc 指标](https://github.com/zhufuyi/grpc_examples/tree/main/metrics)
  - [默认 grpc metrics](https://github.com/zhufuyi/grpc_examples/tree/main/metrics/defaultMetrics)
  - [自定义 grpc metrics](https://github.com/zhufuyi/grpc_examples/tree/main/metrics/customMetrics)
- [在http中调用grpc](https://github.com/zhufuyi/grpc_examples/tree/main/httpToGrpc)
  - [在gin中调用grpc](https://github.com/zhufuyi/grpc_examples/tree/main/httpToGrpc/ginToGrpc)
  - [grpc 网关](https://github.com/zhufuyi/grpc_examples/tree/main/httpToGrpc/grpc-gateway)
- [grpc 封装实践](https://github.com/zhufuyi/grpc_examples/tree/main/usage)

<br>
