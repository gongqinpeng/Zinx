package main

import (
	"ZINX/znet"
	"fmt"
	"io"
	"net"
	"time"
)

func main() {
	fmt.Println("client1 start...")
	time.Sleep(time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start err,exit", err)
		return
	}

	for {
		dp := znet.NewDataPack()
		binaryMsg, err := dp.Pack(znet.NewMsgPackage(1, []byte("Zinxv0.8 client1 Test Message")))
		if err != nil {
			fmt.Println("pack err", err)
			return
		}
		if _, err := conn.Write(binaryMsg); err != nil {
			fmt.Println("write error", err)
			return
		}

		binaryHead := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, binaryHead); err != nil {
			fmt.Println("read head error", err)
			break
		}

		msgHead, err := dp.UnPack(binaryHead)
		if err != nil {
			fmt.Println("client unpack msghead err", err)
			break
		}
		if msgHead.GetMsgLen() > 0 {
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(conn, msg.Data); err != nil {
				fmt.Println("read msg data error", err)
				break
			}

			fmt.Println("recv server msg :id =", msg.Id, "len = ", msg.GetMsgLen(), "data = ", string(msg.Data))

		}
		time.Sleep(time.Second)
	}
}
