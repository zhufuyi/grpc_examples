#!/bin/bash

# 插件版本
# protoc                               v3.20.1
# protoc-gen-go                   v1.28.0

# 服务名称
serverName="rpcerror"

outPath="${serverName}pb"  # 和proto文件的go_package名称一致，也就是文件夹名和包名一致
mkdir -p ${outPath}

# 生成pb.go文件
protoc --go_out=${outPath} --go_opt=paths=source_relative *.proto

echo "${serverName}：生成文件结束"
