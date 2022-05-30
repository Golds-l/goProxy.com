package communication

import (
	"errors"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

type CommunicationConnection struct {
	Conn            *net.Conn
	Id              string
	Communication   bool
	StartTime       int64
	Alive           bool
	IP              string
	CommunicateChan chan int
}

type CloudConnection struct {
	ConnLocal  *net.Conn
	ConnRemote *net.TCPConn
	Id         string
	StartTime  int64
	Alive      bool
}

func GenerateRandomInt() int64 {
	rand.Seed(time.Now().Unix())
	return rand.Int63()
}

func GenerateConnId() string {
	return strconv.FormatInt(GenerateRandomInt(), 10)
}

func (c *CommunicationConnection) Write(s []byte) (int, error) {
	conn := c.Conn
	n, e := (*conn).Write(s)
	return n, e
}

func (c *CommunicationConnection) Close() error {
	conn := c.Conn
	e := (*conn).Close()
	return e
}

func (c *CommunicationConnection) Read(b []byte) (int, error) {
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

func (communicationConn *CommunicationConnection) EstablishCommunicationConnS(serverListener *net.TCPListener) {
	connACK := make([]byte, 512)
	var isLog = false
	for {
		_ = serverListener.SetDeadline(time.Now().Add(10 * time.Second))
		conn, acceptErr := serverListener.Accept()
		if acceptErr != nil {
			if !isLog {
				fmt.Println("Can not establish communication connections,retry in a second")
				isLog = true
			}
			time.Sleep(1 * time.Second)
			continue
		} else {
			isLog = false
		}
		fmt.Printf("accept a communication connection from %v %v\n", conn.RemoteAddr().String(), time.Now().Format("2006-01-02 15:04:05"))
		communicationConn.Id = GenerateConnId()
		mesg := "communication:" + communicationConn.Id + ":xy"
		_, writeErr := conn.Write([]byte(mesg))
		if writeErr != nil {
			fmt.Printf("communication connection write error! %v\n", writeErr)
			fmt.Printf("connection:%v will be closed\n", conn)
			_ = conn.Close()
			continue
		}
		n, readErr := conn.Read(connACK)
		if readErr != nil {
			fmt.Printf("communication connection read error! %v\n", writeErr)
			fmt.Printf("connection:%v will be closed\n", communicationConn.Id)
			_ = conn.Close()
			continue
		}
		mesgACKSlice := strings.Split(string(connACK[:n]), ":")
		if mesgACKSlice[0] == "RCReady" && mesgACKSlice[1] == communicationConn.Id && mesgACKSlice[2] == "wodexinxin" {
			communicationConn.Conn = &conn
			communicationConn.Communication = true
			communicationConn.StartTime = time.Now().Unix()
			communicationConn.IP = strings.Split(conn.RemoteAddr().String(), ":")[0]
			fmt.Printf("cloud server<--->remote client is connected!\nfrom %v id:%v %v\n", conn.RemoteAddr().String(), communicationConn.Id, time.Now().Format("2006-01-02 15:04:05"))
			return
		}
		_ = conn.Close()
	}
}

func EstablishCommunicationConnC(addr string) *CommunicationConnection {
	var communicationConn CommunicationConnection
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
			fmt.Printf("communication connection established. Id: %v %v\n", communicationConn.Id, time.Now().Format("2006-01-02 15:04:05"))
			break
		}
		_ = conn.Close()
	}
	return &communicationConn
}

func (communicationConn *CommunicationConnection) KeepAliveS(listener *net.TCPListener) {
	cache := make([]byte, 1024)
	for {
		select {
		case <-communicationConn.CommunicateChan:
			time.Sleep(2 * time.Second) // Sleep for make new connection
		default:
			_, writeErr := communicationConn.Write([]byte("isAlive"))
			if writeErr != nil {
				fmt.Printf("communication connection write err %v\n", writeErr)
				fmt.Printf("close and reconnect a second later.%v\n", time.Now().Format("2006-01-02 15:04:05"))
				time.Sleep(1 * time.Second)
				_ = communicationConn.Close()
				communicationConn.EstablishCommunicationConnS(listener)
				continue
			}
			n, readErr := communicationConn.Read(cache)
			if readErr != nil {
				fmt.Printf("communication connection read error %v\n", readErr)
				fmt.Printf("close and reconnect a second later.%v\n", time.Now().Format("2006-01-02 15:04:05"))
				time.Sleep(1 * time.Second)
				_ = communicationConn.Close()
				communicationConn.EstablishCommunicationConnS(listener)
				continue
			}
			if string(cache[:n]) == "alive" {
				communicationConn.Alive = true
			}
		}
		time.Sleep(3 * time.Second)
	}
}

func (communicationConn *CommunicationConnection) SendNewConnectionRequest(id string) error {
	communicationConn.StopCheckAlive()
	_, writeErr := communicationConn.Write([]byte("NEWC:" + id)) // make new Connection
	if writeErr != nil {
		fmt.Println(writeErr)
		return errors.New("communication connection write error when send new connection request")
	}
	return nil
}

func (communicationConn *CommunicationConnection) StopCheckAlive() {
	communicationConn.CommunicateChan <- 1
}
