package other

import (
	"fmt"

	"github.com/Golds-l/goproxy/server"
)

func CloseConn(conn1, conn2 *server.CloudConnection) {
	err1, err2 := conn1.Close(), conn2.Close()
	if err1 != nil || err2 != nil {
		fmt.Printf("connect close error.%v %v", err1, err2)
	}
}

func HandleErr(e error) {
	if e != nil {
		fmt.Println(e)
	}
}
