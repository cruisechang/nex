package nex

import (
	"os"
	"os/signal"
	"syscall"
	"github.com/cruisechang/nex/entity"
	"github.com/cruisechang/nex/builtinEvent"
	"github.com/cruisechang/nex/event"
	"github.com/cruisechang/nex/websocket"
	nxhttp "github.com/cruisechang/nex/http"
	nxlog "github.com/cruisechang/nex/log"
	nxRPC "github.com/cruisechang/nex/rpc"
	"time"
	goLog "log"
	"fmt"
	"google.golang.org/grpc"
	"net/http"
	"github.com/juju/errors"
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
	RegisterCommandProcessor(name string, p CommandProcessor) error
	CreateCommand(code int, step int, cmdName string, data string) (*Command, error)
	SendCommand(cmd *Command, sender entity.User, targetConnIDs []string, lowercase bool)

	//RegisterBuiltinEventProcessor registers a procese to handle event
	RegisterBuiltinEventProcessor(code builtinEvent.EventNO, p builtinEvent.EventProcessor) error
	UnRegisterBuiltinEventProcessor(code builtinEvent.EventNO)

	//RegisterBuiltinEventProcessor registers a procese to handle event
	RegisterEventProcessor(event string, p event.EventProcessor) error
	UnRegisterEventProcessor(event string)
	DispatchEvent(event string)

	//RegisterUpdateProcessor
	RegisterUpdateProcessor(name string, p UpdateProcessor) error
	UnRegisterUpdateProcessor(name string) error
	StartUpdateProcessor(name string) (resErr error)
	StopUpdateProcessor(name string) (resErr error)
	IsUpdateProcessorRunning(name string) (bool, error)

	//logger
	GetLogger() nxlog.Logger

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
	//StartGRPCServer(addr, port string, registerServerFn interface{}, serverStruct interface{})
	StartGRPCServer(registerServerFn interface{}, serverStruct interface{}) error
	GetGRPCClient(addr, port string, newClientFunc interface{}) (interface{}, error)
	GetGRPClientConn(address, port string) (*grpc.ClientConn, error)

	//http server
	StartHTTPServer(parameters *nxhttp.ServerParameters) error
	StopHTTPServer()
	GetHTTPServerRequestRealAddr(r *http.Request) (string, error)
	RegisterHTTPServerHandler(pattern string, handler http.Handler)
}

//Nex is server instance used for game zone.
type nex struct {
	active              bool
	version             string
	config              *Configurer //config from config package
	socketService       websocket.WebsocketService
	commandManager      CommandManager
	builtinEventManager builtinEvent.EventManager
	eventManager        event.EventManager
	grpcServer          *grpc.Server
	httpServer          nxhttp.Server

	hallManager       HallManager
	roomManager       RoomManager
	userManager       UserManager
	requester         nxhttp.Client
	packetHandler     PacketHandler
	middleHandler     MiddleHandler
	updateManager     UpdateManager
	logger            nxlog.Logger
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
	nl, err := nxlog.NewLogger()
	if err != nil {
		goLog.Fatalf("NewLogger error:%s\n", err.Error())
		return err
	}

	nx.logger = nl
	nx.logger.SetLevel(nxlog.LevelInfo)

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
			nx.logger.Log(nxlog.LevelError, fmt.Sprintf("NewWebsocketService error:%s\n", err.Error()))
			return err
		}
		nx.socketService = so
		nx.logger.Log(nxlog.LevelInfo, "NewWebsocketService success\n")
	}

	ph, err := NewPacketHandler()
	if err != nil {
		nx.logger.Log(nxlog.LevelError, fmt.Sprintf("NewPacketHandler error:%s\n", err.Error()))
		return err
	}
	nx.packetHandler = ph
	nx.logger.Log(nxlog.LevelInfo, "NewPacketHandler success\n")

	mh, err := NewMiddleHandler()
	if err != nil {
		nx.logger.Log(nxlog.LevelError, fmt.Sprintf("NewMiddleHandler error:%s\n", err.Error()))
		return err
	}
	nx.middleHandler = mh
	nx.logger.Log(nxlog.LevelInfo, "NewMiddleHandler success\n")

	ch, err := NewCommandManager()
	if err != nil {
		nx.logger.Log(nxlog.LevelError, fmt.Sprintf("NewCommandManager error:%s\n", err.Error()))
		return err
	}
	nx.commandManager = ch
	nx.logger.Log(nxlog.LevelInfo, "NewCommandManager success\n")

	em, err := builtinEvent.NewManager()
	if err != nil {
		nx.logger.Log(nxlog.LevelError, fmt.Sprintf("NewBuiltinEventManager error:%s\n", err.Error()))
		return err
	}
	nx.builtinEventManager = em
	nx.logger.Log(nxlog.LevelInfo, "NewBuiltinEventManager success\n")

	eem, err := event.NewManager()
	if err != nil {
		nx.logger.Log(nxlog.LevelError, fmt.Sprintf("NewEventManager error:%s\n", err.Error()))
		return err
	}
	nx.eventManager = eem
	nx.logger.Log(nxlog.LevelInfo, "NewEventManager success")

	httpCAddr, httpCPort := nx.config.HttpClientAddress()
	re, err := nxhttp.NewClient(httpCAddr,
		httpCPort,
		nx.config.HttpClientTCPConnectTimeoutSecond(),
		nx.config.HttpClientHandshakeTimeoutSecond(),
		nx.config.HttpClientRequestTimeoutSecond())
	if err != nil {
		nx.logger.Log(nxlog.LevelError, fmt.Sprintf("NewRequester error:%s\n", err.Error()))
		return err
	}
	nx.requester = re
	nx.logger.Log(nxlog.LevelInfo, "NewRequester success\n")

	hm, err := NewHallManager()
	if err != nil {
		nx.logger.Log(nxlog.LevelError, fmt.Sprintf("NewHallManager error:%s\n", err.Error()))
		return err
	}
	nx.hallManager = hm
	nx.logger.Log(nxlog.LevelInfo, "NewHallManager success\n")

	rm, err := NewRoomManager()
	if err != nil {
		nx.logger.Log(nxlog.LevelError, fmt.Sprintf("NewRoomManager error:%s\n", err.Error()))
		return err
	}
	nx.roomManager = rm
	nx.logger.Log(nxlog.LevelInfo, "NewRoomManager success\n")

	um, err := NewUserManager()
	if err != nil {
		nx.logger.Log(nxlog.LevelError, fmt.Sprintf("NewUserManager error:%s\n", err.Error()))
		return err
	}

	nx.userManager = um
	nx.logger.Log(nxlog.LevelInfo, "NewUserManager success\n")

	upm, err := NewUpdateManager()
	if err != nil {
		nx.logger.Log(nxlog.LevelError, fmt.Sprintf("NewUpdateManager error:%s\n", err.Error()))
		return err
	}
	nx.updateManager = upm
	nx.logger.Log(nxlog.LevelInfo, "NewUpdateManager success\n")
	return nil
}

//Start starts listening
func (nx *nex) Start() {
	nx.receivePacketChan = nx.socketService.GetReceivePacketChan()
	nx.lostConnIDChan = nx.socketService.GetLostConnIDChan()
	nx.socketService.Start() //go accept
	nx.logger.Log(nxlog.LevelInfo, "Nex started...\n")
	nx.mainFlow()
	nx.handleCtrlC()
}

//Stop user manager process.
func (nx *nex) Stop() {
	nx.active = false
	nx.socketService.Stop()
	nx.logger.Log(nxlog.LevelInfo, "Nex stopped")
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

// Followings are for builtin event
func (nx *nex) RegisterBuiltinEventProcessor(code builtinEvent.EventNO, p builtinEvent.EventProcessor) error {
	return nx.builtinEventManager.RegisterProcessor(code, p)
}
func (nx *nex) UnRegisterBuiltinEventProcessor(code builtinEvent.EventNO) {
	nx.builtinEventManager.UnRegisterProcessor(code)
}

// Followings are for event
func (nx *nex) RegisterEventProcessor(event string, p event.EventProcessor) error {
	return nx.eventManager.RegisterProcessor(event, p)
}
func (nx *nex) UnRegisterEventProcessor(event string) {
	nx.eventManager.UnRegisterProcessor(event)
}
func (nx *nex) DispatchEvent(event string) {
	nx.eventManager.DispatchEvent(event)
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
func (nx *nex) GetLogger() nxlog.Logger {
	return nx.logger
}
func (nx *nex) GetConfig() *Configurer {
	return nx.config
}

func (nx *nex) GetHallManager() HallManager {
	return nx.hallManager
}

func (nx *nex) GetRoomManager() RoomManager {
	return nx.roomManager
}

//http server
func (nx *nex) StartHTTPServer(params *nxhttp.ServerParameters) error {
	if nx.httpServer == nil {

		s, err := nxhttp.NewServer(params)
		if err != nil {
			return err
		}
		nx.httpServer = s

	}
	return nx.httpServer.Start()
}
func (nx *nex) StopHTTPServer() {
	nx.httpServer.Stop()
}
func (nx *nex) GetHTTPServerRequestRealAddr(r *http.Request) (string, error) {
	if nx.httpServer == nil {
		return "", errors.New("HTTP server not initialed")
	}
	return nx.httpServer.GetRealAddr(r), nil
}
func (nx *nex) RegisterHTTPServerHandler(pattern string, handler http.Handler) {
	nx.httpServer.RegisterHandler(pattern, handler)
}

//http client
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

//user
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
func (nx *nex) StartGRPCServer(registerServerFunc interface{}, serverStruct interface{}) error {
	addr, port := nx.config.RPCServerAddress()
	nx.grpcServer = nxRPC.NewServer(addr, port, nx.logger)
	return nxRPC.StartServer(registerServerFunc, serverStruct)
}

//registerServerFunc is grpc protobuf serivce register function.
//example pb.RegisterGreeterServer(s, &server{}),  s==*grpc.Server, &server{}=serverStruct{} contain rpc methods.
//func (nx *nex) StartGRPCServer(addr, port string, registerServerFunc interface{}, serverStruct interface{}) {
//
//	go func() {
//		defer func() {
//			if r := recover(); r != nil {
//				nx.logger.Log(nxlog.LevelPanic, fmt.Sprintf("StartGRPCServer panic=%v", r))
//				return
//			}
//
//		}()
//
//		if strings.Index(reflect.TypeOf(registerServerFunc).String(), "func") != 0 {
//			nx.logger.Log(nxlog.LevelError, fmt.Sprintf("StartGRPCServer fn is not function"))
//			return
//		}
//
//		//check
//		if nx.grpcServer == nil {
//			nx.grpcServer = grpc.NewServer()
//		}
//
//		//server listen
//		addr, _ := net.ResolveTCPAddr("tcp", addr+":"+port)
//		lis, err := net.ListenTCP("tcp", addr)
//		if err != nil {
//			nx.logger.Log(nxlog.LevelError, fmt.Sprintf("StartGRPCServer listen error=%s", err.Error()))
//			return
//		}
//
//		value := reflect.ValueOf(registerServerFunc)
//		value.Call([]reflect.Value{reflect.ValueOf(nx.grpcServer), reflect.ValueOf(serverStruct)})
//
//		if err := nx.grpcServer.Serve(lis); err != nil {
//			nx.logger.Log(nxlog.LevelError, fmt.Sprintf("StartGRPCServer serve error=%s", err.Error()))
//			return
//		}
//
//	}()
//}

func (nx *nex) GetGRPCClient(addr, port string, newClientFunc interface{}) (interface{}, error) {
	return nxRPC.GetGRPCClient(addr, port, newClientFunc)
}

//func (nx *nex) GetGRPCClient(addr, port string, newClientFunc interface{}) (interface{}, error) {
//
//	conn, err := grpc.Dial(addr+":"+port, grpc.WithInsecure())
//	if err != nil {
//		return nil, err
//	}
//
//	value := reflect.ValueOf(newClientFunc)
//	if value.Kind().String() != "func" {
//		return nil, errors.New(" passed newClientFunc is no  a function")
//	}
//	cli := value.Call([]reflect.Value{reflect.ValueOf(conn)}) //retrun []Value
//
//	//value to original type
//	//value.Interface().(oriType)
//	return cli[0].Interface(), nil //value to interface
//}

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
			nx.logger.Log(nxlog.LevelPanic, fmt.Sprintf("mainFlow panic:%v", r))
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
					nx.logger.LogFile(nxlog.LevelError, fmt.Sprintf("nex acketHandler.Handle err=%s\n", err.Error()))
					nx.socketService.Disconnect(connUUID)
					continue
				}

				cmd, err := nx.middleHandler.PacketToCommand(dataBody)

				if err != nil {
					nx.logger.LogFile(nxlog.LevelError, fmt.Sprintf("nex middleHandler.PacketToCommand err=%s\n", err.Error()))
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

					nx.builtinEventManager.RunProcessor(&builtinEvent.EventObject{
						Code: builtinEvent.EventUserLost,
						User: user,
					})
				}

			}

		}
	}
}
