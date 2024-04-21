package main

import "ZINX/znet"

//基于Zinx框架来开发的服务器端应用程序

func main() {
	s := znet.NewServer("[zinx v0.1]")
	s.Serve()
}
