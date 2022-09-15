package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	pb "github.com/zhufuyi/grpc_examples/waitForReady/proto/hellopb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func sayHello(ctx context.Context, client pb.GreeterClient, i int) error {
	now := time.Now()

	name := fmt.Sprintf("foo[%d]", i)
	resp, err := client.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		return fmt.Errorf("%s, %v", name, err)
	}

	fmt.Println(resp.Message, time.Since(now).String())
	return nil
}

func main() {
	options := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.Dial("127.0.0.1:8080", options...)
	if err != nil {
		panic(err)
	}
	defer conn.Close() //nolint

	client := pb.NewGreeterClient(conn)

	// 请求超过800毫秒超时
	ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(time.Millisecond*800)) //nolint
	wg := &sync.WaitGroup{}

	// 同时并发5个请求，服务端做了随机0~1000毫秒延时
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			err = sayHello(ctx, client, i)
			if err != nil {
				fmt.Println(err)
			}
		}(i)
	}

	ch := make(chan struct{})

	go func() {
		wg.Wait()
		ch <- struct{}{}
	}()

	// 等待，请求返回小于800毫秒正常，否则出错，超过800毫秒未完成的请求快速失败
	select {
	case <-ch:
	case <-ctx.Done():
		time.Sleep(time.Millisecond)
	}
}
