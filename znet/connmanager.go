package znet

import (
	"ZINX/ziface"
	"errors"
	"fmt"
	"sync"
)

type ConnManager struct {
	connections map[uint32]ziface.IConnection
	connLock    sync.RWMutex
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

func (cm *ConnManager) Add(conn ziface.IConnection) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	cm.connections[conn.GetConnID()] = conn
	fmt.Println("ConnId = ", conn.GetConnID(), "connection add to ConnManager success, conn num = ", cm.Len())
}

func (cm *ConnManager) Remove(conn ziface.IConnection) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	delete(cm.connections, conn.GetConnID())
	fmt.Println("ConnId = ", conn.GetConnID(), "connection remove to ConnManager success, conn num = ", cm.Len())
}

func (cm *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	cm.connLock.RLock()
	defer cm.connLock.RUnlock()
	if conn, ok := cm.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not found!")
	}
}

func (cm *ConnManager) Len() int {
	return len(cm.connections)
}

func (cm *ConnManager) ClearConn() {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	for connID, conn := range cm.connections {
		conn.Stop()
		delete(cm.connections, connID)
	}

	fmt.Println("Clear all connections succ! conn num = ", cm.Len())
}
