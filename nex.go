package nex

import (
	"os"
	"os/signal"
	"syscall"
	"github.com/cruisechang/nex/entity"
	"github.com/cruisechang/nex/websocket"
	"time"
	goLog "log"
	"fmt"
	"google.golang.org/grpc"
	"reflect"
	"net"
	"strings"
	"errors"
)

const (
	version = "2.0"
)

//Nex is used for main process
type Nex interface {
	Start()
	Stop()
	Version() string
	GetConfig() *Configurer

	//RegisterCommandProcessor registers a process to handle command
	RegisterCommandProcessor(key string, p CommandProcessor) error
	CreateCommand(code int, step int, cmdName string, data string) (*Command, error)
	SendCommand(cmd *Command, sender entity.User, targetConnIDs []string, lowercase bool)

	//RegisterEventProcessor registers a procese to handle event
	RegisterEventProcessor(code uint16, p EventProcessor) error
	UnRegisterProcessor(code uint16)
	DispatchEvent(code uint16) error

	//RegisterUpdateProcessor
	RegisterUpdateProcessor(name string, p UpdateProcessor) error
	UnRegisterUpdateProcessor(name string) error
	StartUpdateProcessor(name string) (resErr error)
	StopUpdateProcessor(name string) (resErr error)
	IsUpdateProcessorRunning(name string) (bool, error)

	//logger
	GetLogger() Logger

	//requester
	GetRequesterURL() string
	SetRequesterPostURI(path string, queryPair map[string]string) error
	GetRequesterPostURI() string
	GetRequesterPostQuery() string
	RequesterPost() (string, error)

	//user management
	GetUser(userID int) (entity.User, bool)
	DisconnectUser(userID int)
	RemoveUser(userID int)
	GetUsers() []entity.User
	GetUserConnIDs() []string

	//manager
	GetHallManager() HallManager
	GetRoomManager() RoomManager

	//for grpc server
	//fn is grpc protobuf serivce register function.
	//example pb.RegisterGreeterServer(s, &server{}),  s==*grpc.Server, &server{}=serverStruct{} contain rpc methods.
	StartGRPCServer(addr, port string, registerServerFn interface{}, serverStruct interface{})

	GetGRPCClient(addr, port string, newClientFunc interface{}) (interface{}, error)
	//for grpc client
	GetGRPClientConn(address, port string) (*grpc.ClientConn, error)
}

//Nex is server instance used for game zone.
type nex struct {
	active         bool
	version        string
	config         *Configurer //config from config package
	socketService  websocket.WebsocketService
	commandManager CommandManager
	eventManager   EventManager
	grpcServer     *grpc.Server

	hallManager       HallManager
	roomManager       RoomManager
	userManager       UserManager
	requester         Requester
	packetHandler     PacketHandler
	middleHandler     MiddleHandler
	updateManager     UpdateManager
	logger            Logger
	receivePacketChan <-chan *websocket.ReceivePacketData //Get from socketService read only
	lostConnIDChan    <-chan string                       //websocket connection lost chan contains connID
	shutdown          chan int
}

//NewNex used to new Nex instance.
//Call this to new not Nex{}.
func NewNex(configFilePath string) (Nex, error) {
	gs := &nex{
		version:  version,
		shutdown: make(chan int),
	}

	con, err := NewConfigurer(configFilePath)

	if err != nil {
		return nil, err
	}
	gs.config = con
	err = gs.initial()
	return gs, err
}

//Initial to init Nex
//call this after new Nex
func (nx *nex) initial() error {
	nl, err := NewLogger()
	if err != nil {

		goLog.Fatalf("NewLogger error:%s\n", err.Error())
		return err
	}

	nx.logger = nl
	nx.logger.SetLevel(LogLevelInfo)

	goLog.Printf("NewLogger success")

	if nx.socketService == nil {
		addr, port := nx.config.WebSocketServerAddress()

		so, err := websocket.NewWebsocketService(addr,
			port,
			nx.config.WebSocketServerDefaultAESKey(),
			nx.config.WebSocketServerConnectionCapacity(),
			nx.config.WebSocketServerPacketChanCapacity(),
			nx.config.WebSocketServerConnectionLostChanCapacity(),
			time.Second*time.Duration(nx.config.WebSocketServerAcceptTimeoutSecond()),
			time.Second*time.Duration(nx.config.WebSocketServerAliveTimeoutSecond()))
		if err != nil {
			nx.logger.Log(LogLevelError, fmt.Sprintf("NewWebsocketService error:%s\n", err.Error()))
			return err
		}
		nx.socketService = so
		nx.logger.Log(LogLevelInfo, "NewWebsocketService success\n")
	}

	ph, err := NewPacketHandler()
	if err != nil {
		nx.logger.Log(LogLevelError, fmt.Sprintf("NewPacketHandler error:%s\n", err.Error()))
		return err
	}
	nx.packetHandler = ph
	nx.logger.Log(LogLevelInfo, "NewPacketHandler success\n")

	mh, err := NewMiddleHandler()
	if err != nil {
		nx.logger.Log(LogLevelError, fmt.Sprintf("NewMiddleHandler error:%s\n", err.Error()))
		return err
	}
	nx.middleHandler = mh
	nx.logger.Log(LogLevelInfo, "NewMiddleHandler success\n")

	ch, err := NewCommandManager()
	if err != nil {
		nx.logger.Log(LogLevelError, fmt.Sprintf("NewCommandManager error:%s\n", err.Error()))
		return err
	}
	nx.commandManager = ch
	nx.logger.Log(LogLevelInfo, "NewCommandManager success\n")

	em, err := NewEventManager()
	if err != nil {
		nx.logger.Log(LogLevelError, fmt.Sprintf("NewEventManager error:%s\n", err.Error()))
		return err
	}
	nx.eventManager = em
	nx.logger.Log(LogLevelInfo, "NewEventManager success\n")

	httpCAddr, httpCPort := nx.config.HttpClientAddress()
	re, err := NewRequester(httpCAddr,
		httpCPort,
		nx.config.HttpClientTCPConnectTimeoutSecond(),
		nx.config.HttpClientHandshakeTimeoutSecond(),
		nx.config.HttpClientRequestTimeoutSecond())
	if err != nil {
		nx.logger.Log(LogLevelError, fmt.Sprintf("NewRequester error:%s\n", err.Error()))
		return err
	}
	nx.requester = re
	nx.logger.Log(LogLevelInfo, "NewRequester success\n")

	hm, err := NewHallManager()
	if err != nil {
		nx.logger.Log(LogLevelError, fmt.Sprintf("NewHallManager error:%s\n", err.Error()))
		return err
	}
	nx.hallManager = hm
	nx.logger.Log(LogLevelInfo, "NewHallManager success\n")

	rm, err := NewRoomManager()
	if err != nil {
		nx.logger.Log(LogLevelError, fmt.Sprintf("NewRoomManager error:%s\n", err.Error()))
		return err
	}
	nx.roomManager = rm
	nx.logger.Log(LogLevelInfo, "NewRoomManager success\n")

	um, err := NewUserManager()
	if err != nil {
		nx.logger.Log(LogLevelError, fmt.Sprintf("NewUserManager error:%s\n", err.Error()))
		return err
	}

	nx.userManager = um
	nx.logger.Log(LogLevelInfo, "NewUserManager success\n")

	upm, err := NewUpdateManager()
	if err != nil {
		nx.logger.Log(LogLevelError, fmt.Sprintf("NewUpdateManager error:%s\n", err.Error()))
		return err
	}
	nx.updateManager = upm
	nx.logger.Log(LogLevelInfo, "NewUpdateManager success\n")
	return nil
}
func (nx *nex) GetConfig() *Configurer {
	return nx.config
}

//Start starts listening
func (nx *nex) Start() {
	nx.receivePacketChan = nx.socketService.GetReceivePacketChan()
	nx.lostConnIDChan = nx.socketService.GetLostConnIDChan()
	nx.socketService.Start() //go accept
	nx.logger.Log(LogLevelInfo, "Nex started...\n")
	nx.mainFlow()
	nx.handleCtrlC()
}

//Stop user manager process.
func (nx *nex) Stop() {
	nx.active = false
	nx.socketService.Stop()
	nx.logger.Log(LogLevelInfo, "Nex stopped")
}

func (nx *nex) Version() string {
	return nx.version
}

func (nx *nex) handleCtrlC() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		nx.shutdown <- 0
	}()
}

/*
  Following are for command
 */
//RegisterCommandProcessor is a wrapper for commandHandler.RegisterProcessor
func (nx *nex) RegisterCommandProcessor(key string, p CommandProcessor) error {
	return nx.commandManager.RegisterProcessor(key, p)
}
func (nx *nex) UnRegisterCommandProcessor(key string) {
	nx.commandManager.UnRegisterProcessor(key)
}

//CreateCommand return
func (nx *nex) CreateCommand(code int, step int, cmdName string, data string) (*Command, error) {
	return nx.commandManager.CreateCommand(code, step, cmdName, data)
}

//SendCommand sends command to target connID
//sender is user who send this command
func (nx *nex) SendCommand(cmd *Command, sender entity.User, targetConnIDs []string, lowercase bool) {

	var p []byte
	var err error

	if lowercase {

		nc := CommandLowerCase{
			Code:    cmd.Code,
			Command: cmd.Command,
			Step:    cmd.Step,
			Data:    cmd.Data,
		}

		p, err = nx.middleHandler.CommandToPacketLowercase(&nc)

	} else {
		p, err = nx.middleHandler.CommandToPacket(cmd)
	}

	if err != nil {
		return
	}
	nx.socketService.Send(targetConnIDs, p)
}

/*

 Following are for event
 */
func (nx *nex) RegisterEventProcessor(code uint16, p EventProcessor) error {
	return nx.eventManager.RegisterProcessor(code, p)
}
func (nx *nex) UnRegisterProcessor(code uint16) {
	nx.eventManager.UnRegisterProcessor(code)
}
func (nx *nex) DispatchEvent(code uint16) error {
	return nx.eventManager.DispatchEvent(code)
}

/*
followings are for updater
 */
func (nx *nex) RegisterUpdateProcessor(name string, p UpdateProcessor) error {
	return nx.updateManager.RegisterProcessor(name, p)
}

/*
 unregister 一定要檢查是否stop 否則可能run 不停
 */
func (nx *nex) UnRegisterUpdateProcessor(name string) error {
	return nx.updateManager.UnRegisterProcessor(name)

}

func (nx *nex) StartUpdateProcessor(name string) (resErr error) {
	return nx.updateManager.StartProcessor(name)
}

func (nx *nex) StopUpdateProcessor(name string) (resErr error) {
	return nx.updateManager.StopProcessor(name)

}

func (nx *nex) IsUpdateProcessorRunning(name string) (bool, error) {
	return nx.updateManager.IsRunning(name)
}

//
func (nx *nex) GetHallManager() HallManager {
	return nx.hallManager
}

func (nx *nex) GetRoomManager() RoomManager {
	return nx.roomManager
}

func (nx *nex) GetLogger() Logger {
	return nx.logger
}

func (nx *nex) GetRequesterURL() string {
	return nx.requester.URL()
}
func (nx *nex) SetRequesterPostURI(path string, queryPair map[string]string) error {
	return nx.requester.SetPostURI(path, queryPair)
}
func (nx *nex) GetRequesterPostURI() string {
	return nx.requester.PostURI()
}
func (nx *nex) GetRequesterPostQuery() string {
	return nx.requester.PostQuery()
}
func (nx *nex) RequesterPost() (string, error) {
	return nx.requester.Post()
}
func (nx *nex) RemoveUser(userID int) {
	nx.userManager.RemoveUser(userID)
}

func (nx *nex) DisconnectUser(userID int) {
	if u, ok := nx.userManager.GetUser(userID); ok {
		nx.socketService.Disconnect(u.ConnID())

	}
}

func (nx *nex) GetUser(userID int) (entity.User, bool) {
	return nx.userManager.GetUser(userID)
}

func (nx *nex) GetUsers() []entity.User {
	return nx.userManager.GetUsers()
}

func (nx *nex) GetUserConnIDs() []string {
	users := nx.userManager.GetUsers()

	ids := []string{}
	for i, v := range users {
		ids[i] = v.ConnID()

	}
	return ids
}

//grpc
//registerServerFunc is grpc protobuf serivce register function.
//example pb.RegisterGreeterServer(s, &server{}),  s==*grpc.Server, &server{}=serverStruct{} contain rpc methods.
func (nx *nex) StartGRPCServer(addr, port string, registerServerFunc interface{}, serverStruct interface{}) {

	go func() {
		defer func() {
			if r := recover(); r != nil {
				goLog.Fatalf("StartGRPCServer panic %v", r)
				nx.logger.Log(LogLevelPanic, fmt.Sprintf("StartGRPCServer panic=%v\n", r))
				return
			}

		}()

		if strings.Index(reflect.TypeOf(registerServerFunc).String(), "func") != 0 {
			goLog.Fatalf("fn is not function ")
			nx.logger.Log(LogLevelError, fmt.Sprintf("StartGRPCServer fn is not function\n"))
			return
		}

		//check
		if nx.grpcServer == nil {
			nx.grpcServer = grpc.NewServer()
		}

		//server listen
		addr, _ := net.ResolveTCPAddr("tcp", addr+":"+port)
		lis, err := net.ListenTCP("tcp", addr)
		if err != nil {
			goLog.Fatalf("StartGRPCServer listen error=%s", err.Error())
			nx.logger.Log(LogLevelError, fmt.Sprintf("StartGRPCServer listen error=%s\n", err.Error()))
			return
		}

		value := reflect.ValueOf(registerServerFunc)
		value.Call([]reflect.Value{reflect.ValueOf(nx.grpcServer), reflect.ValueOf(serverStruct)})

		if err := nx.grpcServer.Serve(lis); err != nil {
			goLog.Fatalf("StartGRPCServer serve error=%s", err.Error())
			nx.logger.Log(LogLevelError, fmt.Sprintf("StartGRPCServer serve error=%s\n", err.Error()))
			return
		}

	}()
}

func (nx *nex) GetGRPCClient(addr, port string, newClientFunc interface{}) (interface{}, error) {

	conn, err := grpc.Dial(addr+":"+port, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	value := reflect.ValueOf(newClientFunc)
	if value.Kind().String() != "func" {
		return nil, errors.New(" passed newClientFunc is no  a function")
	}
	cli := value.Call([]reflect.Value{reflect.ValueOf(conn)}) //retrun []Value

	//value to original type
	//value.Interface().(oriType)
	return cli[0].Interface(), nil //value to interface
}

//pass port slice
func (nx *nex) GetGRPClientConn(address, port string) (*grpc.ClientConn, error) {

	// Set up a connection to the server.
	return grpc.Dial(address+port, grpc.WithInsecure())
}

//mainFlow
//blocking
//receive data and process it
func (nx *nex) mainFlow() {
	defer func() {
		if r := recover(); r != nil {
			nx.logger.Log(LogLevelPanic, fmt.Sprintf("mainFlow panic:%v", r))
		}
	}()

	for {
		select {

		//ctrl-c to shtdown
		case <-nx.shutdown:
			nx.Stop()
			os.Exit(0)
		case pd, opend := <-nx.receivePacketChan:

			//check if chan closed
			if opend {
				dataBody, connUUID, err := nx.packetHandler.Handle(pd)
				if err != nil {
					//要斷掉連線
					nx.logger.LogFile(LogLevelError, fmt.Sprintf("nex acketHandler.Handle err=%s\n", err.Error()))
					nx.socketService.Disconnect(connUUID)
					continue
				}

				cmd, err := nx.middleHandler.PacketToCommand(dataBody)

				if err != nil {
					nx.logger.LogFile(LogLevelError, fmt.Sprintf("nex middleHandler.PacketToCommand err=%s\n", err.Error()))
					nx.socketService.Disconnect(connUUID)
					continue
				}

				//fmt.Printf("%#v\n", cmd)
				user, ok := nx.userManager.GetUserByConnID(connUUID)
				if !ok {
					user, err = nx.userManager.CreateUser(connUUID)
				}
				nx.commandManager.RunProcessor(&CommandObject{Cmd: cmd, User: user})
			}

		case connID, opend := <-nx.lostConnIDChan:
			if opend {
				if user, ok := nx.userManager.GetUserByConnID(connID); ok {

					nx.eventManager.RunProcessor(&EventObject{
						Code: EventUserLost,
						User: user,
					})
				}

			}

		}
	}
}
