package nex

import (
	"encoding/json"
	"encoding/base64"
)

type MiddleHandler interface {
	PacketToCommand(packet []byte) (*Command, error)
	CommandToPacket(cmd *Command) ([]byte, error)
	CommandToPacketLowercase(cmd *CommandLowerCase) ([]byte, error)
}
type middleHandler struct {
}

func NewMiddleHandler() (MiddleHandler, error) {
	m := &middleHandler{

	}
	return m, nil
}

func (m *middleHandler) PacketToCommand(packet []byte) (*Command, error) {

	baseStr:=string(packet)

	b,err:=base64.StdEncoding.DecodeString(baseStr)

	if err!=nil{
		return nil,err
	}
	c := &Command{}
	err=json.Unmarshal(b, c)

	return c, err
}

func (m *middleHandler) CommandToPacket(cmd *Command) ([]byte, error) {

	c,err:=json.Marshal(cmd)

	if err!=nil{
		return nil,err
	}

	baseStr:=base64.StdEncoding.EncodeToString(c)


	return []byte(baseStr), nil
}
func (m *middleHandler) CommandToPacketLowercase(cmd *CommandLowerCase) ([]byte, error) {

	c,err:=json.Marshal(cmd)

	if err!=nil{
		return nil,err
	}

	baseStr:=base64.StdEncoding.EncodeToString(c)


	return []byte(baseStr), nil
}