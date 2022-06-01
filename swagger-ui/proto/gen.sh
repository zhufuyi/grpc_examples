#!/bin/bash

# 插件版本
# protoc                               v3.20.1
# protoc-gen-openapiv2        v2.10.0

# 服务名称
serverName="hello"

outPath="${serverName}pb"  # 和proto文件的go_package名称一致，也就是文件夹名和包名一致
mkdir -p ${outPath}

# 生成*.swagger.json文件
protoc --openapiv2_opt=logtostderr=true --openapiv2_out=${outPath} *.proto

echo "生成json文件结束"
