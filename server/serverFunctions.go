package server

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/Golds-l/goproxy/communication"
)

type CloudConnection struct {
	ConnLocal       *net.Conn
	ConnRemote      *net.TCPConn
	Id              string
	StartTime       int64
	Alive           bool
	CommunicateChan chan int
}

func (conn *CloudConnection) CloudServerToLocal() {
	cache := make([]byte, 4096)
	connLocal, connRemote := *conn.ConnLocal, *conn.ConnRemote
	for {
		select {
		case <-conn.CommunicateChan:
			return
		default:
			readNum, readErr := connRemote.Read(cache)
			if readErr != nil {
				continue
			}
			_, writeErr := connLocal.Write(cache[:readNum])
			if writeErr != nil {
				continue
			}
		}
	}
}

func (conn *CloudConnection) LocalToCloudServer() {
	cache := make([]byte, 4096)
	connLocal, connRemote := *conn.ConnLocal, *conn.ConnRemote
	for {
		readNum, readErr := connLocal.Read(cache)
		if readErr != nil {
			CloseCloudConnection(conn)
			conn.CommunicateChan <- 1
			return
		}
		_, writeErr := connRemote.Write(cache[:readNum])
		if writeErr != nil {
			CloseCloudConnection(conn)
			conn.CommunicateChan <- 1
			return
		}
	}
}

func (conn *CloudConnection) Close() error {
	if conn.ConnLocal != nil && conn.ConnRemote != nil {
		lErr := (*conn.ConnLocal).Close()
		rErr := (*conn.ConnRemote).Close()
		if lErr != nil || rErr != nil {
			return errors.New(lErr.Error() + "\n" + rErr.Error())
		}
	} else if conn.ConnLocal != nil && conn.ConnRemote == nil {
		lErr := (*conn.ConnLocal).Close()
		if lErr != nil {
			return lErr
		}
	} else if conn.ConnLocal == nil && conn.ConnRemote != nil {
		rErr := (*conn.ConnRemote).Close()
		if rErr != nil {
			return rErr
		}
	}
	return nil
}

func MakeNewConn(communicationConn *communication.CommunicationConnection, listener *net.TCPListener, conn *CloudConnection) error {
	// readCache := make([]byte, 256)
	conn.Id = communication.GenerateConnId()
	fmt.Printf("stop KeepAliveS ")
	sendErr := communicationConn.SendNewConnectionRequest(conn.Id) // make new Connection
	fmt.Printf("stopped\n")
	if sendErr != nil {
		return sendErr
	}
	fmt.Println("New request has sent, wait connections...")
	for n := 0; n < 8; n++ {
		deadlineErr := listener.SetDeadline(time.Now().Add(3 * time.Second))
		if deadlineErr != nil {
			return errors.New("listener error")
		}
		newConn, newConnectionErr := listener.AcceptTCP() // for loop to establish connection
		if newConnectionErr != nil {
			fmt.Printf("Connection etablished error. %v\n", newConnectionErr)
			// _ = newConn.Close()
			continue
		}
		if strings.Split(newConn.RemoteAddr().String(), ":")[0] != communicationConn.IP {
			fmt.Printf("unkonw connection accpeted:%v\n", newConn.RemoteAddr().String())
			_ = newConn.Close()
			continue
		} else {
			fmt.Println("Establish a connection with a remote client..")
			conn.ConnRemote = newConn
			conn.Alive = true
			conn.StartTime = time.Now().Unix()
			return nil
		}

	}
	return errors.New(fmt.Sprintf("conection accept times out, close all connections."))
}

func CloseCloudConnection(conn *CloudConnection) {
	err := conn.Close()
	if err != nil {
		fmt.Printf("Cloud connection close error! %v\n", err)
		conn.Alive = false
	} else {
		conn.Alive = false
		fmt.Printf("Connection closed. Id:%v Time:%v\n", conn.Id, time.Now().Format("2006-01-02 15:04:05"))
	}
}

func CheckAlive(conns []*CloudConnection) (int, []*CloudConnection) {
	var newConns = make([]*CloudConnection, 0, 15)
	for i := range conns {
		if conns[i].Alive {
			newConns = append(newConns, conns[i])
		} else {
			continue
		}
	}
	return len(newConns), newConns
}
