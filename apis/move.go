package apis

import (
	"MMO_GAME/core"
	__ "MMO_GAME/pb"
	"ZINX/ziface"
	"ZINX/znet"
	"fmt"
	"google.golang.org/protobuf/proto"
)

type MoveApi struct {
	znet.BaseRouter
}

func (m *MoveApi) Handle(request ziface.IRequest) {
	proto_msg := &__.Position{}
	err := proto.Unmarshal(request.GetData(), proto_msg)
	if err != nil {
		fmt.Println("Move Unmarshal error", err)
		return
	}

	pid, err := request.GetConnection().GetProperty("pid")
	if err != nil {
		fmt.Println("GetProperty pid error", err)
		return
	}
	fmt.Printf("Player pid = %d, move(%f,%f,%f,%f)\n", pid, proto_msg.X, proto_msg.Y, proto_msg.Z, proto_msg.V)

	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))
	//广播并更新当前玩家坐标
	player.UpdatePos(proto_msg.X, proto_msg.Y, proto_msg.Z, proto_msg.V)

}
