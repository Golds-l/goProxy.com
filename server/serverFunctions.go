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
	_, _ = communicationConn.Write([]byte("NEWC:" + conn.Id)) // make new connection
	n, _ := communicationConn.Read(readCache)
	mesgSlice := strings.Split(string(readCache[:n]), ":")
	if mesgSlice[0] == "NEW" && mesgSlice[1] == conn.Id {
		newConn, newConnectionErr := listener.Accept()
		if newConnectionErr != nil {
			fmt.Printf("connection etablished error. %v\n", newConnectionErr)
			return nil, newConnectionErr
		}
		fmt.Println("establishing remote..")
		conn.ConnLocal = &connLocal
		conn.ConnRemote = &newConn
		conn.Alive = true
		conn.StartTime = time.Now().Unix()
	}
	return &conn, nil
}

func KeepAliveS(conn *communication.Connection, listener net.Listener) {
	cache := make([]byte, 1024)
	for {
		_, writeErr := conn.Write([]byte("isAlive"))
		if writeErr != nil {
			fmt.Printf("server communication connection write err %v\n", writeErr)
			fmt.Println("close and reconnect..")
			time.Sleep(1 * time.Second)
			_ = conn.Close()
			conn = communication.EstablishCommunicationConnS(listener)
			continue
		}
		n, readErr := conn.Read(cache)
		if readErr != nil {
			fmt.Printf("server communication connection read err %v\n", readErr)
			fmt.Println("close and reconnect in a second..")
			time.Sleep(1 * time.Second)
			_ = conn.Close()
			conn = communication.EstablishCommunicationConnS(listener)
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
		fmt.Printf("cloud connection close error! %v\n", err)
	} else {
		conn.Alive = false
		fmt.Printf("%v closed.\n", conn.Id)
	}
}

func CheckAlive(conns []*CloudConnection) int {
	var num int
	var length = len(conns)
	fmt.Println(length)
	// for i := range conns {
	// 	if i == length-1 && !conns[i].Alive {
	// 		conns = conns[:i]
	// 	}
	// 	if i < length-1 && !conns[i].Alive {
	// 		conns = append(conns[:i], conns[i+1:]...)
	// 	} else {
	// 		num += 1
	// 	}
	// }
	return num
}
