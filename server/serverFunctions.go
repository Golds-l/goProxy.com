package server

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/Golds-l/goproxy/communication"
	"github.com/Golds-l/goproxy/log"
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
	cache := make([]byte, 5120)
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
	cache := make([]byte, 5120)
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
	return fmt.Errorf("conection accept times out, close all connections")
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

// establish connection with nodes,make communication connction. goroutine
func ListenRemotePort(remoteListener net.TCPListener, connPool map[string]*communication.RemoteConnection) {
	// TODO: not only listen communication connection
	mesg := make([]byte, 1024)
	for {
		newConn, listenErr := remoteListener.AcceptTCP() // accept a conn from a node(maybe legal)
		if listenErr != nil {
			if newConn != nil {
				fmt.Printf("remote listener error, close connection.%v", listenErr)
				_ = newConn.Close()
				continue
			} else {
				fmt.Printf("remote listener error.%v", listenErr)
				continue
			}
		}
		fmt.Printf("remote port accpeted a connection from:%v", newConn.RemoteAddr())
		log.LogNow()
		newConn.SetReadDeadline(time.Now().Add(time.Second))
		n, readErr := newConn.Read(mesg)
		if readErr != nil {
			fmt.Printf("error when receivce mesg:%v\n", readErr)
			_ = newConn.Close()
			continue
		}
		newConn.SetReadDeadline(time.Time{})
		mesgSlice := strings.Split(string(mesg[:n]), ":")
		if mesgSlice[0] == "comConn" {
			comConn, mkErr := communication.EstablishCommunicationConnS(newConn) // try to establish communication conn
			if mkErr != nil {
				fmt.Println(mkErr)
				log.LogNow()
				_ = newConn.Close()
				continue
			}
			go comConn.KeepCommunicationConn()
			go ListenLocalPort(comConn, connPool)
			fmt.Printf("Cloud<----->Remote.communication conection established.from:%v", comConn.IP)
			log.LogNow()
		} else if mesgSlice[0] == "conn" {
			var RConn communication.RemoteConnection
			RConn.Conn = newConn
			if len(mesgSlice) > 2 {
				var redundantMesg string
				for i := range mesgSlice[2:] {
					redundantMesg += mesgSlice[2:][i]
				}
				RConn.RedundantMesg = redundantMesg
			}
			connPool[mesgSlice[1]] = &RConn
		} else {
			fmt.Printf("illegal connection from:%v", newConn.RemoteAddr())
			log.LogNow()
			_ = newConn.Close()
		}
	}
}

// accept conn from local client, establish it with remote client. goroutine
func ListenLocalPort(comConn *communication.CommunicationConnection, connPool map[string]*communication.RemoteConnection) {
	for {
		if !comConn.Alive {
			return
		}
		connLocal, listenErr := comConn.Listener.Accept()
		if listenErr != nil {
			if connLocal != nil {
				_ = connLocal.Close()
			}
			fmt.Printf("local listener error %v\n", listenErr)
			continue
		}
		fmt.Printf("port %v receive a request from:%v", connLocal.LocalAddr(), connLocal.RemoteAddr())
		log.LogNow()
		var newConn CloudConnection
		newConn.Id = communication.GenerateConnId()
		connRemote, rErr := comConn.EstablishNewConn(newConn.Id, connPool)
		if rErr != nil {
			fmt.Printf("establish error when connect remote:%v", rErr)
			connLocal.Close()
			log.LogNow()
			continue
		}
		newConn.ConnLocal = &connLocal
		newConn.ConnRemote = connRemote.Conn
		newConn.StartTime = time.Now().Unix()
		newConn.Alive = true
		newConn.CommunicateChan = make(chan int)
		if connRemote.RedundantMesg != "" {
			(*newConn.ConnLocal).Write([]byte(connRemote.RedundantMesg))
		}
		go newConn.LocalToCloudServer()
		go newConn.CloudServerToLocal()
		fmt.Printf("connection established.Id:%v", newConn.Id)
		log.LogNow()
	}
}
