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
	StartTime       int64
	Alive           bool
	IP              string
	CommunicateChan chan int
	Listener        *net.TCPListener
	Port            string // local port
}

type CloudConnection struct {
	ConnLocal  *net.Conn
	ConnRemote *net.TCPConn
	Id         string
	StartTime  int64
	Alive      bool
}

type RemoteConnection struct {
	Conn          *net.TCPConn
	RedundantMesg string
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
	c.Alive = false
	c.StartTime = 0
	e := (*c.Conn).Close()
	(*c.Listener).Close()
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

func (communicationConn *CommunicationConnection) EstablishCommunicationConnS(conn net.Conn) {
	connACK := make([]byte, 512)
	// var isLog = false
	for {
		// _ = serverListener.SetDeadline(time.Now().Add(10 * time.Second))
		// conn, acceptErr := serverListener.Accept()
		// if acceptErr != nil {
		// 	if !isLog {
		// 		fmt.Println("Can not establish communication connections,retry in a second")
		// 		isLog = true
		// 	}
		// 	time.Sleep(1 * time.Second)
		// 	continue
		// } else {
		// 	isLog = false
		// }
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
			communicationConn.StartTime = time.Now().Unix()
			communicationConn.IP = strings.Split(conn.RemoteAddr().String(), ":")[0]
			fmt.Printf("cloud server<--->remote client is connected!\nfrom %v id:%v %v\n", conn.RemoteAddr().String(), communicationConn.Id, time.Now().Format("2006-01-02 15:04:05"))
			return
		}
		_ = conn.Close()
	}
}

func EstablishCommunicationConnS(conn net.Conn) (*CommunicationConnection, error) {
	var comConn CommunicationConnection
	comConn.Id = GenerateConnId()
	mesg := make([]byte, 1024)
	_, writeErr := conn.Write([]byte("communication:" + comConn.Id + ":xy"))
	if writeErr != nil {
		return nil, fmt.Errorf("communication connection establish error when write")
	}
	n, readErr := conn.Read(mesg)
	if readErr != nil {
		return nil, fmt.Errorf("communication connection establish error when read")
	}
	mesgSlice := strings.Split(string(mesg[:n]), ":")
	if len(mesgSlice) >= 3 && mesgSlice[0] == "RCReady" && mesgSlice[1] == comConn.Id && mesgSlice[2] == "wodexinxin" {
		comConn.Conn = &conn
		comConn.StartTime = time.Now().Unix()
		comConn.IP = strings.Split(conn.RemoteAddr().String(), ":")[0]
		comConn.CommunicateChan = make(chan int)
		comConn.Alive = true
		localListener, listenerErr := StartListener(mesgSlice[3])
		if listenerErr != nil {
			fmt.Println(listenerErr)
			return nil, listenerErr
		}
		comConn.Listener = localListener
		return &comConn, nil
	} else {
		return nil, fmt.Errorf("wrong messages when establish communication connection:%v", string(mesg[:n]))
	}
}

func (communicationConn *CommunicationConnection) EstablishNewConn(id string, connPool map[string]*RemoteConnection) (*RemoteConnection, error) {
	fmt.Printf("stop KeepAliveS ")
	sendErr := communicationConn.SendNewConnectionRequest(id) // make new Connection
	fmt.Printf("stopped\n")
	if sendErr != nil {
		return nil, sendErr
	}
	fmt.Printf("New request has sent, find connections...")
	for n := 0; n < 10; n++ {
		if conn, ok := connPool[id]; ok {
			delete(connPool, id)
			fmt.Println("found connections...")
			return conn, nil
		} else {
			time.Sleep(1 * time.Second)
		}
	}
	return nil, fmt.Errorf("times out")
}

func EstablishCommunicationConnC(addr string, localPort string) *CommunicationConnection {
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
		_, writeErr := conn.Write([]byte("comConn"))
		if writeErr != nil {
			fmt.Println(writeErr)
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
			communicationConn.Alive = true
			communicationConn.StartTime = time.Now().Unix()
			communicationConn.CommunicateChan = make(chan int)
			communicationConn.Port = localPort
			_, writeErr := communicationConn.Write([]byte("RCReady:" + communicationConn.Id + ":wodexinxin:" + localPort))
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

// goroutine
func (communicationConn *CommunicationConnection) KeepCommunicationConn() {
	mesg := make([]byte, 512)
	for {
		if !communicationConn.Alive {
			fmt.Printf("communication connection closed %v\n", communicationConn.IP)
			return
		}
		select {
		case <-communicationConn.CommunicateChan:
			time.Sleep(3 * time.Second) // Sleep for make new connection
		default:
			time.Sleep(1 * time.Second)
			_, writeErr := communicationConn.Write([]byte("isAlive"))
			if writeErr != nil {
				fmt.Printf("communication connection write err %v\n", writeErr)
				_ = communicationConn.Close()
				return
			}
			n, readErr := communicationConn.Read(mesg)
			if readErr != nil {
				fmt.Printf("communication connection read error %v\n", readErr)
				_ = communicationConn.Close()
				return
			}
			if string(mesg[:n]) == "alive" {
				communicationConn.Alive = true
				continue
			}
			if string(mesg[:n]) == "closexy" {
				_ = communicationConn.Close() // TODO:remove from comConn pool
				return
			}
		}
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

func StartListener(port string) (*net.TCPListener, error) {
	listenAddr := net.IPv4(0, 0, 0, 0)
	portInt, convErr := strconv.Atoi(port)
	if convErr != nil || portInt < 0 || portInt > 65535 {
		return nil, fmt.Errorf("port error %v", port)
	}
	tcpAddr := net.TCPAddr{IP: listenAddr, Port: portInt, Zone: ""}
	listenRemote, listenErr := net.ListenTCP("tcp", &tcpAddr)
	if listenErr != nil {
		return nil, fmt.Errorf("listen error:%v", listenErr)
	}
	return listenRemote, nil
}
