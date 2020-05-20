package builtinEvent

import (
	"errors"
	"github.com/cruisechang/nex/entity"
	"sync"
)

//Notify code is defined by nex
//user use this code to receive notify
const (
//UserLost uint16 = 0
)

type EventNO uint8

const (
	EventUserLost EventNO = 0
)

type EventObject struct {
	Code EventNO
	User entity.User
}

//EventProcessor is a interface
//user implements this interface to create notify processor
//and register the process to nex
//when notify fired , processor will be executed
type EventProcessor interface {
	Run(obj *EventObject) error
}

type EventManager interface {
	RunProcessor(object *EventObject) error
	RegisterProcessor(code EventNO, p EventProcessor) error
	UnRegisterProcessor(code EventNO)
}

type eventManager struct {
	//Processor container, use Command name as key
	processorTable map[EventNO]EventProcessor
	mutex          *sync.RWMutex
}

func NewManager() (EventManager, error) {

	c := &eventManager{
		processorTable: make(map[EventNO]EventProcessor),
		mutex:          &sync.RWMutex{},
	}
	return c, nil
}

func (cs *eventManager) RegisterProcessor(code EventNO, p EventProcessor) error {
	if cs.containProcessor(code) {
		return errors.New("processor already in table")
	}

	codeOK := false
	//check code illigal
	if code == EventUserLost {
		codeOK = true
	}

	if codeOK {
		cs.mutex.Lock()
		defer cs.mutex.Unlock()
		cs.processorTable[code] = p
		return nil
	}
	return errors.New("code is not illegal")
}

func (cs *eventManager) UnRegisterProcessor(code EventNO) {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()
	delete(cs.processorTable, code)
}

func (cs *eventManager) containProcessor(code EventNO) bool {
	cs.mutex.RLock()
	defer cs.mutex.RUnlock()
	_, ok := cs.processorTable[code]
	return ok
}
func (cs *eventManager) RunProcessor(obj *EventObject) error {

	p, ok := cs.processorTable[obj.Code]
	if !ok {
		return errors.New("processor not found")
	}

	return p.Run(obj)
}
