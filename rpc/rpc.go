package rpc

import (
	"errors"
	"fmt"
	"net"
	"reflect"
	"strings"

	"github.com/cruisechang/nex/log"
	"google.golang.org/grpc"
)

var (
	server     *grpc.Server
	svrAddress string
	svrPort    string
	logger     log.Logger
)

func NewServer(address, port string, log log.Logger) *grpc.Server {
	server = grpc.NewServer()
	svrAddress = address
	svrPort = port
	logger = log
	return server
}
func StartServer(registerServerFunc interface{}, serverStruct interface{}) error {

	if strings.Index(reflect.TypeOf(registerServerFunc).String(), "func") != 0 {
		return errors.New("StartGRPCServer fn is not function")
	}

	//server listen
	addr, _ := net.ResolveTCPAddr("tcp", svrAddress+":"+svrPort)
	lis, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return errors.New(fmt.Sprintf("StartGRPCServer listen error=%s", err.Error()))
	}

	value := reflect.ValueOf(registerServerFunc)
	value.Call([]reflect.Value{reflect.ValueOf(server), reflect.ValueOf(serverStruct)})

	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Log(log.LevelPanic, fmt.Sprintf("StartGRPCServer panic=%v", r))
			}
		}()

		if err := server.Serve(lis); err != nil {
			logger.Log(log.LevelError, fmt.Sprintf("StartGRPCServer serve error=%s", err.Error()))
			return
		}
	}()

	return nil
}

func GetGRPCClient(addr, port string, newClientFunc interface{}) (interface{}, error) {

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
