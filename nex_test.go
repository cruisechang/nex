package nex

import (
	"testing"
	"path/filepath"
	"os"
	"bytes"
	"log"
)

func TestNewGameServer(t *testing.T) {
	nx, err := NewNex("nexConfig.json")
	if err != nil {
		t.Errorf("TestNewGameServer error %s \n", err.Error())
		return
	}

	log.Printf("%v",nx.GetHallManager())
	log.Printf("%v",nx.GetRoomManager())
	log.Printf("get requester url :%s\n", nx.GetRequesterURL())


	c:=nx.GetConfig()
	wssAddr,wssPort:=c.WebSocketServerAddress()
	log.Printf("Address:%s, port=%s \n",wssAddr,wssPort)
	log.Printf("Socket connectionCapacity:%d\n",c.WebSocketServerConnectionCapacity())
	log.Printf("Socket ConnectionLostChanCapacity:%d\n",c.WebSocketServerConnectionLostChanCapacity())
	log.Printf("Socket PacketChanCapacity:%d\n",c.WebSocketServerPacketChanCapacity())
	log.Printf("Socket AcceptTimeoutSecond:%d\n",c.WebSocketServerAcceptTimeoutSecond())
	log.Printf("Socket AliveTimeoutSecond:%d\n",c.WebSocketServerAliveTimeoutSecond())
	log.Printf("Socket DefaultAESKeyString:%s\n",c.WebSocketServerDefaultAESKeyString())
	log.Printf("Socket DefaultAESKey:%v\n",c.WebSocketServerDefaultAESKey())

	addr,port:=c.HttpClientAddress()
	log.Printf("HttpClientAddress Address:%s, port=%s \n",addr,port)
	log.Printf("HttpClientHandshakeTimeoutSecond =%d \n",c.HttpClientHandshakeTimeoutSecond())
	log.Printf("HttpClientRequestTimeoutSecond =%d \n",c.HttpClientRequestTimeoutSecond())
	log.Printf("HttpClientTCPConnectTimeoutSecond =%d \n",c.HttpClientTCPConnectTimeoutSecond())


	addr,port=c.RPCServerAddress()
	log.Printf("RPCServerAddress Address:%s, port=%s \n",addr,port)


	log.Printf("RPCClient =%v,  \n",c.GetRPCClients())

}

func loadConfigFile(fileName string)string{
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
	}

	var buf bytes.Buffer
	buf.WriteString(dir)
	buf.WriteString("/")
	buf.WriteString(fileName)

	return buf.String()
}