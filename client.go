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
	var connections []*client.RemoteConnection
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
	for {
		n, readErr := communicationConn.Read(cache)
		if readErr != nil {
			fmt.Printf("client communication connection read err. %v\n", readErr)
			fmt.Println("close and reconnect in a second..")
			_ = communicationConn.Close()
			time.Sleep(1 * time.Second)
			communicationConn = communication.EstablishCommunicationConnC(addrCloud)
			continue
		}
		mesg := string(cache[:n])
		if mesg == "isAlive" {
			_, writeErr := communicationConn.Write([]byte("alive"))
			if writeErr != nil {
				fmt.Printf("client communication connection write err. %v\n", writeErr)
				fmt.Println("close and reconnect in a second..")
				time.Sleep(1 * time.Second)
				communicationConn = communication.EstablishCommunicationConnC(addrCloud)
				continue
			}
		}
		mesgSlice := strings.Split(string(cache[:n]), ":")
		if mesgSlice[0] == "NEWC" {
			_, writeErr := communicationConn.Write([]byte("NEW:" + mesgSlice[1]))
			if writeErr != nil {
				fmt.Printf("client communication connection write err. %v\n", writeErr)
				fmt.Println("close and reconnect in a second..")
				time.Sleep(1 * time.Second)
				communicationConn = communication.EstablishCommunicationConnC(addrCloud)
				continue
			}
			fmt.Println("cloud connection established, connecting local end system")
			conn, mkErr := client.MakeNewClient(addrCloud, addrLocal, mesgSlice[1])
			if mkErr != nil {
				fmt.Println("can not establish connection to the local end system process,please check the port.")
				os.Exit(0)
			}
			go conn.RemoteClientToCloudServer()
			go conn.CloudServerToRemoteClient()
			fmt.Println("local end system connection established")
			connections = append(connections, conn)
		}
	}
}
