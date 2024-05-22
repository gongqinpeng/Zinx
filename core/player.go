package core

import (
	__ "MMO_GAME/pb"
	"ZINX/ziface"
	"fmt"
	"google.golang.org/protobuf/proto"
	"math/rand"
	"sync"
)

type Player struct {
	Pid  int32
	Conn ziface.IConnection //当前玩家与服务器的链接
	X    float32            // 平面X坐标
	Y    float32            // 高度
	Z    float32            // 平面Y坐标
	V    float32            // 旋转角度 0-360
}

var PidGen int32 = 1 //ID计数器
var IdLock sync.Mutex

func NewPlayer(conn ziface.IConnection) *Player {
	IdLock.Lock()
	id := PidGen
	PidGen++
	IdLock.Unlock()

	p := &Player{
		Pid:  id,
		Conn: conn,
		X:    float32(160 + rand.Intn(10)), //随机在160坐标点 基于X若干偏移
		Y:    0,
		Z:    float32(140 + rand.Intn(20)),
		V:    0,
	}

	return p
}

func (p *Player) SendMsg(msgId uint32, data proto.Message) {
	msg, err := proto.Marshal(data)
	if err != nil {
		fmt.Println("marshal err", err)
		return
	}

	if p.Conn == nil {
		fmt.Println("Conn is null")
		return
	}

	if err := p.Conn.SendMsg(msgId, msg); err != nil {
		fmt.Println("Player send msg error", err)
		return
	}
	return
}

func (p *Player) SyncPid() {
	data := &__.SyncPid{
		Pid: p.Pid,
	}
	p.SendMsg(1, data) //告知玩家id msgid = 1
}

func (p *Player) BroadCastStartPosition() {
	data := &__.BroadCast{
		Pid: p.Pid,
		Tp:  2, //广播位置
		Data: &__.BroadCast_P{
			P: &__.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}

	p.SendMsg(200, data) //告知玩家位置 msgid = 200
}

func (p *Player) Talk(content string) {
	proto_msg := &__.BroadCast{
		Pid:  p.Pid,
		Tp:   1, //聊天广播
		Data: &__.BroadCast_Content{Content: content},
	}

	players := WorldMgrObj.GetAllPlayers()
	for _, player := range players {
		player.SendMsg(200, proto_msg)
	}
}

func (p *Player) SyncSurrounding() {
	//获取周围玩家
	//pids := WorldMgrObj.AoiMgr.GetPidsByPos(p.X, p.Z)
	//
	//players := make([]*Player, 0, len(pids))
	//for _, pid := range pids {
	//	players = append(players, WorldMgrObj.GetPlayerByPid(int32(pid)))
	//}
	players := p.GetSurroundingPlayers()
	//将当前玩家位置通过msgid200发送给周围玩家
	proto_msg := &__.BroadCast{
		Pid: p.Pid,
		Tp:  2, // 广播坐标
		Data: &__.BroadCast_P{
			P: &__.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			}},
	}

	for _, player := range players {
		player.SendMsg(200, proto_msg)
	}

	//周围玩家发送信息给当前玩家msgid202
	players_proto_msg := make([]*__.Player, 0, len(players))
	for _, player := range players {
		p := &__.Player{
			Pid: player.Pid,
			P: &__.Position{
				X: player.X,
				Y: player.Y,
				Z: player.Z,
				V: player.V,
			},
		}
		players_proto_msg = append(players_proto_msg, p)
	}

	SyncPlayers_proto_msg := &__.SyncPlayers{Ps: players_proto_msg[:]}

	p.SendMsg(202, SyncPlayers_proto_msg)

}

// 广播当前玩家位置移动信息
func (p *Player) UpdatePos(x, y, z, v float32) {
	p.X = x
	p.Y = y
	p.Z = z
	p.V = v

	proto_msg := &__.BroadCast{
		Pid: p.Pid,
		Tp:  4, //表示移动后坐标
		Data: &__.BroadCast_P{P: &__.Position{
			X: p.X,
			Y: p.Y,
			Z: p.Z,
			V: p.V,
		}},
	}

	players := p.GetSurroundingPlayers()

	for _, player := range players {
		player.SendMsg(200, proto_msg)
	}

}

func (p *Player) GetSurroundingPlayers() []*Player {
	pids := WorldMgrObj.AoiMgr.GetPidsByPos(p.X, p.Z)

	players := make([]*Player, 0, len(pids))
	for _, pid := range pids {
		players = append(players, WorldMgrObj.GetPlayerByPid(int32(pid)))
	}
	return players
}

func (p *Player) Offline() {
	players := p.GetSurroundingPlayers()
	proto_msg := &__.SyncPid{Pid: p.Pid}

	for _, player := range players {
		player.SendMsg(201, proto_msg)
	}

	WorldMgrObj.AoiMgr.RemoveFromGridByPos(int(p.Pid), p.X, p.Z)
	WorldMgrObj.RemovePlayer(p.Pid)
}
