package main

import (
	"errors"
	"log"
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
	cache := make([]byte, 5120)
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
	cache := make([]byte, 5120)
	for {
		readNum, connProReadErr := connCloud.Read(cache)
		if connProReadErr != nil {
			CloseRemoteConnection(conn)
			q <- 1
			return
		}
		_, connCloudWriteErr := connPro.Write(cache[:readNum])
		if connCloudWriteErr != nil {
			CloseRemoteConnection(conn)
			q <- 1
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

func MakeNewClient(serverAddr, localAddr, id string, host string) (*RemoteConnection, error) {
	var conn RemoteConnection
	for i := 0; i < 5; i++ {
		log.Printf("try to connect...")
		connServer, connServerErr := net.DialTimeout("tcp", serverAddr, 3*time.Second)
		if connServerErr != nil {
			log.Println("connection establish error", connServerErr)
			continue
		}
		connServer.Write([]byte("conn:" + id + ":"))
		log.Println("connections establish,connect local service", time.Now().Format("2006-01-02 15:04:05"))
		connLocal, connLocalErr := net.Dial("tcp", localAddr) // connect ssh server
		if connLocalErr != nil {
			_ = connServer.Close()
			return nil, connLocalErr
		}
		conn.Id = id
		conn.ConnCloud = &connServer
		conn.ConnProcess = &connLocal
		conn.StartTime = time.Now().Unix()
		conn.Alive = true
		return &conn, nil
	}
	return nil, errors.New("try out")
}

func KeepAliveC(conn *communication.CommunicationConnection, addr string) {
	cache := make([]byte, 512)
	for {
		n, readErr := conn.Read(cache)
		if readErr != nil {
			log.Printf("communication connection read err. %v\n", readErr)
			log.Println("close and reconnect in a second..")
			_ = conn.Close()
			time.Sleep(1 * time.Second)
			conn = communication.EstablishCommunicationConnC(addr, conn.Port)
		}
		if string(cache[:n]) == "isAlive" {
			log.Printf("connection %v alive.. %v\n", conn.Id, time.Now())
		}
		_, writeErr := conn.Write([]byte("alive"))
		if writeErr != nil {
			log.Printf("client communication connection %v write err. %v\n", conn.Id, writeErr)
			log.Println("close and reconnect in a second..")
			time.Sleep(1 * time.Second)
			conn = communication.EstablishCommunicationConnC(addr, conn.Port)
		}
	}
}

func CloseRemoteConnection(conn *RemoteConnection) {
	err := conn.Close()
	if err != nil {
		log.Printf("remote connection close error! Id:%v\n%v\n", conn.Id, err)
	} else {
		conn.Alive = false
		log.Printf("Connection closed. Id:%v Time:%v\n", conn.Id, time.Now().Format("2006-01-02 15:04:05"))
	}
}

func CheckAlive(conns []*RemoteConnection) (int, []*RemoteConnection) {
	var newConns = make([]*RemoteConnection, 0, 15)
	for i := range conns {
		if conns[i].Alive {
			newConns = append(newConns, conns[i])
		} else {
			continue
		}
	}
	return len(newConns), newConns
}
