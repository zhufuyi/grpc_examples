#!/bin/bash

# 插件版本
# protoc                               v3.20.1
# protoc-gen-go                   v1.28.0

# 服务名称
serverName="account"

outPath="${serverName}pb"  # 和proto文件的go_package名称一致，也就是文件夹名和包名一致
mkdir -p ${outPath}

# 生成pb.go和grpc.pb.go文件，如果要兼容旧版本protoc-gen-go生成代码，需要添加参数--go-grpc_opt=require_unimplemented_servers=false，或者在实现接口添加pb.Unimplemented***Server
protoc --go_out=${outPath} --go_opt=paths=source_relative --go-grpc_out=${outPath} --go-grpc_opt=paths=source_relative *.proto

# 生成validate文件
protoc --validate_opt=paths=source_relative --validate_out=lang=go:${outPath} *.proto

echo "${serverName}：生成文件结束"
