package websocket

import (
	"testing"
	"time"
)

func TestNewConnData(t *testing.T) {

	//conn, _ := gorillaWebsocketUpgrader.Upgrade(nil, nil, nil)
	cd,err:=NewConnData(nil, make(chan string, 10),make(chan *ReceivePacketData, 100),time.Duration(time.Second*10))



	if err != nil {
		t.Error(err)
	}
	cd.Start()
	cd.Stop()

}