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

func (this *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("Call Router PreHandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping..."))
	if err != nil {
		fmt.Println("call back before ping error")
	}

}
func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping ping  ping..."))
	if err != nil {
		fmt.Println("call back ping ping ping error")
	}
}
func (this *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("Call Router PostHandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping..."))
	if err != nil {
		fmt.Println("call back after ping error")
	}
}
func main() {
	s := znet.NewServer("[zinx v0.3]")

	s.AddRouter(&PingRouter{})

	s.Serve()
}
