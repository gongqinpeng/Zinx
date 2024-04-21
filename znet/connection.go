package znet

import (
	"ZINX/utils"
	"ZINX/ziface"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
)

type Connection struct {
	//当前connection属于哪个server
	TcpServer ziface.IServer
	Conn      *net.TCPConn
	ConnID    uint32
	isClosed  bool

	ExitChan chan bool

	msgChan    chan []byte
	MsgHandler ziface.IMsgHandler
	//链接属性集合
	property     map[string]interface{}
	propertyLock sync.RWMutex
}

func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandler) *Connection {
	c := &Connection{
		TcpServer: server,
		Conn:      conn,
		ConnID:    connID,
		isClosed:  false,

		ExitChan:   make(chan bool, 1),
		msgChan:    make(chan []byte),
		MsgHandler: msgHandler,
		property:   make(map[string]interface{}),
	}
	c.TcpServer.GetConnManager().Add(c)
	return c
}

func (c *Connection) StartReader() {
	fmt.Println("Reader GoRoutine is running...")
	defer fmt.Println("connID = ", c.ConnID, "Reader is exit, remote addr is", c.RemoteAddr().String())
	defer c.Stop()

	for {
		//buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		//_, err := c.Conn.Read(buf)
		//if err != nil {
		//	fmt.Println("recv buf err", err)
		//	continue
		//}

		dp := NewDataPack()
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head error")
			break
		}

		msg, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("unpack error", err)
			break
		}

		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error", err)
				break
			}

		}

		msg.SetData(data)
		//if err := c.handleAPI(c.Conn, buf, cnt); err != nil {
		//	fmt.Println("ConnID", c.ConnID, "handler is error", err)
		//	break
		//}

		req := Request{
			conn: c,
			msg:  msg,
		}

		if utils.GlobalObject.WorkerPoolSize > 0 {
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			go c.MsgHandler.DoMsgHandler(&req)
		}
	}
}

func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is running...]")
	defer fmt.Println(c.RemoteAddr().String(), "[conn writer exit]")

	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("send data error", err)
				return
			}
		case <-c.ExitChan:
			return
		}
	}
}

func (c *Connection) Start() {
	fmt.Println("Conn start.. ConnID =", c.ConnID)
	go c.StartReader()
	go c.StartWriter()

	c.TcpServer.CallOnConnStart(c)
}

func (c *Connection) Stop() {
	fmt.Println("Conn stop.. ConnID=", c.ConnID)
	if c.isClosed == true {
		return
	}
	c.isClosed = true
	c.TcpServer.CallOnConnStop(c)
	c.Conn.Close()
	c.ExitChan <- true
	c.TcpServer.GetConnManager().Remove(c)
	close(c.ExitChan)
	close(c.msgChan)
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed {
		return errors.New("connection closed")
	}
	dp := NewDataPack()
	binaryMsg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack error msg id ", msgId)
		return errors.New("Pack error msg")
	}
	//if _, err := c.Conn.Write(binaryMsg); err != nil {
	//	fmt.Println("Write msg id", msgId, "error", err)
	//	return errors.New("conn write error")
	//}
	c.msgChan <- binaryMsg
	return nil
}

func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	c.property[key] = value
}

func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()
	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("no property found")
	}
}

func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	delete(c.property, key)
}
