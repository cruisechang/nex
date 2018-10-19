package event

import (
	"sync"
	"errors"
	"github.com/cruisechang/nex/entity"
)


type EventObject struct {
	Event string
	User  entity.User
}

//NotifyProcessor is a interface
//user implements this interface to create notify processor
//and register the process to nex
//when notify fired , processor will be executed
type EventProcessor interface {
	Run(obj *EventObject) error
}

type EventManager interface {
	RegisterProcessor(event string, p EventProcessor) error
	UnRegisterProcessor(event string)
	DispatchEvent(event string) error
}

type eventManager struct {
	//Processor container, use Command name as key
	processorTable map[string]EventProcessor
	mutex          *sync.RWMutex
}

func NewManager() (EventManager, error) {

	c := &eventManager{
		processorTable: make(map[string]EventProcessor),
		mutex:          &sync.RWMutex{},
	}
	return c, nil
}

func (cs *eventManager) RegisterProcessor(event string, p EventProcessor) error {
	if cs.containProcessor(event) {
		return errors.New("processor already in table")
	}

	cs.mutex.Lock()
	defer cs.mutex.Unlock()
	cs.processorTable[event] = p
	return nil
}

func (cs *eventManager) UnRegisterProcessor(event string) {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()
	delete(cs.processorTable, event)
}

func (cs *eventManager) containProcessor(event string) bool {
	cs.mutex.RLock()
	defer cs.mutex.RUnlock()
	_, ok := cs.processorTable[event]
	return ok
}
func (cs *eventManager) DispatchEvent(event string) error {

	p, ok := cs.processorTable[event]
	if !ok {
		return errors.New("processor not found")
	}

	return p.Run(&EventObject{event, nil})
}

