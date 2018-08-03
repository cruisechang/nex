package nex

import (
	"sync"
	"errors"
)

type EventManager interface {
	RunProcessor(object *EventObject) error
	RegisterProcessor(code uint16, p EventProcessor) error
	UnRegisterProcessor(code uint16)
	DispatchEvent(code uint16) error
}

type eventManager struct {
	//Processor container, use Command name as key
	processorTable map[uint16]EventProcessor
	mutex          *sync.Mutex
}

func NewEventManager() (EventManager, error) {

	c := &eventManager{
		processorTable: make(map[uint16]EventProcessor),
		mutex:          &sync.Mutex{},
	}
	return c, nil
}

func (cs *eventManager) RegisterProcessor(code uint16, p EventProcessor) error {
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
		cs.processorTable[code] = p
		cs.mutex.Unlock()
		return nil
	}
	return errors.New("code is not illegal")
}

func (cs *eventManager) UnRegisterProcessor(code uint16)  {
		cs.mutex.Lock()
		delete(cs.processorTable,code)
		cs.mutex.Unlock()
}

func (cs *eventManager) containProcessor(code uint16) bool {
	cs.mutex.Lock()
	_, ok := cs.processorTable[code]
	cs.mutex.Unlock()
	return ok
}
func (cs *eventManager) DispatchEvent(code uint16) error {

	p, ok := cs.processorTable[code]
	if !ok {
		return errors.New("processor not found")
	}

	return p.Run(&EventObject{code,nil})
}

func (cs *eventManager) RunProcessor(obj *EventObject) error {

	p, ok := cs.processorTable[obj.Code]
	if !ok {
		return errors.New("processor not found")
	}

	return p.Run(obj)
}
