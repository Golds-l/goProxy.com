package main

import (
	"fmt"
	"net"
	"os"

	"goProxy.com/client"
	"goProxy.com/communication"
	"goProxy.com/other"
)

func main() {
	argsMap, ok := other.GetArgsRemoteClient()
	fmt.Printf("server:%v ", argsMap["CloudServer"]+":"+argsMap["cloudServerPort"])
	fmt.Printf("host:%v\n", argsMap["remoteHost"]+":"+argsMap["remoteHostPort"])
	if !ok {
		fmt.Println("args illegal")
		os.Exit(0)
	}
	addrCloud := argsMap["CloudServer"] + ":" + argsMap["cloudServerPort"]
	addrLocal := argsMap["remoteHost"] + ":" + argsMap["remoteHostPort"]
	communicationConn, communicationConnErr := net.Dial("tcp", addrCloud)
	other.HandleErr(communicationConnErr)
	communicationConn.Write([]byte("RCReady"))
	fmt.Println("client ready")
	cache := make([]byte, 10240)
	for {
		n, readErr := communicationConn.Read(cache)
		if readErr != nil {
			fmt.Printf("%v\n", readErr)
			continue
		}
		if string(cache[:n]) == "NEWXX" {
			connCloud, connLocal := client.MakeNewClient(addrCloud, addrLocal)
			go communication.RemoteClientToCloudServer(connCloud, connLocal)
			go communication.CloudServerToRemoteClient(connCloud, connLocal)
		}
	}
}
