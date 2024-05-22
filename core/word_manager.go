package core

import (
	"fmt"
	"sync"
)

type WordManager struct {
	AoiMgr  *AOIManager
	Players map[int32]*Player
	pLock   sync.RWMutex
}

var WorldMgrObj *WordManager

func init() {
	WorldMgrObj = &WordManager{
		AoiMgr:  NewAOIManager(AOI_MIN_X, AOI_MAX_X, AOI_CNTS_X, AOI_MIN_Y, AOI_MAX_Y, AOI_CNTS_Y), //创建世界aoi地图
		Players: make(map[int32]*Player),
		pLock:   sync.RWMutex{},
	}
}

func (wm *WordManager) AddPlayer(player *Player) {
	wm.pLock.Lock()
	wm.Players[player.Pid] = player
	wm.pLock.Unlock()

	wm.AoiMgr.AddToGridByPos(int(player.Pid), player.X, player.Z)
}

func (wm *WordManager) RemovePlayer(pid int32) {
	wm.pLock.Lock()
	if _, ok := wm.Players[pid]; !ok {
		fmt.Println("The player is not exist!")
		return
	}

	wm.AoiMgr.RemoveFromGridByPos(int(pid), wm.Players[pid].X, wm.Players[pid].Z)
	delete(wm.Players, pid)
	wm.pLock.Unlock()
}

func (wm *WordManager) GetPlayerByPid(pid int32) *Player {
	wm.pLock.RLock()
	defer wm.pLock.RUnlock()
	if player, ok := wm.Players[pid]; !ok {
		fmt.Println("The player is not exist!")
		return nil
	} else {
		return player
	}
}

func (wm *WordManager) GetAllPlayers() []*Player {
	wm.pLock.RLock()
	defer wm.pLock.RUnlock()
	players := make([]*Player, 0)

	for _, v := range wm.Players {
		players = append(players, v)
	}
	return players
}
