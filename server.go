package main

import (
	"fmt"
	"net"
	"os"

	"github.com/Golds-l/goproxy/communication"
	"github.com/Golds-l/goproxy/other"
	"github.com/Golds-l/goproxy/server"
)

func main() {
	var communicationConn *communication.Connection
	var connections []*server.CloudConnection
	var aliveNum int
	argsMap, ok := other.GetArgsCloudServer()
	if !ok { // need update
		fmt.Println("args error")
		os.Exit(0)
	}
	listenLocal, err := net.Listen("tcp", ":"+argsMap["localPort"])
	if err != nil {
		fmt.Println("listen err", err)
	}
	listenRemote, err := net.Listen("tcp", ":"+argsMap["remotePort"])
	if err != nil {
		fmt.Println("listen err", err)
	}
	fmt.Printf("begin listen... local port:%v remote port:%v\n", argsMap["localPort"], argsMap["remotePort"])
	communicationConn = communication.EstablishCommunicationConnS(listenRemote)
	go server.KeepAliveS(communicationConn, listenRemote)
	for {
		connLocal, connLocalErr := listenLocal.Accept()
		fmt.Printf("connection from %v\n", connLocal.RemoteAddr())
		if connLocalErr != nil {
			fmt.Printf("connection from %v error!\n.close and listening", connLocal.RemoteAddr())
			continue
		}
		conn, mkErr := server.MakeNewConn(communicationConn, listenRemote, connLocal)
		if mkErr != nil {
			_ = conn.Close()
			continue
		}
		go conn.LocalToCloudServer()
		go conn.CloudServerToLocal()
		fmt.Printf("connection etablished. id: %v\n", conn.Id)
		connections = append(connections, conn)
		aliveNum, connections = server.CheckAlive(connections)
		fmt.Printf("connections num: %v\n", aliveNum)
	}
}
