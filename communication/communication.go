package communication

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/Golds-l/goproxy/other"
)

type Connection struct {
	Conn          *net.Conn
	Id            string
	Communication bool
	StartTime     int64
}

func (c *Connection) Write(s string) (int, error) {
	conn := c.Conn
	n, e := (*conn).Write([]byte(s))
	return n, e
}

func (c *Connection) Close() error {
	conn := c.Conn
	e := (*conn).Close()
	return e
}

func (c *Connection) Read(b []byte) (int, error) {
	conn := c.Conn
	n, e := (*conn).Read(b)
	return n, e
}

func CloudServerToLocal(CloudServerR, CloudServerL net.Conn) {
	defer other.CloseConn(CloudServerR, CloudServerL)
	for {
		cache := make([]byte, 10240)
		readNum, _ := CloudServerR.Read(cache)
		_, _ = CloudServerL.Write(cache[:readNum])
	}
}

func LocalToCloudServer(CloudServerR, CloudServerL net.Conn) {
	defer other.CloseConn(CloudServerR, CloudServerL)
	for {
		cache := make([]byte, 10240)
		readNum, _ := CloudServerL.Read(cache)
		_, _ = CloudServerR.Write(cache[:readNum])
	}
}

func RemoteClientToCloudServer(RemoteClient, SSHRemoteClient net.Conn) {
	defer other.CloseConn(RemoteClient, SSHRemoteClient)
	for {
		cache := make([]byte, 10240)
		readNum, _ := SSHRemoteClient.Read(cache)
		_, _ = RemoteClient.Write(cache[:readNum])
	}
}

func CloudServerToRemoteClient(RemoteClient, SSHRemoteClient net.Conn) {
	defer other.CloseConn(RemoteClient, SSHRemoteClient)
	for {
		cache := make([]byte, 10240)
		readNum, _ := RemoteClient.Read(cache)
		_, _ = SSHRemoteClient.Write(cache[:readNum])
	}
}

func WriteAlive(conn *net.Conn, s string) {
	for {
		_, _ = (*conn).Write([]byte(s))
		time.Sleep(3 * time.Second)
	}
}

func EstablishCommunicationConnS(serverListener net.Listener) *Connection {
	var communicationConn Connection
	connACK := make([]byte, 512)
	for {
		conn, acceptErr := serverListener.Accept()
		if acceptErr != nil {
			fmt.Println("Can not connect cloud server... Retry in a second")
			time.Sleep(1 * time.Second)
			continue
		}
		communicationConn.Id = other.GenerateConnId()
		mesg := "communication:" + communicationConn.Id + ":xy"
		_, writeErr := conn.Write([]byte(mesg))
		if writeErr != nil {
			fmt.Printf("connection write err! %v\n", writeErr)
			fmt.Printf("connection:%v will be closed\n", conn)
			_ = conn.Close()
		}
		n, readErr := conn.Read(connACK)
		if readErr != nil {
			fmt.Printf("connection read err! %v\n", writeErr)
			fmt.Printf("connection:%v will be closed\n", communicationConn.Id)
			_ = conn.Close()
		}
		mesgACKSlice := strings.Split(string(connACK[:n]), ":")
		if mesgACKSlice[0] == "RCReady" && mesgACKSlice[1] == communicationConn.Id {
			communicationConn.Conn = &conn
			communicationConn.Communication = true
			communicationConn.StartTime = time.Now().Unix()
			fmt.Println("cloud server<--->remote client is connected!")
			break
		}
		_ = conn.Close()
	}
	return &communicationConn
}

func EstablishCommunicationConnC(addr string) *Connection {
	var communicationConn Connection
	communicationConnACK := make([]byte, 512)
	for {
		conn, connErr := net.Dial("tcp", addr)
		if connErr != nil {
			fmt.Println("Can not connect cloud server... Retry in a second")
			time.Sleep(1 * time.Second)
			continue
		}
		n, readErr := conn.Read(communicationConnACK)
		if readErr != nil {
			fmt.Printf("coonection read err!%v\n", readErr)
			_ = conn.Close()
			fmt.Println("close and retry in a second")
			time.Sleep(1 * time.Second)
			continue
		}
		mesSlice := strings.Split(string(communicationConnACK[:n]), ":")
		if mesSlice[0] == "communication" && mesSlice[2] == "xy" {
			communicationConn.Conn = &conn
			communicationConn.Id = mesSlice[1]
			communicationConn.Communication = true
			communicationConn.StartTime = time.Now().Unix()
			_, writeErr := communicationConn.Write("RCReady:" + communicationConn.Id)
			if writeErr != nil {
				fmt.Printf("communication connection write error!%v\n", writeErr)
				_ = conn.Close()
				fmt.Println("close and retry in a second")
				time.Sleep(1 * time.Second)
				continue
			}
			fmt.Printf("connection %v established\n", communicationConn.Id)
			break
		}
		_ = conn.Close()
	}
	return &communicationConn
}
