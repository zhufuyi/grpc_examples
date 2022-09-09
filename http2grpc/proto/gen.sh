#!/bin/bash

# 插件版本
# protoc                               v3.20.1
# protoc-gen-go                   v1.28.0
# protoc-gen-go-grpc            v1.2.0
# protoc-gen-grpc-gateway   v2.10.0
# protoc-gen-openapiv2        v2.10.0

# 服务名称
serverName="account"

outPath="pb"  # 和proto文件的go_package名称一致，也就是文件夹名和包名一致
mkdir -p ${outPath}

# 生成pb.go和grpc.pb.go文件，如果要兼容旧版本protoc-gen-go生成代码，需要添加参数--go-grpc_opt=require_unimplemented_servers=false，或者在实现接口添加pb.Unimplemented***Server
protoc --go_out=${outPath} --go_opt=paths=source_relative --go-grpc_out=${outPath} --go-grpc_opt=paths=source_relative *.proto

# 生成pb.gw.go文件，http的api接口文件
protoc --grpc-gateway_opt=paths=source_relative --grpc-gateway_out=${outPath} *.proto

# 生成swagger.json文件
protoc --openapiv2_opt=logtostderr=true --openapiv2_out=${outPath} *.proto

echo "${serverName}：生成文件结束"
