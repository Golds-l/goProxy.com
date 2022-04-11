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
		case "-h":
			fmt.Println("this program run in remote machine\n-cS or -cloudServer\nserver address(ipv4 only)\n-cSP or -cloudServerPort\nserver port\n-rH or remoteHost\nLAN network address(ipv4 only)\n-rHP or -remoteHostPort\nLAN service port\nexp: ./remote -cS x.x.x.x -cSP 2000 -rH 127.0.0.1 -rHP 22 // ssh service")
			os.Exit(0)
		case "-help":
			fmt.Println("this program run in remote machine\n-cS or -cloudServer\nserver address(ipv4 only)\n-cSP or -cloudServerPort\nserver port\n-rH or remoteHost\nLAN network address(ipv4 only)\n-rHP or -remoteHostPort\nLAN service port\nexp: ./remote -cS x.x.x.x -cSP 2000 -rH 127.0.0.1 -rHP 22 // ssh service")
			os.Exit(0)
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
		case "-remoteHostPort":
			if isPort(os.Args[i+1]) {
				args["remoteHostPort"] = os.Args[i+1]
			} else {
				fmt.Println(os.Args[i+1], "port illegal")
			}
		case "-rHP":
			if isPort(os.Args[i+1]) {
				args["remoteHostPort"] = os.Args[i+1]
			} else {
				fmt.Println(os.Args[i+1], "port illegal")
			}
		case "-remoteHost":
			if isAddr(os.Args[i+1]) {
				args["remoteHost"] = os.Args[i+1]
			} else {
				fmt.Println(os.Args[i+1], "ip addr illegal")
			}
		case "-rH":
			if isAddr(os.Args[i+1]) {
				args["remoteHost"] = os.Args[i+1]
			} else {
				fmt.Println(os.Args[i+1], "ip addr illegal")
			}
		}
	}
	_, cloudServerAddrExist := args["CloudServer"]
	_, localHostAddrExist := args["remoteHost"]
	_, cloudServerPortExist := args["cloudServerPort"]
	_, localhostPortExist := args["remoteHostPort"]
	if cloudServerAddrExist && cloudServerPortExist && localhostPortExist && localHostAddrExist {
		return args, true
	} else {
		return args, false
	}
}
func GetArgsCloudServer() (map[string]string, bool) {
	args := make(map[string]string)
	for i := range os.Args {
		switch os.Args[i] {
		case "-h":
			fmt.Println("this program run in server\n-rp or -remotePort\nport for remote client\n-lP or -localPort\nport for connector")
			os.Exit(0)
		case "-help":
			fmt.Println("this program run in server\n-rp or -remotePort\nport for remote client\n-lP or -localPort\nport for connector")
			os.Exit(0)
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
