package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/Golds-l/goproxy/communication"
	"github.com/Golds-l/goproxy/other"
)

func main() {
	var communicationConn *communication.CommunicationConnection
	var connections []*RemoteConnection
	var aliveNum int
	argsMap, ok := other.GetArgsRemoteClient()
	log.Printf("server:%v ", argsMap["CloudServer"]+":"+argsMap["cloudServerPort"])
	log.Printf("host:%v\n", argsMap["remoteHost"]+":"+argsMap["remoteHostPort"])
	if !ok {
		log.Println("args illegal, check and restart")
		os.Exit(0)
	}
	logErr := other.InitLog()
	if logErr != nil {
		log.Println(logErr)
		os.Exit(0)
	}
	addrCloud := argsMap["CloudServer"] + ":" + argsMap["cloudServerPort"]
	addrLocal := argsMap["remoteHost"] + ":" + argsMap["remoteHostPort"]
	cache := make([]byte, 10240)
	communicationConn = communication.EstablishCommunicationConnC(addrCloud, argsMap["localPort"])
	log.Println("connected to the server!")
	for {
		n, readErr := communicationConn.Read(cache)
		if readErr != nil {
			log.Printf("client communication connection read error. %v\n", readErr)
			log.Println("close and reconnect in a second..")
			_ = communicationConn.Close()
			time.Sleep(1 * time.Second)
			communicationConn = communication.EstablishCommunicationConnC(addrCloud, argsMap["localPort"])
			log.Println("reconnect successfully!")
			continue
		}
		mesg := string(cache[:n])
		if mesg == "isAlive" {
			_, writeErr := communicationConn.Write([]byte("alive"))
			if writeErr != nil {
				log.Printf("communication connection write error. %v\n", writeErr)
				log.Printf("close and reconnect a second later. %v\n", time.Now().Format("2006-01-02 15:04:05"))
				time.Sleep(1 * time.Second)
				communicationConn = communication.EstablishCommunicationConnC(addrCloud, argsMap["localPort"])
				continue
			}
			continue
		}
		mesgSlice := strings.Split(string(cache[:n]), ":")
		if mesgSlice[0] == "NEWC" {
			log.Println("receive new connection request, establish connection..", time.Now().Format("2006-01-02 15:04:05"))
			conn, mkErr := MakeNewClient(addrCloud, addrLocal, mesgSlice[1], argsMap["remoteHost"]+":"+argsMap["remoteHostPort"])
			if mkErr != nil {
				log.Println("can not establish connection.", mkErr, time.Now().Format("2006-01-02 15:04:05"))
				continue
			}
			q := make(chan int)
			go conn.RemoteClientToCloudServer(q)
			go conn.CloudServerToRemoteClient(q)
			log.Printf("connection established. Id:%v. Time:%v\n", conn.Id, time.Now().Format("2006-01-02 15:04:05"))
			connections = append(connections, conn)
			aliveNum, connections = CheckAlive(connections)

			log.Printf("Number of connections: %v\n", aliveNum)
			continue
		}
	}
}
