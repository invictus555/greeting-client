package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	. "github.com/invictus555/greeting/api/greeting"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net/http"
)

var (
	httpPort     *int    // http server服务端口
	greetingPort *int    // greetingServer的监听端口
	greetingAddr *string // greetingServer的监听地址
)

func init() {
	// 获取GreetingServer的监听端口与监听地址
	httpPort = flag.Int("httpServerPort", 8080, "the port of http server")
	greetingPort = flag.Int("greetingServerPort", 50050, "the port of greeting server")
	greetingAddr = flag.String("greetingServerAddress", "localhost", "the address of greeting server")
	flag.Parse() // 解析命令行参数
}

func NewGreetingServer(addr string) *grpc.ClientConn {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials())) // 连接gRPC服务器
	if err != nil {
		fmt.Printf("Failed to dial to: " + addr)
	}

	if conn == nil {
		panic("Failed to get connection from the server")
	}

	return conn
}

func goodBye(w http.ResponseWriter, r *http.Request) {
	greetingServer := NewGreetingServer(fmt.Sprintf("%s:%d", *greetingAddr, *greetingPort))
	c := NewGreetingClient(greetingServer) // 初始化客户端

	empty := new(emptypb.Empty)
	resp, err := c.SayBye(context.Background(), empty)
	if err != nil {
		fmt.Printf("Fail to call SayBye method, err: %v", err)
	}

	greetingServer.Close()
	fmt.Println(resp.Message)    // 打印消息
	w.WriteHeader(http.StatusOK) // 设置响应状态码为 200
	fmt.Fprintf(w, resp.Message) // 发送响应到客户端
}

func sayHello(w http.ResponseWriter, r *http.Request) {
	greetingServer := NewGreetingServer(fmt.Sprintf("%s:%d", *greetingAddr, *greetingPort))
	c := NewGreetingClient(greetingServer) // 初始化客户端

	reqBody := new(HelloRequest)
	reqBody.Name = "gRPC Server"
	resp, err := c.SayHello(context.Background(), reqBody) // 调用方法
	if err != nil {
		fmt.Printf("Fail to call SayHello method, err: %v", err)
	}

	greetingServer.Close()
	fmt.Println(resp.Message)    // 打印消息
	w.WriteHeader(http.StatusOK) // 设置响应状态码为 200
	fmt.Fprintf(w, resp.Message) // 发送响应到客户端
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/bye", goodBye).Methods(http.MethodGet)
	r.HandleFunc("/hello", sayHello).Methods(http.MethodGet)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *httpPort), r))
}
