package ziface

type IMsgHandler interface {
	DoMsgHandler(request IRequest)
	AddRouter(msgID uint32, router IRouter)
	StartWorkerPool()
	SendMsgToTaskQueue(request IRequest)
}
