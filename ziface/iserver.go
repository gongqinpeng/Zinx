package ziface

//定义服务器接口

type IServer interface {
	Start()
	Stop()
	Serve()
	AddRouter(msgId uint32, router IRouter)
	GetConnManager() IConnManager
	SetOnConnStart(func(connection IConnection))
	SetOnConnStop(func(connection IConnection))
	CallOnConnStart(connection IConnection)
	CallOnConnStop(connection IConnection)
}
