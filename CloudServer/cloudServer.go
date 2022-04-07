package main

import (
	"fmt"
	"goProxy.com/args"
	"net"
	"os"
)

func CloseConn(conn1, conn2 net.Conn) {
	err1, err2 := conn1.Close(), conn2.Close()
	if err1 != nil || err2 != nil {
		fmt.Println("connect close err")
	}
}

func CloudServerToLocal(CloudServerR, CloudServerL net.Conn) {
	defer CloseConn(CloudServerR, CloudServerL)
	for {
		cache := make([]byte, 10240)
		readNum, _ := CloudServerR.Read(cache)
		_, _ = CloudServerL.Write(cache[:readNum])
	}
}

func LocalToCloudServer(CloudServerR, CloudServerL net.Conn) {
	defer CloseConn(CloudServerR, CloudServerL)
	for {
		cache := make([]byte, 10240)
		readNum, _ := CloudServerL.Read(cache)
		_, _ = CloudServerR.Write(cache[:readNum])
	}
}

func main() {
	argsMap, ok := args.GetArgsCloudServer()
	if ok == false {
		fmt.Println("args error")
		os.Exit(0)
	}
	fmt.Println(":"+argsMap["portLocal"], ":"+argsMap["portRemote"])
	os.Exit(0)
	listenLocal, err := net.Listen("tcp", ":"+argsMap["portLocal"])
	if err != nil {
		fmt.Println("err")
	}
	listenRemote, err := net.Listen("tcp", ":"+argsMap["portRemote"])
	if err != nil {
		fmt.Println("err")
	}
	fmt.Println("begin listen")
	for {
		cloudServerR, Err := listenRemote.Accept()
		if Err != nil {
			fmt.Println("err")
		}
		fmt.Println("remote client connect...")
		cloudServerL, _ := listenLocal.Accept()
		fmt.Println("local client connect...\nall connect")
		go CloudServerToLocal(cloudServerR, cloudServerL)
		go LocalToCloudServer(cloudServerR, cloudServerL)
	}
}
