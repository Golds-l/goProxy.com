package server

import (
	"fmt"
	"github.com/Golds-l/goproxy/communication"
	"net"
	"time"
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

func KeepAliveS(conn *net.Conn, listener net.Listener) {
	cache := make([]byte, 1024)
	//go communication.WriteAlive(conn, "isAlive")
	for {
		_, writeErr := (*conn).Write([]byte("isAlive"))
		if writeErr != nil {
			fmt.Printf("server communication connection write err %v\n", writeErr)
			fmt.Println("close and reconnect..")
			time.Sleep(1 * time.Second)
			_ = (*conn).Close()
			conn = communication.EstablishCommunicationConnS(listener)
			continue
		}
		n, readErr := (*conn).Read(cache)
		if readErr != nil {
			fmt.Printf("server communication connection read err %v\n", readErr)
			fmt.Println("close and reconnect in a second..")
			time.Sleep(1 * time.Second)
			_ = (*conn).Close()
			conn = communication.EstablishCommunicationConnS(listener)
			continue
		}
		if string(cache[:n]) == "alive" {
			fmt.Printf("server alive.. %v\n", time.Now())
		}
		time.Sleep(3 * time.Second)
	}
}
