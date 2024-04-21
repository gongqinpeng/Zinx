package main

import (
	"ZINX/znet"
	"fmt"
)
import "ZINX/ziface"

//基于Zinx框架来开发的服务器端应用程序

type PingRouter struct {
	znet.BaseRouter
}

//	func (this *PingRouter) PreHandle(request ziface.IRequest) {
//		fmt.Println("Call Router PreHandle...")
//		_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping..."))
//		if err != nil {
//			fmt.Println("call back before ping error")
//		}
//
// }
func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call PingRouter Handle...")
	fmt.Println("recv from client msgId=", request.GetMsgID(), "data = ", string(request.GetData()))

	err := request.GetConnection().SendMsg(200, []byte("ping...ping..."))
	if err != nil {
		fmt.Println(err)
	}
}

type HelloZinxRouter struct {
	znet.BaseRouter
}

func (this *HelloZinxRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call HelloRouter Handle...")
	fmt.Println("recv from client msgId=", request.GetMsgID(), "data = ", string(request.GetData()))

	err := request.GetConnection().SendMsg(201, []byte("Hello Zinx"))
	if err != nil {
		fmt.Println(err)
	}
}

//	func (this *PingRouter) PostHandle(request ziface.IRequest) {
//		fmt.Println("Call Router PostHandle...")
//		_, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping..."))
//		if err != nil {
//			fmt.Println("call back after ping error")
//		}
//	}
func main() {
	s := znet.NewServer("[zinx v0.6]")

	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloZinxRouter{})

	s.Serve()
}
