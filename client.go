package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Golds-l/goproxy/client"
	"github.com/Golds-l/goproxy/communication"
	"github.com/Golds-l/goproxy/other"
)

func main() {
	var communicationConn *communication.Connection
	var connections []*communication.Connection
	argsMap, ok := other.GetArgsRemoteClient()
	fmt.Printf("server:%v ", argsMap["CloudServer"]+":"+argsMap["cloudServerPort"])
	fmt.Printf("host:%v\n", argsMap["remoteHost"]+":"+argsMap["remoteHostPort"])
	if !ok {
		fmt.Println("invalid option. Try '--help' for more information.")
		os.Exit(0)
	}
	addrCloud := argsMap["CloudServer"] + ":" + argsMap["cloudServerPort"]
	addrLocal := argsMap["remoteHost"] + ":" + argsMap["remoteHostPort"]
	cache := make([]byte, 10240)
	communicationConn = communication.EstablishCommunicationConnC(addrCloud)
	fmt.Println("client ready")
	for {
		n, readErr := communicationConn.Read(cache)
		if readErr != nil {
			fmt.Printf("communication connection read error. %v\n", readErr)
			fmt.Println("close and reconnect in a second..")
			_ = communicationConn.Close()
			time.Sleep(1 * time.Second)
			communicationConn = communication.EstablishCommunicationConnC(addrCloud)
			continue
		}
		mesg := string(cache[:n])
		if mesg == "isAlive" {
			//fmt.Printf("server alive.. %v\n", time.Now())
			_, writeErr := communicationConn.Write("alive")
			if writeErr != nil {
				fmt.Printf("communication connection write error. %v\n", writeErr)
				fmt.Println("close and reconnect in a second..")
				time.Sleep(1 * time.Second)
				communicationConn = communication.EstablishCommunicationConnC(addrCloud)
				continue
			}
		}
		mesgSlice := strings.Split(string(cache[:n]), ":")
		if mesgSlice[0] == "NEWC" {
			var conn communication.Connection
			conn.Id = mesgSlice[1]
			_, writeErr := communicationConn.Write("NEW:" + mesgSlice[1])
			if writeErr != nil {
				fmt.Printf("communication connection write error. %v\n", writeErr)
				fmt.Println("close and reconnect in a second..")
				time.Sleep(1 * time.Second)
				communicationConn = communication.EstablishCommunicationConnC(addrCloud)
				continue
			}
			fmt.Println("cloud connection established, connecting local end system")
			connCloud, connLocal := client.MakeNewClient(addrCloud, addrLocal)
			go communication.RemoteClientToCloudServer(connCloud, connLocal)
			go communication.CloudServerToRemoteClient(connCloud, connLocal)
			conn.Conn = &connLocal
			fmt.Println("local end system connection established")
			connections = append(connections, &conn)
		}
	}
}
