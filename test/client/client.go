package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	gorillaWebsocket "github.com/gorilla/websocket"
	"github.com/cruisechang/nex"
	"encoding/json"
	"encoding/binary"
	"encoding/base64"
)

var addr = flag.String("addr", "localhost:11000", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: ""}
	log.Printf("connecting to %s", u.String())

	//connect
	c, _, err := gorillaWebsocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()




	//wait for 3 sec
	time.Sleep(time.Duration(time.Second * 3))


	//login command
	//c.WriteMessage(gorillaWebsocket.BinaryMessage, createCompletePacket())
	c.WriteMessage(gorillaWebsocket.TextMessage, createTextMessage())


	go func() {

		//recover panic
		defer func() {
			if r := recover(); r != nil {
				c.Close()
				c=nil
			}
		}()

		for {

			_, resBuff, _ := c.ReadMessage()
			b,_:=base64.StdEncoding.DecodeString(string(resBuff))

			cmd := &nex.Command{}

			unMarshalErr:=json.Unmarshal(b,cmd)
			if unMarshalErr!=nil{
				log.Printf("unmarshal error:%s",unMarshalErr.Error())
			}else{
				switch cmd.Command {
				case "Login":
				case "Heartbeat":
				case "UserInfo":
				case "GameInfo":

				}
			}


		}
	}()



	/*
	done := make(chan struct{})

	//ticker
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			err := c.WriteMessage(gorillaWebsocket.TextMessage, []byte(t.String()))
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(gorillaWebsocket.CloseMessage, gorillaWebsocket.FormatCloseMessage(gorillaWebsocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
	*/

}

func createLoginMessage(cmd string,) []byte {
	d, _ := createLoginCommandData("slot0", "sessionID","zh-hans")

	c, _ := createCommand(0, 1, "Login", d)

	ce, _ := json.Marshal(c)

	str:=base64.StdEncoding.EncodeToString(ce)

	return []byte(str)

}

func createTextMessage() []byte {
	d, _ := createLoginCommandData("slot0", "sessionID","zh-hans")

	c, _ := createCommand(0, 1, "Login", d)

	ce, _ := json.Marshal(c)

	str:=base64.StdEncoding.EncodeToString(ce)

	return []byte(str)

}

func createCompletePacket() []byte {
	d, _ := createLoginCommandData("gameId", "sessionID","lang")
	c, _ := createCommand(0, 1, "Login", d)
	ce, _ := json.Marshal(c)
	return composePacket(0, ce)

}
func composePacket(packetType uint16, packetBody []byte) []byte {
	pl := uint32(len(packetBody))

	headBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(headBytes, pl)

	typeBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(typeBytes, packetType)

	res := make([]byte, pl+4+1)
	copy(res[0:4], headBytes[:])
	copy(res[4:5], typeBytes[1:])
	copy(res[5:], packetBody[:])

	return res

}

func createCommand(code int, step int, cmdName string, data string) (*nex.Command, error) {

	cmd := &nex.Command{
		Code:    code,
		Step:    step,
		Command: cmdName,
		Data:    data,
	}
	return cmd, nil
}

func createLoginCommandData(gameID, sessinID ,lang string) (string, error) {

	l := LoginCmdData{
		GameId:    gameID,
		SessionId: sessinID,
		Lang:lang,
	}
	d := []LoginCmdData{l}

	b, err := json.Marshal(d)
	if err != nil {
		return "", err
	}

	//[]byte encode to base64 string
	return base64.StdEncoding.EncodeToString(b), nil
}

type LoginCmdData struct {
	GameId    string
	SessionId string
	Lang string
}
