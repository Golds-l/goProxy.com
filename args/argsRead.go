package args

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func isAddr(ipv4Addr string) bool {
	var addr = strings.Split(ipv4Addr, ".")
	if len(addr) < 4 {
		return false
	}
	for i := range addr {
		intI, err := strconv.Atoi(addr[i])
		if err != nil {
			return false
		}
		if intI < 0 || intI > 255 {
			return false
		}
	}
	return true
}

func isPort(port string) bool {
	portInt, err := strconv.Atoi(port)
	if err != nil {
		fmt.Println(err)
		return false
	}
	if portInt < 1023 || portInt > 65535 {
		return false
	}
	return true
}

func GetArgsRemoteClient() (map[string]string, bool) {
	args := make(map[string]string)
	for i := range os.Args {
		switch os.Args[i] {
		case "-CloudServer":
			if isAddr(os.Args[i+1]) {
				args["CloudServer"] = os.Args[i+1]
			} else {
				fmt.Println(os.Args[i+1], "ip addr illegal")
			}
		case "-cS":
			if isAddr(os.Args[i+1]) {
				args["CloudServer"] = os.Args[i+1]
			} else {
				fmt.Println(os.Args[i+1], "ip addr illegal")
			}
		case "-port":
			if isPort(os.Args[i+1]) {
				args["port"] = os.Args[i+1]
			} else {
				fmt.Println(os.Args[i+1], "port illegal")
			}
		case "-p":
			if isPort(os.Args[i+1]) {
				args["port"] = os.Args[i+1]
			} else {
				fmt.Println(os.Args[i+1], "port illegal")
			}
		}
	}
	_, addrExist := args["CloudServer"]
	_, portExist := args["port"]
	if addrExist && portExist {
		return args, true
	} else {
		return args, false
	}
}
func GetArgsCloudServer() (map[string]string, bool) {
	args := make(map[string]string)
	for i := range os.Args {
		switch os.Args[i] {
		case "-portLocal":
			if isPort(os.Args[i+1]) {
				args["portLocal"] = os.Args[i+1]
			} else {
				fmt.Println(os.Args[i+1], "local port illegal")
			}
		case "-portRemote":
			if isPort(os.Args[i+1]) {
				args["portRemote"] = os.Args[i+1]
			} else {
				fmt.Println(os.Args[i+1], "remote port illegal")
			}
		case "-pL":
			if isPort(os.Args[i+1]) {
				args["portLocal"] = os.Args[i+1]
			} else {
				fmt.Println(os.Args[i+1], "local port illegal")
			}
		case "-pR":
			if isPort(os.Args[i+1]) {
				args["portRemote"] = os.Args[i+1]
			} else {
				fmt.Println(os.Args[i+1], "remote port illegal")
			}
		}
	}
	_, portLocalExist := args["portLocal"]
	_, portRemoteExist := args["portRemote"]
	if portLocalExist && portRemoteExist && args["portLocal"] != args["portRemote"] {
		return args, true
	} else {
		return args, false
	}
}
