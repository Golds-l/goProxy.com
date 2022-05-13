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

func (conn *RemoteConnection) RemoteClientToCloudServer(q chan int) {
	connCloud, connPro := *conn.ConnCloud, *conn.ConnProcess
	cache := make([]byte, 1440)
	for {
		select {
		case <-q:
			return
		default:
			readNum, connProReadErr := connPro.Read(cache)
			if connProReadErr != nil {
				continue
			}
			_, connCloudWriteErr := connCloud.Write(cache[:readNum])
			if connCloudWriteErr != nil {
				continue
			}
		}
	}
}

func (conn *RemoteConnection) CloudServerToRemoteClient(q chan int) {
	connCloud, connPro := *conn.ConnCloud, *conn.ConnProcess
	cache := make([]byte, 1440)
	for {
		readNum, connProReadErr := connCloud.Read(cache)
		if string(cache[:readNum]) == "XYEOF" {
			connCloud.Write(cache[:readNum])
			q <- 1
			CloseRemoteConnection(conn)
			return
		}
		if connProReadErr != nil {
			q <- 1
			CloseRemoteConnection(conn)
			return
		}
		_, connCloudWriteErr := connPro.Write(cache[:readNum])
		if connCloudWriteErr != nil {
			q <- 1
			CloseRemoteConnection(conn)
			return
		}
	}
}

func (conn *RemoteConnection) Close() error {
	connCloud, connPro := *conn.ConnCloud, *conn.ConnProcess
	connCloudCloseErr, connProCloseErr := connCloud.Close(), connPro.Close()
	if connCloudCloseErr != nil || connProCloseErr != nil {
		if connCloudCloseErr == nil {
			return connProCloseErr
		} else if connProCloseErr == nil {
			return connCloudCloseErr
		} else {
			return errors.New(connCloudCloseErr.Error() + "\n" + connProCloseErr.Error())
		}
	} else {
		return nil
	}
}

func MakeNewClient(serverAddr, localAddr, id string) (*RemoteConnection, error) {
	var conn RemoteConnection
	connServer, connServerErr := net.Dial("tcp", serverAddr)
	if connServerErr != nil {
		return nil, connServerErr
	}
	connLocal, connLocalErr := net.Dial("tcp", localAddr)
	if connLocalErr != nil {
		return nil, connLocalErr
	}
	_, serverWriteErr := connServer.Write([]byte(id + ":xy"))
	if serverWriteErr != nil {
		return nil, serverWriteErr
	}
	conn.Id = id
	conn.ConnCloud = &connServer
	conn.ConnProcess = &connLocal
	conn.StartTime = time.Now().Unix()
	conn.Alive = true
	return &conn, nil
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
		fmt.Printf("remote connection close error! Id:%v\n%v\n", conn.Id, err)
	} else {
		conn.Alive = false
		fmt.Printf("Connection closed. Id:%v Time:%v\n", conn.Id, time.Now().Format("2006-01-02 15:04:05"))
	}
}
