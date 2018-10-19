package nex

import (
	"github.com/cruisechang/nex/entity"
	"time"
)

type Command struct {
	//Code is success code, you can defined yoursef
	Code int
	//Command is string command name,
	Command string
	//Step is a serial number you can use to check command sequence
	Step int
	//Data is real data which is string format, you can define the content
	Data string
}
type CommandLowerCase struct {
	//Code is success code, you can defined yoursef
	Code int `json:"code"`
	//Command is string command name,
	Command string `json:"command"`
	//Step is a serial number you can use to check command sequence
	Step int `json:"step"`
	//Data is real data which is string format, you can define the content
	Data string `json:"data"`
}

//CommandProcessor is a interface
//user implements this interface to create command processor
//and register the processor  to nex
//when command received , processor will be executed
type CommandProcessor interface {
	Run(obj *CommandObject) error
}

//ProcessorObject is a container
//contains command and user (send this command)
type CommandObject struct {
	Cmd  *Command
	User entity.User
}

/*

Following is for event

 */

//Notify code is defined by nex
//user use this code to receive notify
//const (
//	EventUserLost uint16 = 0
//)
//
//type EventObject struct {
//	Code uint16
//	User entity.User
//}

//NotifyProcessor is a interface
//user implements this interface to create notify processor
//and register the process to nex
//when notify fired , processor will be executed
//type EventProcessor interface {
//	Run(obj *EventObject) error
//}

/*
 updateProcessor is a scheduled processon
 */

type UpdateObject struct {
	Name string
	User entity.User
}

type UpdateProcessor interface {
	Run() error
	StopChan() chan bool     //get a chan used for stop routine
	Duration() time.Duration //get update duration
	Running() bool           //return if processor is running
	SetRunning(bool)         //return if processor is running
}
