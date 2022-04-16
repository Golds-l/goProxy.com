package client

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/Golds-l/goproxy/communication"
)

type RemoteConnection struct {
	ConnCloud   *net.Conn
	ConnProcess *net.Conn
	Id          string
	StartTime   int64
	Alive       bool
}

func (conn *RemoteConnection) RemoteClientToCloudServer() {
	connCloud, connPro := *conn.ConnCloud, *conn.ConnProcess
	cache := make([]byte, 4096)
	for {
		readNum, connProReadErr := connPro.Read(cache)
		_, connCloudWriteErr := connCloud.Write(cache[:readNum])
		if connProReadErr != nil || connCloudWriteErr != nil {
			break
		}
	}
}

func (conn *RemoteConnection) CloudServerToRemoteClient() {
	defer CloseRemoteConnection(conn)
	connCloud, connPro := *conn.ConnCloud, *conn.ConnProcess
	cache := make([]byte, 4096)
	for {
		readNum, connProReadErr := connCloud.Read(cache)
		_, connCloudWriteErr := connPro.Write(cache[:readNum])
		if connProReadErr != nil || connCloudWriteErr != nil {
			break
		}
	}
}

func (conn *RemoteConnection) Close() error {
	connCloud, connPro := *conn.ConnCloud, *conn.ConnProcess
	connCloudCloseErr, connProCloseErr := connCloud.Close(), connPro.Close()
	if connCloudCloseErr != nil || connProCloseErr != nil {
		return errors.New(connCloudCloseErr.Error() + "/n" + connProCloseErr.Error())
	} else {
		return nil
	}
}

func MakeNewClient(serverAddr, localAddr, id string) (*RemoteConnection, error) {
	var conn RemoteConnection
	connServer, _ := net.Dial("tcp", serverAddr)
	connLocal, connLocalErr := net.Dial("tcp", localAddr)
	conn.Id = id
	conn.ConnCloud = &connServer
	conn.ConnProcess = &connLocal
	conn.StartTime = time.Now().Unix()
	conn.Alive = true
	return &conn, connLocalErr
}

func KeepAliveC(conn *communication.Connection, addr string) {
	cache := make([]byte, 512)
	for {
		n, readErr := conn.Read(cache)
		if readErr != nil {
			fmt.Printf("client communication connection read err. %v\n", readErr)
			fmt.Println("close and reconnect in a second..")
			_ = conn.Close()
			time.Sleep(1 * time.Second)
			conn = communication.EstablishCommunicationConnC(addr)
		}
		if string(cache[:n]) == "isAlive" {
			fmt.Printf("connection %v alive.. %v\n", conn.Id, time.Now())
		}
		_, writeErr := conn.Write([]byte("alive"))
		if writeErr != nil {
			fmt.Printf("client communication connection %v write err. %v\n", conn.Id, writeErr)
			fmt.Println("close and reconnect in a second..")
			time.Sleep(1 * time.Second)
			conn = communication.EstablishCommunicationConnC(addr)
		}
	}
}

func CloseRemoteConnection(conn *RemoteConnection) {
	err := conn.Close()
	if err != nil {
		fmt.Printf("remote connection close error! %v\n", err)
	} else {
		conn.Alive = false
		fmt.Printf("%v closed.\n", conn.Id)
	}
}
