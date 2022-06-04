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
	argsMap, ok := other.GetArgsCloudServer()
	if !ok {
		fmt.Println(argsMap)
		fmt.Println("args error!")
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
	connPool := make(map[string]*communication.RemoteConnection)
	fmt.Printf("Start listening. Local port:%v Remote port:%v\n", argsMap["localPort"], argsMap["remotePort"])
	fmt.Printf("time: %v\n", time.Now().Format("2006-01-02 15:04:05"))
	go server.ListenRemotePort(*listenRemote, connPool)
	for {

	}
}
