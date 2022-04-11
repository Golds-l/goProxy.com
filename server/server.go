package server

import (
	"fmt"
	"net"
	"os"

	"github.com/Golds-l/goproxy/communication"
	"github.com/Golds-l/goproxy/other"
)

func main() {
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
	communicationConn, communicationConnErr := listenRemote.Accept()
	other.HandleErr(communicationConnErr)
	ACKCache := make([]byte, 1024)
	n, _ := communicationConn.Read(ACKCache)
	if string(ACKCache[:n]) == "RCReady" {
		fmt.Println("cloud server --- remote client is connected!")
	}
	for {
		connLocal, connLocalErr := listenLocal.Accept()
		fmt.Printf("connect from %v\n", connLocal.RemoteAddr())
		other.HandleErr(connLocalErr)
		connRemote := MakeNewConn(communicationConn, listenRemote)
		go communication.CloudServerToLocal(connRemote, connLocal)
		go communication.LocalToCloudServer(connRemote, connLocal)
	}
}
