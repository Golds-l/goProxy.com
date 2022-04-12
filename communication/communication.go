package communication

import (
	"fmt"
	"net"
	"time"

	"github.com/Golds-l/goproxy/other"
)

func CloudServerToLocal(CloudServerR, CloudServerL net.Conn) {
	defer other.CloseConn(CloudServerR, CloudServerL)
	for {
		cache := make([]byte, 10240)
		readNum, _ := CloudServerR.Read(cache)
		_, _ = CloudServerL.Write(cache[:readNum])
	}
}

func LocalToCloudServer(CloudServerR, CloudServerL net.Conn) {
	defer other.CloseConn(CloudServerR, CloudServerL)
	for {
		cache := make([]byte, 10240)
		readNum, _ := CloudServerL.Read(cache)
		_, _ = CloudServerR.Write(cache[:readNum])
	}
}

func RemoteClientToCloudServer(RemoteClient, SSHRemoteClient net.Conn) {
	defer other.CloseConn(RemoteClient, SSHRemoteClient)
	for {
		cache := make([]byte, 10240)
		readNum, _ := SSHRemoteClient.Read(cache)
		_, _ = RemoteClient.Write(cache[:readNum])
	}
}

func CloudServerToRemoteClient(RemoteClient, SSHRemoteClient net.Conn) {
	defer other.CloseConn(RemoteClient, SSHRemoteClient)
	for {
		cache := make([]byte, 10240)
		readNum, _ := RemoteClient.Read(cache)
		_, _ = SSHRemoteClient.Write(cache[:readNum])
	}
}

func EstablishCommunicationConnC(addr string) net.Conn {
	var communicationConn net.Conn
	communicationConnACK := make([]byte, 512)
	for {
		conn, connErr := net.Dial("tcp", addr)
		if connErr != nil {
			fmt.Println("Can not connect cloud server... Retry in a second")
			time.Sleep(1 * time.Second)
			continue
		}
		n, readErr := conn.Read(communicationConnACK)
		if readErr != nil {
			fmt.Printf("coon read err!%v\n", readErr)
			_ = conn.Close()
			continue
		}
		if string(communicationConnACK[:n]) == "cloudXX" {
			fmt.Println("established")
			communicationConn = conn
			_, _ = communicationConn.Write([]byte("RCReady"))
			break
		}
		_ = conn.Close()
	}
	return communicationConn
}

func EstablishCommunicationConnS(serverListener net.Listener) net.Conn {
	var communicationConn net.Conn
	connACK := make([]byte, 512)
	for {
		conn, e := serverListener.Accept()
		if e != nil {
			fmt.Println("Can not connect cloud server... Retry in a second")
			time.Sleep(1 * time.Second)
			continue
		}
		_, writeErr := conn.Write([]byte("cloudXX"))
		if writeErr != nil {
			fmt.Printf("connection write err! %v\n", writeErr)
			fmt.Printf("connection:%v will be closed\n", conn)
			_ = conn.Close()
		}
		n, readErr := conn.Read(connACK)
		if readErr != nil {
			fmt.Printf("connection read err! %v\n", writeErr)
			fmt.Printf("connection:%v will be closed\n", conn)
			_ = conn.Close()
		}
		if string(connACK[:n]) == "RCReady" {
			fmt.Println("cloud server<--->remote client is connected!")
			break
		}
		_ = conn.Close()
	}
	return communicationConn
}
