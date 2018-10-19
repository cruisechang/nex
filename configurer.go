package nex

import (
	"io/ioutil"
	"encoding/json"
	"errors"
	"time"
)

//Config config main struct
type Configurer struct {
	data confData
}

//NewConfig make a new config struct
func NewConfigurer(filePath string) (*Configurer, error) {

	cf := &Configurer{
		data: confData{},
	}
	if err := cf.loadConfig(filePath); err != nil {
		return nil, err
	}
	return cf, nil
}

//LoadConfig loads config file.
func (c *Configurer) loadConfig(filePath string) error {

	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	//unmarshal to struct
	if err := json.Unmarshal(b, &c.data); err != nil {
		return err
	}

	//aeskey string to []byte
	c.SetWebSocketServerDefaultAESKey([]byte(c.WebSocketServerDefaultAESKeyString()))

	return nil
}

type confData struct {
	WebSocketServer *webSocketServerConf
	HttpServer      *httpServerConf
	HttpClient      *httpClientConf
	RPCServer       *rpcServerConf
	RPCClient       []rpcClientConf
}
type webSocketServerConf struct {
	Address string
	Port    string
	//ConnectionCapacity defines max of connection
	ConnectionCapacity int
	//ConnectionLostChanCapacity defines lost connection chan capacity
	//chan exceed its capacity wiil blocking until others read item from chan.
	ConnectionLostChanCapacity int
	//PacketChanCapacity defines max of packet chan
	PacketChanCapacity int

	DefaultAESKeyString string

	// for default aes setted in gameServer LoadConfig
	DefaultAESKey []byte

	AcceptTimeoutSecond int
	AliveTimeoutSecond  int
}

type httpServerConf struct {
	Address            string
	Port               string
	ReadTimeoutSecond  int
	WriteTimeoutSecond int
	IdleTimeoutSecond  int
	MaxHeaderBytes     int
}

//Requester is request handerl
//used to request backstage http api
//no web service for outter
type httpClientConf struct {
	Address                 string
	Port                    string
	TCPConnectTimeoutSecond int //connect timeout
	HandshakeTimeoutSecond  int //handshake timeout
	RequestTimeoutSecond    int //request timeout
}

type rpcServerConf struct {
	Address string
	Port    string
}

type rpcClientConf struct {
	Address string
	Port    string
	Name    string
}

//websocket
func (c *Configurer) WebSocketServerAddress() (addr, port string) {
	return c.data.WebSocketServer.Address, c.data.WebSocketServer.Port
}
func (c *Configurer) WebSocketServerConnectionCapacity() int {
	return c.data.WebSocketServer.ConnectionCapacity
}
func (c *Configurer) WebSocketServerConnectionLostChanCapacity() int {
	return c.data.WebSocketServer.ConnectionLostChanCapacity
}
func (c *Configurer) WebSocketServerPacketChanCapacity() int {
	return c.data.WebSocketServer.PacketChanCapacity
}
func (c *Configurer) WebSocketServerDefaultAESKeyString() string {
	return c.data.WebSocketServer.DefaultAESKeyString
}
func (c *Configurer) WebSocketServerDefaultAESKey() []byte {
	return c.data.WebSocketServer.DefaultAESKey
}
func (c *Configurer) SetWebSocketServerDefaultAESKey(key []byte) {
	c.data.WebSocketServer.DefaultAESKey = key
}
func (c *Configurer) WebSocketServerAcceptTimeoutSecond() int {
	return c.data.WebSocketServer.AcceptTimeoutSecond
}
func (c *Configurer) WebSocketServerAliveTimeoutSecond() int {
	return c.data.WebSocketServer.AliveTimeoutSecond
}

//http
func (c *Configurer) HttpServerAddress() (addr, port string) {
	return c.data.HttpServer.Address, c.data.HttpServer.Port
}
func (c *Configurer) HttpServerReadTimeout() time.Duration {
	return time.Second*time.Duration(c.data.HttpServer.ReadTimeoutSecond)
}
func (c *Configurer) HttpServerWriteTimeout() time.Duration {
	return time.Second*time.Duration(c.data.HttpServer.WriteTimeoutSecond)
}
func (c *Configurer) HttpServerIdleTimeout() time.Duration {
	return time.Second*time.Duration(c.data.HttpServer.IdleTimeoutSecond)
}
func (c *Configurer) HttpServerMaxHeaderBytes() int {
	return c.data.HttpServer.MaxHeaderBytes
}
//http client

func (c *Configurer) HttpClientAddress() (addr, port string) {
	return c.data.HttpClient.Address, c.data.HttpClient.Port
}
func (c *Configurer) HttpClientTCPConnectTimeoutSecond() int {
	return c.data.HttpClient.TCPConnectTimeoutSecond
}
func (c *Configurer) HttpClientHandshakeTimeoutSecond() int {
	return c.data.HttpClient.HandshakeTimeoutSecond
}

func (c *Configurer) HttpClientRequestTimeoutSecond() int {
	return c.data.HttpClient.RequestTimeoutSecond
}

//rpc server
func (c *Configurer) RPCServerAddress() (addr, port string) {
	return c.data.RPCServer.Address, c.data.RPCServer.Port
}

//rpc client
func (c *Configurer) GetRPCClients() []rpcClientConf {
	return c.data.RPCClient
}

func (c *Configurer) GetRPCClientByIndex(index int) (rpcClientConf, error) {
	if index >= len(c.data.RPCClient) {
		return rpcClientConf{}, errors.New("index out of range")
	}
	return c.data.RPCClient[index], nil
}

func (c *Configurer) GetRPCClientByName(name string) (rpcClientConf, error) {
	for _, v := range c.data.RPCClient {
		if v.Name == name {
			return v, nil
		}
	}
	return rpcClientConf{}, errors.New("rpcClient not found")
}
