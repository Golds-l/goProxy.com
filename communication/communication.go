package communication

import (
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

type Connection struct {
	Conn          *net.Conn
	Id            string
	Communication bool
	StartTime     int64
	Alive         bool
}

func GenerateRandomInt() int64 {
	rand.Seed(time.Now().Unix())
	return rand.Int63()
}

func GenerateConnId() string {
	return strconv.FormatInt(GenerateRandomInt(), 10)
}

func (c *Connection) Write(s []byte) (int, error) {
	conn := c.Conn
	n, e := (*conn).Write(s)
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

func WriteAlive(conn *net.Conn, s string) {
	for {
		_, _ = (*conn).Write([]byte(s))
		time.Sleep(3 * time.Second)
	}
}

func EstablishCommunicationConnS(serverListener net.Listener, communicationConn *Connection) {
	connACK := make([]byte, 512)
	for {
		conn, acceptErr := serverListener.Accept()
		if acceptErr != nil {
			fmt.Println("Can not connect cloud server... Retry in a second")
			time.Sleep(1 * time.Second)
			continue
		}
		communicationConn.Id = GenerateConnId()
		mesg := "communication:" + communicationConn.Id + ":xy"
		_, writeErr := conn.Write([]byte(mesg))
		if writeErr != nil {
			fmt.Printf("connection write error! %v\n", writeErr)
			fmt.Printf("connection:%v will be closed\n", conn)
			_ = conn.Close()
		}
		n, readErr := conn.Read(connACK)
		if readErr != nil {
			fmt.Printf("connection read error! %v\n", writeErr)
			fmt.Printf("connection:%v will be closed\n", communicationConn.Id)
			_ = conn.Close()
		}
		mesgACKSlice := strings.Split(string(connACK[:n]), ":")
		if mesgACKSlice[0] == "RCReady" && mesgACKSlice[1] == communicationConn.Id && mesgACKSlice[2] == "wodexinxin" {
			communicationConn.Conn = &conn
			communicationConn.Communication = true
			communicationConn.StartTime = time.Now().Unix()
			fmt.Printf("cloud server<--->remote client is connected!\nid:%v\n", communicationConn.Id)
			break
		}
		_ = conn.Close()
	}
}

func EstablishCommunicationConnC(addr string) *Connection {
	var communicationConn Connection
	var isLog = false
	communicationConnACK := make([]byte, 512)
	for {
		conn, connErr := net.Dial("tcp", addr)
		if connErr != nil {
			if !isLog {
				fmt.Println("Can not connect cloud server... Retrying")
			}
			isLog = true
			time.Sleep(1 * time.Second)
			continue
		} else {
			isLog = false
		}
		n, readErr := conn.Read(communicationConnACK)
		if readErr != nil {
			fmt.Printf("coonection read error!%v\n", readErr)
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
			_, writeErr := communicationConn.Write([]byte("RCReady:" + communicationConn.Id + ":wodexinxin"))
			if writeErr != nil {
				fmt.Printf("communication connection write error!%v\n", writeErr)
				_ = conn.Close()
				fmt.Println("close and retry in a second")
				time.Sleep(1 * time.Second)
				continue
			}
			fmt.Printf("Connection established. Id: %v\n", communicationConn.Id)
			break
		}
		_ = conn.Close()
	}
	return &communicationConn
}
