package nex

import (
	"testing"
	"fmt"
)


func TestEventManager(t *testing.T) {

	nm, err := NewEventManager()

	if err != nil {
		t.Errorf("TestEventManager NewNotifyManager error:%s\n", err.Error())
	}

	p:=createEventProcessor()

	if err=nm.RegisterProcessor(EventUserLost, p);err!=nil{
		s := "TestEventManager RegisterProcessor error:%s\n"
		t.Errorf(s,err.Error())
	}

	if err=nm.RegisterProcessor(999,p);err==nil{
		s := "TestEventManager RegisterProcessor must error but not :%s\n"
		t.Errorf(s,err.Error())
	}

	n:=&EventObject{
		Code:EventUserLost,
		User:nil,
	}
	if err=nm.RunProcessor(n) ;err!=nil{

		s := "TestEventManager RunProcessor error but not :%s\n"
		t.Errorf(s,err.Error())
	}


	um,_:=NewUserManager()
	u,_:=um.CreateUser("xxxx")

	t.Logf("%v:",u)

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
