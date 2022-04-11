package communication

import (
	"net"

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
