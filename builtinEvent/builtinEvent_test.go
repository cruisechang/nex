package builtinEvent

import (
	"fmt"
	"testing"
)

func TestEventManager(t *testing.T) {

	nm, err := NewManager()

	if err != nil {
		t.Errorf("TestEventManager NewNotifyManager error:%s\n", err.Error())
	}

	p := createEventProcessor()

	if err = nm.RegisterProcessor(EventUserLost, p); err != nil {
		s := "TestEventManager RegisterProcessor error:%s\n"
		t.Errorf(s, err.Error())
	}

	if err = nm.RegisterProcessor(1, p); err == nil {
		s := "TestEventManager RegisterProcessor must error but not :%s\n"
		t.Errorf(s, err.Error())
	}

	n := &EventObject{
		Code: EventUserLost,
		User: nil,
	}
	if err = nm.RunProcessor(n); err != nil {

		s := "TestEventManager RunProcessor error but not :%s\n"
		t.Errorf(s, err.Error())
	}

}

func createEventProcessor() Processor {
	return &processor{}
}

type Processor interface {
	Run(obj *EventObject) error
}

type processor struct {
}

func (p *processor) Run(obj *EventObject) error {
	fmt.Printf("Processor run!")
	return nil
}
