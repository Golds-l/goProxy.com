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
	ConnLocal  *net.Conn
	ConnRemote *net.Conn
	Id         string
	StartTime  int64
	Alive      bool
}

func (conn *CloudConnection) CloudServerToLocal(q chan int) {
	cache := make([]byte, 1440)
	connLocal, connRemote := *conn.ConnLocal, *conn.ConnRemote
	for {
		select {
		case <-q:
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

func (conn *CloudConnection) LocalToCloudServer(q chan int) {
	cache := make([]byte, 1440)
	connLocal, connRemote := *conn.ConnLocal, *conn.ConnRemote
	for {
		readNum, readErr := connLocal.Read(cache)
		if readErr != nil {
			CloseCloudConnection(conn)
			q <- 1
			return
		}
		_, writeErr := connRemote.Write(cache[:readNum])
		if writeErr != nil {
			CloseCloudConnection(conn)
			q <- 1
			return
		}
	}
}

func (conn *CloudConnection) Close() error {
	connLocal, connRemote := *conn.ConnLocal, *conn.ConnRemote
	connLocalCloseErr := connLocal.Close()
	connRemoteCloseErr := connRemote.Close()
	if connRemoteCloseErr != nil || connLocalCloseErr != nil {
		if connRemoteCloseErr == nil {
			return connLocalCloseErr
		} else if connLocalCloseErr == nil {
			return connRemoteCloseErr
		} else {
			return errors.New(connLocalCloseErr.Error() + "\n" + connRemoteCloseErr.Error())
		}
	} else {
		return nil
	}
}

func MakeNewConn(communicationConn *communication.Connection, listener net.Listener, connLocal net.Conn) (*CloudConnection, error) {
	readCache := make([]byte, 256)
	var conn CloudConnection
	conn.Id = communication.GenerateConnId()
	_, writeErr := communicationConn.Write([]byte("NEWC:" + conn.Id)) // make new Connection
	if writeErr != nil {
		fmt.Println(writeErr)
		return nil, errors.New("communication connection write error")
	}
	n, communicationReadErr := communicationConn.Read(readCache)
	if communicationReadErr != nil {
		fmt.Println(writeErr)
		return nil, errors.New("communication connection write error")
	}
	mesgSlice := strings.Split(string(readCache[:n]), ":")
	if mesgSlice[0] == "NEW" && mesgSlice[1] == conn.Id {
		ack := make([]byte, 1024)
		newConn, newConnectionErr := listener.Accept()
		if newConnectionErr != nil {
			fmt.Printf("connection etablished error. %v\n", newConnectionErr)
			_ = newConn.Close()
			return nil, newConnectionErr
		}
		n, readErr := newConn.Read(ack)
		if readErr != nil {
			_ = newConn.Close()
			return nil, errors.New("read error, connection establist failed.")
		}
		mesgStr := string(ack[:n])
		if mesgStr == conn.Id+":xy" {
			fmt.Println("Establish a connection with a remote client..")
			conn.ConnLocal = &connLocal
			conn.ConnRemote = &newConn
			conn.Alive = true
			conn.StartTime = time.Now().Unix()
			return &conn, nil
		} else {
			_ = newConn.Close()
			fmt.Printf("unknow ip, refused! %v\n", newConn.RemoteAddr().String())
			return nil, errors.New("unkonw ip")
		}
	} else {
		fmt.Println(mesgSlice, "can not establish with remote client. wrong mesg")
		return nil, errors.New("remote client error")
	}
}

func KeepAliveS(conn *communication.Connection, listener net.Listener) {
	cache := make([]byte, 1024)
	for {
		_, writeErr := conn.Write([]byte("isAlive"))
		if writeErr != nil {
			fmt.Printf("communication connection write err %v\n", writeErr)
			fmt.Println("close and reconnect in a second...")
			time.Sleep(1 * time.Second)
			_ = conn.Close()
			communication.EstablishCommunicationConnS(listener, conn)
			continue
		}
		n, readErr := conn.Read(cache)
		if readErr != nil {
			fmt.Printf("server communication connection read error %v\n", readErr)
			fmt.Println("close and reconnecting..")
			time.Sleep(1 * time.Second)
			_ = conn.Close()
			communication.EstablishCommunicationConnS(listener, conn)
			continue
		}
		if string(cache[:n]) == "alive" {
			conn.Alive = true
		}
		time.Sleep(3 * time.Second)
	}
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
