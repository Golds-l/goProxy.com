package server

import (
	"fmt"
	"net"
	"time"
)

func MakeNewConn(communicationConn net.Conn, listener net.Listener) net.Conn {
	_, err := communicationConn.Write([]byte("NEWXX")) // make new connection
	if err != nil {
		fmt.Printf("client connection error. %v\n", err)
	}
	newConn, newConnectionErr := listener.Accept()
	if newConnectionErr != nil {
		fmt.Printf("connection made error. %v\n", err)
	}
	return newConn
}

func KeepAliveS(conn net.Conn) {
	cache := make([]byte, 1024)
	for {
		_, err := conn.Write([]byte("isAlive"))
		if err != nil {
			fmt.Println("communication connection err", err)
		}
		n, err := conn.Read(cache)
		if string(cache[:n]) == "XX" {
			time.Sleep(3 * time.Second)
			continue
		}
	}
}
