package nex

import (
	"github.com/juju/errors"
	"sync"
)

type CommandManager interface {
	RunProcessor(object *CommandObject) error
	RegisterProcessor(cmdName string, p CommandProcessor) error
	UnRegisterProcessor(cmdName string)
	CreateCommand(code int, step int, cmdName string, data string) (*Command, error)
}

type commandManager struct {
	//Processor container, use Command name as key
	processorTable map[string]CommandProcessor
	mutex          *sync.Mutex
}

func NewCommandManager() (CommandManager, error) {

	c := &commandManager{
		processorTable: make(map[string]CommandProcessor),
		mutex:          &sync.Mutex{},
	}
	return c, nil
}

//CreateCommand
//code is success code
//cmdName is a string of command name
//data is base64 string encoded in []byte
func (cs *commandManager) CreateCommand(code int, step int, cmdName string, data string) (*Command, error) {

	cmd := &Command{
		Code:    code,
		Step:    step,
		Command: cmdName,
		Data:    data,
	}
	return cmd, nil
}

func (cs *commandManager) RegisterProcessor(cmdName string, p CommandProcessor) error {
	if cs.containProcessor(cmdName) {
		return errors.New("processor already in table")
	}
	cs.mutex.Lock()
	cs.processorTable[cmdName] = p
	cs.mutex.Unlock()

	return nil
}
func (cs *commandManager) UnRegisterProcessor(cmdName string) {
	delete(cs.processorTable, cmdName)
}

func (cs *commandManager) containProcessor(cmdName string) bool {
	cs.mutex.Lock()
	_, ok := cs.processorTable[cmdName]
	cs.mutex.Unlock()
	return ok
}
func (cs *commandManager) RunProcessor(obj *CommandObject) error {

	p, ok := cs.processorTable[obj.Cmd.Command]
	if !ok {
		return errors.New("process not found")
	}

	return p.Run(obj)
}
