package znet

import (
	"ZINX/utils"
	"ZINX/ziface"
	"fmt"
	"net"
)

type Server struct {
	Name      string
	IPVersion string
	IP        string
	Port      int
	//Router    ziface.IRouter
	MsgHandler  ziface.IMsgHandler
	ConnManager ziface.IConnManager
	OnConnStart func(conn ziface.IConnection)
	OnConnStop  func(conn ziface.IConnection)
}

//func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
//	fmt.Println("[Conn Handle] CallbackToClient...")
//	if _, err := conn.Write(data[:cnt]); err != nil {
//		fmt.Println("write back error", err)
//		return errors.New("CallBackToClient error")
//	}
//	return nil
//}

func (s *Server) Start() {
	fmt.Printf("[Zinx] Server Name %s, Listenner at IP %s, Port %d is starting", utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.TcpPort)
	fmt.Printf("[Zinx] Version %s ,MaxConn %d, MaxPacketSize %d\n", utils.GlobalObject.Version, utils.GlobalObject.MaxConn, utils.GlobalObject.MaxPackageSize)

	fmt.Printf("[Start] Server Listenner at IP: %s. Port %d, is starting\n", s.IP, s.Port)

	go func() {
		s.MsgHandler.StartWorkerPool()

		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error:", err)
		}

		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen ", s.IPVersion, "err,", err)
		}

		fmt.Println("start Zinx server success,", s.Name, "success listening...")
		var connID uint32
		connID = 1
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}

			if s.ConnManager.Len() >= utils.GlobalObject.MaxConn {
				fmt.Println("Too many Connections! MaxConn =", utils.GlobalObject.MaxConn)
				conn.Close()
				continue
			}

			dealConn := NewConnection(s, conn, connID, s.MsgHandler)
			connID++

			go dealConn.Start()
		}
	}()

}

func (s *Server) Stop() {
	//TODO 将一些服务器资源 状态 或者已经开辟的链接信息进行停止 或者回收
	fmt.Println("Server Stop")
	s.ConnManager.ClearConn()
}

func (s *Server) Serve() {
	s.Start()

	//TODO 做一些启动服务器后的额外业务
	select {}
}

func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
	fmt.Println("Add Router succ")
}

func (s *Server) GetConnManager() ziface.IConnManager {
	return s.ConnManager
}

/*初始化server*/

func (s *Server) SetOnConnStart(hookFunc func(connection ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

func (s *Server) SetOnConnStop(hookFunc func(connection ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("--> Call OnConnStart.")
		s.OnConnStart(conn)
	}
}

func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("--> Call OnConnStop.")
		s.OnConnStop(conn)
	}
}

func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:        utils.GlobalObject.Name,
		IPVersion:   "tcp4",
		IP:          utils.GlobalObject.Host,
		Port:        utils.GlobalObject.TcpPort,
		MsgHandler:  NewMsgHandler(),
		ConnManager: NewConnManager(),
	}
	return s
}
