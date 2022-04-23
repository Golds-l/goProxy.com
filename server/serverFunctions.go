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

func (conn *CloudConnection) CloudServerToLocal() {
	cache := make([]byte, 4096)
	connLocal, connRemote := *conn.ConnLocal, *conn.ConnRemote
	for {
		readNum, readErr := connRemote.Read(cache)
		_, writeErr := connLocal.Write(cache[:readNum])
		if writeErr != nil || readErr != nil {
			break
		}
	}
}

func (conn *CloudConnection) LocalToCloudServer() {
	defer CloseCloudConnection(conn)
	cache := make([]byte, 4096)
	connLocal, connRemote := *conn.ConnLocal, *conn.ConnRemote
	for {
		readNum, readErr := connLocal.Read(cache)
		_, writeErr := connRemote.Write(cache[:readNum])
		if writeErr != nil || readErr != nil {
			break
		}
	}
}

func (conn *CloudConnection) Close() error {
	connLocal, connRemote := *conn.ConnLocal, *conn.ConnRemote
	connLocalCloseErr := connLocal.Close()
	connRemoteCloseErr := connRemote.Close()
	if connRemoteCloseErr != nil || connLocalCloseErr != nil {
		return errors.New(connLocalCloseErr.Error() + "\n" + connRemoteCloseErr.Error())
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
	}
	n, _ := communicationConn.Read(readCache)
	mesgSlice := strings.Split(string(readCache[:n]), ":")
	if mesgSlice[0] == "NEW" && mesgSlice[1] == conn.Id {
		newConn, newConnectionErr := listener.Accept()
		if newConnectionErr != nil {
			fmt.Printf("connection etablished error. %v\n", newConnectionErr)
			return nil, newConnectionErr
		}
		fmt.Println("Establish a connection with a remote client..")
		conn.ConnLocal = &connLocal
		conn.ConnRemote = &newConn
		conn.Alive = true
		conn.StartTime = time.Now().Unix()
		return &conn, nil
	} else {
		fmt.Println(mesgSlice, "can not establish with remote client.")
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
	} else {
		conn.Alive = false
		fmt.Printf("Connection closed. Id: %v\n", conn.Id)
	}
}

func CheckAlive(conns []*CloudConnection) (int, []*CloudConnection) {
	var newConns []*CloudConnection
	for i := range conns {
		if conns[i].Alive {
			newConns = append(newConns, conns[i])
		} else {
			continue
		}
	}
	return len(newConns), newConns
}
