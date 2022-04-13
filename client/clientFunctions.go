package client

import (
	"fmt"
	"net"
	"time"

	"github.com/Golds-l/goproxy/communication"

	"github.com/Golds-l/goproxy/other"
)

func MakeNewClient(serverAddr, localAddr string) (net.Conn, net.Conn) {
	connServer, err := net.Dial("tcp", serverAddr)
	other.HandleErr(err)
	fmt.Println("make new client")
	other.HandleErr(err)
	connLocal, connLocalErr := net.Dial("tcp", localAddr)
	other.HandleErr(connLocalErr)
	return connServer, connLocal
}

func KeepAliveC(conn *communication.Connection, addr string) {
	cache := make([]byte, 512)
	for {
		n, readErr := conn.Read(cache)
		if readErr != nil {
			fmt.Printf("client communication connection read err. %v\n", readErr)
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
			fmt.Printf("client communication connection %v write err. %v\n", conn.Id, writeErr)
			fmt.Println("close and reconnect in a second..")
			time.Sleep(1 * time.Second)
			conn = communication.EstablishCommunicationConnC(addr)
		}
	}
}
