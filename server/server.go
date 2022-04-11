package server

import (
	"fmt"
	"net"
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
