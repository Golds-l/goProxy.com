package client

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/Golds-l/goproxy/communication"
)

func MakeNewClient(serverAddr, localAddr string) (net.Conn, net.Conn) {
	connServer, _ := net.Dial("tcp", serverAddr)
	connLocal, connLocalErr := net.Dial("tcp", localAddr)
	if connLocalErr != nil {
		fmt.Println("can not etablish connection to the local end system process,please check the port.")
		os.Exit(0)
	}
	return connServer, connLocal
}

func KeepAliveC(conn *communication.Connection, addr string) {
	cache := make([]byte, 512)
	for {
		n, readErr := conn.Read(cache)
		if readErr != nil {
			fmt.Printf("communication connection read error. %v\n", readErr)
			fmt.Println("close and reconnect in a second..")
			_ = conn.Close()
			time.Sleep(1 * time.Second)
			conn = communication.EstablishCommunicationConnC(addr)
		}
		if string(cache[:n]) == "isAlive" {
			fmt.Printf("connection %v alive.. %v\n", conn.Id, time.Now())
		}
		_, writeErr := conn.Write("alive")
		if writeErr != nil {
			fmt.Printf("communication connection %v write error. %v\n", conn.Id, writeErr)
			fmt.Println("close and reconnect in a second..")
			time.Sleep(1 * time.Second)
			conn = communication.EstablishCommunicationConnC(addr)
		}
	}
}
