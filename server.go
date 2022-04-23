package main

import (
	"fmt"
	"net"
	"os"
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
	if !ok { // need update
		fmt.Println("args error")
		os.Exit(0)
	}
	listenLocal, err := net.Listen("tcp", ":"+argsMap["localPort"])
	if err != nil {
		fmt.Println("listen error, please check the port.", err)
		os.Exit(0)
	}
	listenRemote, err := net.Listen("tcp", ":"+argsMap["remotePort"])
	if err != nil {
		fmt.Println("listen error, please check the port.", err)
		os.Exit(0)
	}
	fmt.Printf("Start listening... Local client connection port:%v Remote client connection port:%v\n", argsMap["localPort"], argsMap["remotePort"])
	communication.EstablishCommunicationConnS(listenRemote, communicationConn)
	go server.KeepAliveS(communicationConn, listenRemote)
	for {
		connLocal, connLocalErr := listenLocal.Accept()
		fmt.Printf("Connection from %v\n", connLocal.RemoteAddr())
		if connLocalErr != nil {
			fmt.Printf("Connection from %v error! %v\n", connLocal.RemoteAddr(), time.Now().String())
			connLocal.Close()
			continue
		}
		conn, mkErr := server.MakeNewConn(communicationConn, listenRemote, connLocal)
		if mkErr != nil {
			if conn != nil {
				_ = conn.Close()
				continue
			} else {
				continue
			}
		}
		go conn.LocalToCloudServer()
		go conn.CloudServerToLocal()
		fmt.Printf("Connection etablished. id: %v\n", conn.Id)
		connections = append(connections, conn)
		aliveNum, connections = server.CheckAlive(connections)
		fmt.Printf("Number of connections: %v\n", aliveNum)
	}
}
