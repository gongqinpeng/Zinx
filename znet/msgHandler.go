package znet

import (
	"ZINX/utils"
	"ZINX/ziface"
	"fmt"
	"strconv"
)

type MsgHandler struct {
	Apis           map[uint32]ziface.IRouter
	WorkerPoolSize uint32
	TaskQueue      []chan ziface.IRequest
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{Apis: make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize)}
}

func (mh *MsgHandler) DoMsgHandler(request ziface.IRequest) {
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgId=", request.GetMsgID(), "is not found!")
	}

	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

func (mh *MsgHandler) AddRouter(msgID uint32, router ziface.IRouter) {
	if _, ok := mh.Apis[msgID]; ok {
		panic("repeat api,msgId=" + strconv.Itoa(int(msgID)))
	}

	mh.Apis[msgID] = router
	fmt.Println("Add api msgId=", msgID, "succ!")
}

func (mh *MsgHandler) StartWorkerPool() {
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

func (mh *MsgHandler) StartOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Println("worker ID", workerID, "is started...")
	for {
		select {
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

func (mh *MsgHandler) SendMsgToTaskQueue(request ziface.IRequest) {
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("Add ConnID = ", request.GetConnection().GetConnID(), "request MsgID=", request.GetMsgID(),
		"to WorkerId=", workerID)
	mh.TaskQueue[workerID] <- request
}
