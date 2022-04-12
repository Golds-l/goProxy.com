package main

import (
	"fmt"
	"github.com/Golds-l/goproxy/communication"
	"github.com/Golds-l/goproxy/other"
	"github.com/Golds-l/goproxy/server"
	"net"
	"os"
)

func main() {
	var communicationConn *net.Conn
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
		fmt.Printf("connect from %v\n", connLocal.RemoteAddr())
		other.HandleErr(connLocalErr)
		connRemote := server.MakeNewConn(communicationConn, listenRemote)
		go communication.CloudServerToLocal(connRemote, connLocal)
		go communication.LocalToCloudServer(connRemote, connLocal)
	}
}
