package communication

import (
	"fmt"
	"net"
)

func PrintSocket(conn net.Conn) {
	// for {
	cache := make([]byte, 1024)
	n, err := conn.Read(cache)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(cache[:n]), conn.LocalAddr(), conn.RemoteAddr())
	// }
}

func MakeNewConn(conn net.Conn, listener net.Listener) net.Conn {
	_, err := conn.Write([]byte("MKNC")) // make new connection
	if err != nil {
		fmt.Printf("client connection error. %v\n", err)
	}
	newConn, newConnectionErr := listener.Accept()
	if newConnectionErr != nil {
		fmt.Printf("connection made error. %v\n", err)
	}
	return newConn
}

func main() {
	listenerForLC, connectionForLCErr := net.Listen("tcp", "127.0.0.1:2000")
	if connectionForLCErr != nil {
		fmt.Printf("connectionForLC listen err. %v\n", connectionForLCErr)
	}
	communicationListener, communicationListenerErr := net.Listen("tcp", "127.0.0.1:1128")
	if communicationListenerErr != nil {
		fmt.Printf("communication listener err. %v\n", communicationListenerErr)
	}
	communicationConn, communicationConnErr := communicationListener.Accept()
	if communicationConnErr != nil {
		fmt.Printf("communication connection err. %v\n", communicationConnErr)
	}
	for {
		connForLC, err := listenerForLC.Accept()
		if err != nil {
			fmt.Printf("connForLC err, %v\n", connForLC)
		}
		newConn := MakeNewConn(communicationConn, communicationListener)
		_, _ = newConn.Write([]byte("connection made successfully"))
		fmt.Println(newConn)
	}
}
