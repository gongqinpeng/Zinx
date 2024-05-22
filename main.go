package main

import (
	"MMO_GAME/apis"
	"MMO_GAME/core"
	"ZINX/ziface"
	"ZINX/znet"
	"fmt"
)

func OnConnectionAdd(conn ziface.IConnection) {
	//给客户端发送MsgId 1 200的消息
	player := core.NewPlayer(conn)
	player.SyncPid()
	player.BroadCastStartPosition()
	core.WorldMgrObj.AddPlayer(player)
	conn.SetProperty("pid", player.Pid) // 创建玩家时给连接绑定pid属性

	//同步周边玩家
	player.SyncSurrounding()

	fmt.Println("----->Player pid = ", player.Pid, "is arrived-----<")
}

func OnConnectionLost(conn ziface.IConnection) {
	pid, _ := conn.GetProperty("pid")
	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))

	player.Offline()
	fmt.Println("----->Player pid =", player.Pid, "is offline-------<")
}

func main() {
	s := znet.NewServer("MMO_ZINX_GAME")

	//连接创建和销毁的HOOK
	s.SetOnConnStart(OnConnectionAdd)
	s.SetOnConnStop(OnConnectionLost)
	s.AddRouter(2, &apis.WorldChatApi{})
	s.AddRouter(3, &apis.MoveApi{})
	s.Serve()

}
