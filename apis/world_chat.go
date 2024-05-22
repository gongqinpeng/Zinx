package apis

import (
	"MMO_GAME/core"
	__ "MMO_GAME/pb"
	"ZINX/ziface"
	"ZINX/znet"
	"fmt"
	"google.golang.org/protobuf/proto"
)

//世界聊天的路由

type WorldChatApi struct {
	znet.BaseRouter
}

func (wc *WorldChatApi) Handle(request ziface.IRequest) {
	//解析客户端传来的proto协议
	proto_msg := &__.Talk{}
	err := proto.Unmarshal(request.GetData(), proto_msg)
	if err != nil {
		fmt.Println("Talk Unmarshal err", err)
		return
	}
	//当前聊天是哪个玩家发送的
	pid, err := request.GetConnection().GetProperty("pid")

	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))

	player.Talk(proto_msg.Content)
}
