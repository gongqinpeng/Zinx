package ziface

import "net"

type IConnection interface {
	Start()
	Stop()
	GetTCPConnection() *net.TCPConn
	GetConnID() uint32
	RemoteAddr() net.Addr
	SendMsg(msgId uint32, data []byte) error
	SetProperty(key string, value interface{})
	GetProperty(key string) (interface{}, error)
	RemoveProperty(key string)
}

type HandleFunc func(*net.TCPConn, []byte, int) error
