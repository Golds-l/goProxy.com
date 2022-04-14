package server

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/Golds-l/goproxy/communication"
	"github.com/Golds-l/goproxy/other"
)

func MakeNewConn(communicationConn *communication.Connection, listener net.Listener) (*communication.Connection, error) {
	readCache := make([]byte, 256)
	var conn communication.Connection
	conn.Id = other.GenerateConnId()
	_, _ = communicationConn.Write("NEWC:" + conn.Id) // make new connection
	n, _ := communicationConn.Read(readCache)
	mesgSlice := strings.Split(string(readCache[:n]), ":")
	if mesgSlice[0] == "NEW" && mesgSlice[1] == conn.Id {
		newConn, newConnectionErr := listener.Accept()
		if newConnectionErr != nil {
			fmt.Printf("connection etablished error. %v\n", newConnectionErr)
			return nil, newConnectionErr
		}
		conn.Conn = &newConn
		conn.Communication = false
		conn.StartTime = time.Now().Unix()
	}
	return &conn, nil
}

func KeepAliveS(conn *communication.Connection, listener net.Listener) {
	cache := make([]byte, 1024)
	for {
		_, writeErr := conn.Write("isAlive")
		if writeErr != nil {
			fmt.Printf("server communication connection write error. %v\n", writeErr)
			fmt.Println("close and reconnect..")
			time.Sleep(1 * time.Second)
			_ = conn.Close()
			conn = communication.EstablishCommunicationConnS(listener)
			continue
		}
		n, readErr := conn.Read(cache)
		if readErr != nil {
			fmt.Printf("server communication connection read error. %v\n", readErr)
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
