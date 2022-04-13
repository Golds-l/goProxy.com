package server

import (
	"fmt"
	"net"
	"time"

	"github.com/Golds-l/goproxy/communication"
)

func MakeNewConn(communicationConn *net.Conn, listener net.Listener) net.Conn {
	_, err := (*communicationConn).Write([]byte("NEWXX")) // make new connection
	if err != nil {
		fmt.Printf("client connection error. %v\n", err)
	}
	newConn, newConnectionErr := listener.Accept()
	if newConnectionErr != nil {
		fmt.Printf("connection made error. %v\n", err)
	}
	return newConn
}

func KeepAliveS(conn *communication.Connection, listener net.Listener) {
	cache := make([]byte, 1024)
	for {
		_, writeErr := conn.Write("isAlive")
		if writeErr != nil {
			fmt.Printf("server communication connection write err %v\n", writeErr)
			fmt.Println("close and reconnect..")
			time.Sleep(1 * time.Second)
			_ = conn.Close()
			conn = communication.EstablishCommunicationConnS(listener)
			continue
		}
		n, readErr := conn.Read(cache)
		if readErr != nil {
			fmt.Printf("server communication connection read err %v\n", readErr)
			fmt.Println("close and reconnect in a second..")
			time.Sleep(1 * time.Second)
			_ = conn.Close()
			conn = communication.EstablishCommunicationConnS(listener)
			continue
		}
		if string(cache[:n]) == "alive" {
			fmt.Printf("communication %v alive.. %v\n", conn.Id, time.Now())
		}
		time.Sleep(3 * time.Second)
	}
}
