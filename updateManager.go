package nex

import (
	"sync"
	"errors"
	"time"
	"fmt"
)

type UpdateManager interface {
	RegisterProcessor(name string, p UpdateProcessor) error
	UnRegisterProcessor(name string)error

	StartProcessor(name string) (resErr error)
	StopProcessor(name string) (resErr error)
	IsRunning(name string)(bool,error)
}

type updateManager struct {
	//Processor container, use Command name as key
	processorTable map[string]UpdateProcessor
	mutex          *sync.Mutex
}

func NewUpdateManager() (UpdateManager, error) {

	c := &updateManager{
		processorTable: make(map[string]UpdateProcessor),
		mutex:          &sync.Mutex{},
	}
	return c, nil
}

func (cs *updateManager) RegisterProcessor(name string, p UpdateProcessor) error {
	if cs.containProcessor(name) {
		return errors.New("processor already in table")
	}

	//check if stopChan is initialized
	ch:=p.StopChan()
	if cap(ch) == 0 {
		ch=make(chan bool,1)
	}

	cs.mutex.Lock()
	cs.processorTable[name] = p
	cs.mutex.Unlock()
	return nil
}

/*
 unregister 一定要檢查是否stop 否則可能run 不停
 */
func (cs *updateManager) UnRegisterProcessor(name string) error {

	_, ok := cs.processorTable[name]
	if !ok {
		return errors.New("processor is not in table")
	}
	//
	//if p.Running() {
	//	return errors.New("processor is running")
	//}

	cs.mutex.Lock()
	delete(cs.processorTable, name)
	cs.mutex.Unlock()

	return nil
}

func (cs *updateManager) containProcessor(name string) bool {
	cs.mutex.Lock()
	_, ok := cs.processorTable[name]
	cs.mutex.Unlock()
	return ok
}

func (cs *updateManager) StartProcessor(name string) (resErr error) {

	defer func() {
		if r := recover(); r != nil {

			resErr = errors.New(fmt.Sprintf("updateManager StartProcessor panic:%v", r))
		}
	}()

	processor, ok := cs.processorTable[name]
	if !ok {
		resErr = errors.New("StartProcessor processor not found")
		return
	}

	update(processor)

	return
}

func (cs *updateManager) StopProcessor(name string) (resErr error) {

	defer func() {
		if r := recover(); r != nil {
			resErr = errors.New(fmt.Sprintf("updateManager StopProcessor panic:%v", r))
		}
	}()

	processor, ok := cs.processorTable[name]
	if !ok {
		resErr = errors.New("StopProcessor processor not found")
		return
	}

	ch := processor.StopChan()
	if len(ch) == 0 {
		ch <- true
	}

	return
}

func update(processor UpdateProcessor) chan bool {
	//ticker := time.NewTicker(duration)
	stop := make(chan bool, 1)

	go func() {
		ticker := time.NewTicker(processor.Duration())
		processor.SetRunning(true)

		for {
			select {
			case <-ticker.C:
				processor.Run()
			case <-processor.StopChan():
				processor.SetRunning(false)
				return
			}
		}
	}()

	return stop
}

func (cs *updateManager) IsRunning(name string) (resBool bool, resErr error) {

	defer func() {
		if r := recover(); r != nil {
			resErr = errors.New(fmt.Sprintf("updateManager StopProcessor panic:%v", r))
		}
	}()

	processor, ok := cs.processorTable[name]
	if !ok {
		return false,nil
	}

	return processor.Running(),nil
}