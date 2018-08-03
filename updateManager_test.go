package nex

import (
	"testing"
	"fmt"
	"time"
)

func TestUpdateManager(t *testing.T) {

	cm, _ := NewUpdateManager()

	name := "login"

	np:=newUpdateProcessor(name, 10*time.Second)

	//register
	err := cm.RegisterProcessor(name, np)

	if err != nil {
		t.Fatalf("TestUpdateManager registerProcessor err=%s", err.Error())
	}

	cm.StartProcessor(name)

	//register the same cmd name, got error
	err = cm.RegisterProcessor(name, newUpdateProcessor(name, 1*time.Second))
	if err == nil {
		t.Fatalf("TestUpdateManager registerProcessor should err but not")
	}

	time.Sleep(5 * time.Second)

	cm.StopProcessor(name)

	//unregister
	err=cm.UnRegisterProcessor(name)


	err = cm.StopProcessor(name)
	if err == nil {
		t.Logf("TestUpdateManager StopProcessor should err but not")

	}

}

func newUpdateProcessor(name string,duration time.Duration) UpdateProcessor {
	return &upProcessor{
		name: name,
		stopChan:make(chan bool,1),
		running:false,
		duration:duration,
	}
}

type upProcessor struct {
	name     string
	stopChan chan bool
	duration time.Duration
	running bool
}

func (p *upProcessor) Run() error {
	fmt.Printf("update Processor run!")
	return nil
}

func (p *upProcessor) StopChan() chan bool {
	return p.stopChan
}
func (p *upProcessor)Duration()time.Duration{
	return p.duration
}
func (p *upProcessor)Running()bool{
	return p.running
}
func (p *upProcessor)SetRunning(isRun bool){
	p.running=isRun
}
