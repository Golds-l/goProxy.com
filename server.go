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
	var connections []*communication.Connection
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
		other.HandleErr(connLocalErr)
		connRemote, mkErr := server.MakeNewConn(communicationConn, listenRemote)
		if mkErr != nil {
			connLocal.Close()
			continue
		}
		go communication.CloudServerToLocal(*connRemote.Conn, connLocal)
		go communication.LocalToCloudServer(*connRemote.Conn, connLocal)
		fmt.Printf("connection %v etablished.\n", connRemote.Id)
		connections = append(connections, connRemote)
		fmt.Println(connections)
	}
}
