package websocket

import (
	"testing"
	"time"
)

var (
	ws WebsocketService
)
func TestNewWebSocket(t *testing.T) {

	so, err := NewWebsocketService("127.0.0.1",
		"9000",
		[]byte("qwertyuiasdfghjk12345678!@#$%^&*"),
		500,
		500,
		500,
		time.Second*time.Duration(20),
		time.Second*time.Duration(20))

	if err != nil {
		t.Error(err)
	}else{
		ws=so
	}


}

func TestStart(t *testing.T){
	if ws!=nil {

		err := ws.Start()
		if err != nil {
			t.Error(err)
		}
	}else{
		t.Error("ws is nil")
	}
}