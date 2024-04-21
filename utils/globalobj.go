package utils

import (
	"ZINX/ziface"
	"encoding/json"
	"io/ioutil"
)

type GlobalObj struct {
	TcpServer ziface.IServer
	Host      string
	TcpPort   int
	Name      string

	Version        string
	MaxConn        int
	MaxPackageSize uint32

	WorkerPoolSize   uint32
	MaxWorkerTaskLen uint32
}

var GlobalObject *GlobalObj

func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("Demo/Zinx1.0/conf/zinx.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

func init() {
	GlobalObject = &GlobalObj{
		TcpServer:        nil,
		Host:             "0.0.0.0",
		TcpPort:          8999,
		Name:             "ZinxServerApp",
		Version:          "V0.8",
		MaxConn:          1000,
		MaxPackageSize:   4096,
		WorkerPoolSize:   10,   //worker 工作池的队列的个数
		MaxWorkerTaskLen: 1024, //每个worker对应的消息队列的任务的数量最大值
	}

	GlobalObject.Reload()
}
