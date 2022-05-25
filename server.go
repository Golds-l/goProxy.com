package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/Golds-l/goproxy/communication"
	"github.com/Golds-l/goproxy/other"
	"github.com/Golds-l/goproxy/server"
)

func main() {
	var communicationConn = new(communication.Connection)
	var connections []*server.CloudConnection
	var aliveNum int
	argsMap, ok := other.GetArgsCloudServer()
	if !ok {
		fmt.Println("args error")
		os.Exit(0)
	}
	listenLocal, err := net.Listen("tcp", ":"+argsMap["localPort"])
	if err != nil {
		fmt.Println("listen error, please check the port.", err)
		os.Exit(0)
	}
	listenAddr := net.IPv4(0, 0, 0, 0)
	listenPort, _ := strconv.Atoi(argsMap["remotePort"])
	tcpAddr := net.TCPAddr{IP: listenAddr, Port: listenPort, Zone: ""}
	listenRemote, err := net.ListenTCP("tcp", &tcpAddr)
	if err != nil {
		fmt.Println("listen error, please check the port.", err)
		os.Exit(0)
	}
	fmt.Printf("Start listening. Local port:%v Remote port:%v\n", argsMap["localPort"], argsMap["remotePort"])
	fmt.Printf("time: %v\n", time.Now().Format("2006-01-02 15:04:05"))
	communication.EstablishCommunicationConnS(listenRemote, communicationConn)
	s := make(chan int) // MakeNewConn communicates with KeepAliveS
	go server.KeepAliveS(communicationConn, listenRemote, s)
	for {
		connLocal, connLocalErr := listenLocal.Accept()
		fmt.Printf("Connection from %v. %v ", connLocal.RemoteAddr(), time.Now().Format("2006-01-02 15:04:05"))
		if connLocalErr != nil {
			fmt.Printf("Connection from %v refused! %v\n", connLocal.RemoteAddr(), time.Now().Format("2006-01-02 15:04:05"))
			if connLocalErr != nil {
				fmt.Println(connLocalErr)
			} else {
				fmt.Println("unknow ip!", connLocal.RemoteAddr())
			}
			connLocal.Close()
			continue
		}
		fmt.Println("establish new connection")
		conn, mkErr := server.MakeNewConn(communicationConn, listenRemote, connLocal, s)
		if mkErr != nil {
			fmt.Println(mkErr, time.Now().Format("2006-01-02 15:04:05"))
			if conn != nil {
				_ = conn.Close()
				connLocal.Close()
				continue
			} else {
				connLocal.Close()
				continue
			}
		}
		q := make(chan int)
		go conn.CloudServerToLocal(q)
		go conn.LocalToCloudServer(q)
		fmt.Printf("Connection etablished. Id: %v Time:%v\n", conn.Id, time.Now().Format("2006-01-02 15:04:05"))
		connections = append(connections, conn)
		aliveNum, connections = server.CheckAlive(connections)
		fmt.Printf("Number of connections: %v\n", aliveNum)
	}
}
