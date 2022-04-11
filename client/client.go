package client

import (
	"fmt"
	"net"

	"github.com/Golds-l/goproxy/other"
)

func MakeNewClient(serverAddr, localAddr string) (net.Conn, net.Conn) {
	connServer, err := net.Dial("tcp", serverAddr)
	other.HandleErr(err)
	fmt.Println("make new client")
	other.HandleErr(err)
	connLocal, connLOcalErr := net.Dial("tcp", localAddr)
	other.HandleErr(connLOcalErr)
	return connServer, connLocal
}
