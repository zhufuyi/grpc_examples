package main

import (
	"encoding/json"
	"fmt"

	"github.com/golang/protobuf/proto"
	pb "github.com/zhufuyi/grpc_examples/protobuf/proto/userpb"
)

func main() {
	info := &pb.User{
		Id:    123,
		Name:  "foo",
		Email: "foo@bar.com",
	}
	fmt.Println("source data:", info)

	// json序列化数据
	jsonData, _ := json.Marshal(info)
	fmt.Printf("json data: %X,  size:%d\n", jsonData, len(jsonData))

	// 序列化
	data, err := proto.Marshal(info)
	if err != nil {
		panic(err)
	}
	fmt.Printf("  pb data: %X,  size:%d\n", data, len(data))

	// 反序列化，还原出原始数据
	info2 := &pb.User{}
	err = proto.Unmarshal(data, info2)
	if err != nil {
		panic(err)
	}
	fmt.Println("proto unmarshal data:", info2)
}

/*
输出：
source data: id:123  name:"abc"  email:"abc@123.com"
json data: 7B226964223A3132332C226E616D65223A22616263222C22656D61696C223A22616263403132332E636F6D227D,  size:45
  pb data: 087B12036162631A0B616263403132332E636F6D,  size:20
proto unmarshal data: id:123  name:"abc"  email:"abc@123.com"
*/
