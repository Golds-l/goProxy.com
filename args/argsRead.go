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
	if portInt < 0 || portInt > 65535 {
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
		case "-cloudServerPort":
			if isPort(os.Args[i+1]) {
				args["cloudServerPort"] = os.Args[i+1]
			} else {
				fmt.Println(os.Args[i+1], "port illegal")
			}
		case "-cSP":
			if isPort(os.Args[i+1]) {
				args["cloudServerPort"] = os.Args[i+1]
			} else {
				fmt.Println(os.Args[i+1], "port illegal")
			}
		case "-localhostPort":
			if isPort(os.Args[i+1]) {
				args["localhostPort"] = os.Args[i+1]
			} else {
				fmt.Println(os.Args[i+1], "port illegal")
			}
		case "lP":
			if isPort(os.Args[i+1]) {
				args["localhostPort"] = os.Args[i+1]
			} else {
				fmt.Println(os.Args[i+1], "port illegal")
			}
		}
	}
	_, addrExist := args["CloudServer"]
	_, cloudServerPortExist := args["cloudServerPort"]
	_, localhostPortExist := args["localhostPort"]
	if addrExist && cloudServerPortExist && localhostPortExist {
		return args, true
	} else {
		return args, false
	}
}
func GetArgsCloudServer() (map[string]string, bool) {
	args := make(map[string]string)
	for i := range os.Args {
		switch os.Args[i] {
		case "-localPort":
			if isPort(os.Args[i+1]) {
				args["localPort"] = os.Args[i+1]
			} else {
				fmt.Println(os.Args[i+1], "local port illegal")
			}
		case "-remotePort":
			if isPort(os.Args[i+1]) {
				args["remotePort"] = os.Args[i+1]
			} else {
				fmt.Println(os.Args[i+1], "remote port illegal")
			}
		case "-lP":
			if isPort(os.Args[i+1]) {
				args["localPort"] = os.Args[i+1]
			} else {
				fmt.Println(os.Args[i+1], "local port illegal")
			}
		case "-rP":
			if isPort(os.Args[i+1]) {
				args["remotePort"] = os.Args[i+1]
			} else {
				fmt.Println(os.Args[i+1], "remote port illegal")
			}
		}
	}
	_, portLocalExist := args["localPort"]
	_, portRemoteExist := args["remotePort"]
	if portLocalExist && portRemoteExist && args["localPort"] != args["remotePort"] {
		return args, true
	} else {
		return args, false
	}
}
