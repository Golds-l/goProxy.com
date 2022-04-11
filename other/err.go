package other

import (
	"fmt"
	"net"
)

func CloseConn(conn1, conn2 net.Conn) {
	err1, err2 := conn1.Close(), conn2.Close()
	if err1 != nil || err2 != nil {
		fmt.Println("connect close err")
	}
}

func HandleErr(e error) {
	if e != nil {
		fmt.Println(e)
	}
}
