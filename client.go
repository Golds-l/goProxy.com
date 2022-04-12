package main

import (
	"fmt"
	"github.com/Golds-l/goproxy/client"
	"github.com/Golds-l/goproxy/communication"
	"github.com/Golds-l/goproxy/other"
	"net"
	"os"
	"time"
)

func main() {
	var communicationConn *net.Conn
	argsMap, ok := other.GetArgsRemoteClient()
	fmt.Printf("server:%v ", argsMap["CloudServer"]+":"+argsMap["cloudServerPort"])
	fmt.Printf("host:%v\n", argsMap["remoteHost"]+":"+argsMap["remoteHostPort"])
	if !ok {
		fmt.Println("args illegal")
		os.Exit(0)
	}
	addrCloud := argsMap["CloudServer"] + ":" + argsMap["cloudServerPort"]
	addrLocal := argsMap["remoteHost"] + ":" + argsMap["remoteHostPort"]
	cache := make([]byte, 10240)
	communicationConn = communication.EstablishCommunicationConnC(addrCloud)
	fmt.Println("client ready")
	//go client.KeepAliveC(communicationConn, addrCloud)
	for {
		n, readErr := (*communicationConn).Read(cache)
		if readErr != nil {
			fmt.Printf("client communication connection read err. %v\n", readErr)
			fmt.Println("close and reconnect in a second..")
			_ = (*communicationConn).Close()
			time.Sleep(1 * time.Second)
			communicationConn = communication.EstablishCommunicationConnC(addrCloud)
			continue
		}
		if string(cache[:n]) == "isAlive" {
			//fmt.Printf("server alive.. %v\n", time.Now())
			_, writeErr := (*communicationConn).Write([]byte("alive"))
			if writeErr != nil {
				fmt.Printf("client communication connection write err. %v\n", writeErr)
				fmt.Println("close and reconnect in a second..")
				time.Sleep(1 * time.Second)
				communicationConn = communication.EstablishCommunicationConnC(addrCloud)
				continue
			}
		}
		if string(cache[:n]) == "NEWXX" {
			fmt.Println("cloud connected, connecting local end system")
			connCloud, connLocal := client.MakeNewClient(addrCloud, addrLocal)
			go communication.RemoteClientToCloudServer(connCloud, connLocal)
			go communication.CloudServerToRemoteClient(connCloud, connLocal)
			fmt.Println("connect build")
		}
	}
}
