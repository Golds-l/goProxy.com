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
	ConnRemote *net.TCPConn
	Id         string
	StartTime  int64
	Alive      bool
}

func (conn *CloudConnection) CloudServerToLocal(q chan int) {
	cache := make([]byte, 4096)
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
	cache := make([]byte, 4096)
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

func MakeNewConn(communicationConn *communication.Connection, listener *net.TCPListener, connLocal net.Conn, q chan int) (*CloudConnection, error) {
	readCache := make([]byte, 256)
	var conn CloudConnection
	conn.Id = communication.GenerateConnId()
	fmt.Printf("stop KeepAliveS ")
	q <- 1 // clear communication connection
	fmt.Printf("stopped\n")
	_, writeErr := communicationConn.Write([]byte("NEWC:" + conn.Id)) // make new Connection
	if writeErr != nil {
		fmt.Println(writeErr)
		return nil, errors.New("communication connection write error when establish new connection")
	}
	fmt.Println("New request has sent...")
	n, communicationReadErr := communicationConn.Read(readCache)
	if communicationReadErr != nil {
		fmt.Println(writeErr)
		return nil, errors.New("communication connection read error when establish new connection")
	}
	mesgSlice := strings.Split(string(readCache[:n]), ":")
	if mesgSlice[0] == "NEW" && mesgSlice[1] == conn.Id {
		ack := make([]byte, 1024)
		fmt.Println("begin shakehand...")
		for i := 0; i < 8; i++ { // accept 10 connections
			deadlineErr := listener.SetDeadline(time.Now().Add(10 * time.Second))
			if deadlineErr != nil {
				return nil, errors.New("listener error")
			}
			newConn, newConnectionErr := listener.AcceptTCP() // for loop to establish connection
			if newConnectionErr != nil {
				fmt.Printf("Connection etablished error. %v\n", newConnectionErr)
				continue
			}
			n, readErr := newConn.Read(ack)
			if readErr != nil {
				fmt.Printf("New connection read error.From %v\n", newConn.RemoteAddr().String())
				_ = newConn.Close()
				continue
			}
			fmt.Printf("Received a connection from %v\n", newConn.RemoteAddr().String())
			mesgStr := string(ack[:n])
			if mesgStr == conn.Id+":xy" {
				fmt.Println("Establish a connection with a remote client..")
				_, newConnWriteErr := newConn.Write([]byte(conn.Id + ":xy" + ":wode")) // return a mesg for establish ssh server
				if newConnWriteErr != nil {
					return nil, errors.New(fmt.Sprintf("New connection write when send mesg"))
				}
				conn.ConnLocal = &connLocal
				conn.ConnRemote = newConn
				conn.Alive = true
				conn.StartTime = time.Now().Unix()
				return &conn, nil
			} else {
				fmt.Printf("wrong mseg:%v.from:%v\n", mesgStr, newConn.RemoteAddr().String())
				_ = newConn.Close()
				continue
			}
		}
	} else {
		fmt.Println(mesgSlice, "Can not establish with remote client. wrong mesg")
		return nil, errors.New("remote client error")
	}
	return nil, errors.New(fmt.Sprintf("conection accept times out, close all connections."))
}

func KeepAliveS(conn *communication.Connection, listener *net.TCPListener, q chan int) {
	cache := make([]byte, 1024)
	for {
		select {
		case <-q:
			time.Sleep(2 * time.Second) // Sleep for make new connection
		default:
			_, writeErr := conn.Write([]byte("isAlive"))
			if writeErr != nil {
				fmt.Printf("communication connection write err %v\n", writeErr)
				fmt.Printf("close and reconnect a second later.%v\n", time.Now().Format("2006-01-02 15:04:05"))
				time.Sleep(1 * time.Second)
				_ = conn.Close()
				communication.EstablishCommunicationConnS(listener, conn)
				continue
			}
			n, readErr := conn.Read(cache)
			if readErr != nil {
				fmt.Printf("communication connection read error %v\n", readErr)
				fmt.Printf("close and reconnect a second later.%v\n", time.Now().Format("2006-01-02 15:04:05"))
				time.Sleep(1 * time.Second)
				_ = conn.Close()
				communication.EstablishCommunicationConnS(listener, conn)
				continue
			}
			if string(cache[:n]) == "alive" {
				conn.Alive = true
			}
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
