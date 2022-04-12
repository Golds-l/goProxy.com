package client

import (
	"fmt"
	"github.com/Golds-l/goproxy/communication"
	"net"
	"time"

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

func KeepAliveC(conn *net.Conn, addr string) {
	cache := make([]byte, 512)
	//go communication.WriteAlive(conn, "alive")
	for {
		n, readErr := (*conn).Read(cache)
		fmt.Println("client read")
		if readErr != nil {
			fmt.Printf("client communication connection read err. %v\n", readErr)
			fmt.Println("close and reconnect in a second..")
			_ = (*conn).Close()
			time.Sleep(1 * time.Second)
			conn = communication.EstablishCommunicationConnC(addr)
		}
		if string(cache[:n]) == "isAlive" {
			fmt.Printf("server alive.. %v\n", time.Now())
		}
		_, writeErr := (*conn).Write([]byte("alive"))
		fmt.Println("client write")
		if writeErr != nil {
			fmt.Printf("client communication connection write err. %v\n", writeErr)
			fmt.Println("close and reconnect in a second..")
			time.Sleep(1 * time.Second)
			conn = communication.EstablishCommunicationConnC(addr)
		}
	}
}
