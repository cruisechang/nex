package nex

import (
	"testing"
	"fmt"
)


func TestCommandManager(t *testing.T) {

	cm, _ := NewCommandManager()


	cmdStr:="login"


	obj := CommandObject{
		Cmd:  &Command{
			Code: 0,
			Command: cmdStr,
			Step: 0,
			Data: "data",
		},
		User: nil,
	}

	//register
	cm.RegisterProcessor(cmdStr, newCmdProcessor())

	cm.RunProcessor(&obj)

	//register the same cmd name
	cm.RegisterProcessor(cmdStr, newCmdProcessor())

	cm.RunProcessor(&obj)
	//unregister
	cm.UnRegisterProcessor(cmdStr)
}

func newCmdProcessor() CmdProcessor {
	return &cmdProcessor{}
}

type CmdProcessor interface {
	Run(obj *CommandObject) error
}

type cmdProcessor struct {
}

func (p *cmdProcessor) Run(obj *CommandObject) error {
	fmt.Printf("CommandProcessor is running!\n")
	return nil
}

