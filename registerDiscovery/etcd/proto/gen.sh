#!/bin/bash

# 插件版本
# protoc                               v3.20.1
# protoc-gen-gogofaster        v1.3.2

# 服务名称
serverName="hello"

outPath="${serverName}pb"  # 和proto文件的go_package名称一致，也就是文件夹名和包名一致
mkdir -p ${outPath}

# 生成pb.go和grpc.pb.go文件，使用protoc-gen-gogofaster插件，支持添加自定义tag，并且序列化和反序列化都比protoc-gen-go更快
protoc --gogofaster_out=${outPath} --gogofaster_opt=paths=source_relative  --go-grpc_out=${outPath} --go-grpc_opt=paths=source_relative *.proto

echo "${serverName}：生成文件结束"
