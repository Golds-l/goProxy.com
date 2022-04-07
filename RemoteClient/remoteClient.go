package main

import (
	"fmt"
	"goProxy.com/args"
	"net"
	"os"
	"time"
)

func CloseConn(conn1, conn2 net.Conn) {
	err1, err2 := conn1.Close(), conn2.Close()
	if err1 != nil || err2 != nil {
		fmt.Println("connect close err")
	}
}

func RemoteClientToCloudServer(RemoteClient, SSHRemoteClient net.Conn) {
	defer CloseConn(RemoteClient, SSHRemoteClient)
	for {
		cache := make([]byte, 10240)
		readNum, _ := SSHRemoteClient.Read(cache)
		_, _ = RemoteClient.Write(cache[:readNum])
	}
}

func CloudServerToRemoteClient(RemoteClient, SSHRemoteClient net.Conn) {
	defer CloseConn(RemoteClient, SSHRemoteClient)
	for {
		cache := make([]byte, 10240)
		readNum, _ := RemoteClient.Read(cache)
		_, _ = SSHRemoteClient.Write(cache[:readNum])
	}
}

func main() {
	argsMap, ok := args.GetArgsRemoteClient()
	fmt.Printf("server:%v ", argsMap["CloudServer"]+":"+argsMap["cloudServerPort"])
	fmt.Printf("localhostPort:%v\n", argsMap["localhostPort"])
	if ok == false {
		fmt.Println("args illegal")
		os.Exit(0)
	}
	address := argsMap["CloudServer"] + ":" + argsMap["cloudServerPort"]
	ssh, _ := net.Dial("tcp", "127.0.0.1:"+argsMap["localhostPort"])
	remote, _ := net.Dial("tcp", address)
	go RemoteClientToCloudServer(remote, ssh)
	go CloudServerToRemoteClient(remote, ssh)
	for {
		time.Sleep(1 * time.Second)
	}
}
