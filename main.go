package main

import (
	"context"
	"flag"
	"fmt"
	. "github.com/invictus555/greeting/api/greeting"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	port *int    // greetingServer的监听端口
	addr *string // greetingServer的监听地址
)

func init() {
	// 获取GreetingServer的监听端口与监听地址
	port = flag.Int("port", 50050, "the port of greeting server")
	addr = flag.String("svrAddr", "localhost", "the address of greeting server")
	flag.Parse() // 解析命令行参数
}

func main() {
	// 连接gRPC服务器
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", *addr, *port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("Failed to dial to: " + fmt.Sprintf("%s:%d", *addr, *port))
	}

	if conn == nil {
		panic("Failed to get connection from the server")
	}

	defer conn.Close()
	c := NewGreetingClient(conn) // 初始化客户端

	// 调用方法
	reqBody := new(HelloRequest)
	reqBody.Name = "gRPC Server"
	r, err := c.SayHello(context.Background(), reqBody)
	if err != nil {
		fmt.Printf("Fail to call SayHello method, err: %v", err)
	}

	fmt.Println(r.Message)

	empty := new(emptypb.Empty)
	bye, err := c.SayBye(context.Background(), empty)
	if err != nil {
		fmt.Printf("Fail to call SayBye method, err: %v", err)
	}
	fmt.Println(bye.Message)
}
