package main

import (
	"github.com/xin6211/goproxy/util"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/xin6211/goproxy/communication"
)

func main() {
	logErr := util.InitLog()
	if logErr != nil {
		log.Println(logErr)
		os.Exit(0)
	}
	argsMap, ok := util.GetArgsCloudServer()
	if !ok {
		log.Println(argsMap)
		log.Println("args error!")
		os.Exit(0)
	}
	listenAddr := net.IPv4(0, 0, 0, 0)
	listenPort, _ := strconv.Atoi(argsMap["remotePort"])
	tcpAddr := net.TCPAddr{IP: listenAddr, Port: listenPort, Zone: ""}
	listenRemote, err := net.ListenTCP("tcp", &tcpAddr)
	if err != nil {
		log.Println("listen error, please check the port.", err)
		os.Exit(0)
	}
	connPool := make(map[string]*communication.RemoteConnection)
	log.Printf("Start listening. Local port:%v Remote port:%v\n", argsMap["localPort"], argsMap["remotePort"])
	log.Printf("time: %v\n", time.Now().Format("2006-01-02 15:04:05"))
	go ListenRemotePort(*listenRemote, connPool)
	select {}
}
