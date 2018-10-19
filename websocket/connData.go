package websocket

import (
	"strconv"
	"sync"
	"time"

	gorillaWebsocket "github.com/gorilla/websocket"
	"github.com/cruisechang/util/uuid"
	"github.com/cruisechang/util/security"
	"github.com/cruisechang/nex/common"
	"fmt"
)

//ConnData is connection object
//which do receive and send data
//using gorilla websocket conn
type ConnData interface {
	Start()
	Stop()
	GetConnID() string
	Send(packetBody []byte)
}

type connData struct {
	active            bool
	conn              *gorillaWebsocket.Conn
	aesKey            []byte
	sendChan          chan []byte //chan contains data to send
	receiveDone       bool        //is receive done
	sendDone          bool        // is send done
	mutex             *sync.Mutex
	connID            string                  //conn uuid for get specific connData
	connectLostChan   chan string             //contains lost connID  send to websocket
	receivePacketChan chan *ReceivePacketData //send to websocket
	aliveTimeout      time.Duration

}

//newConnData returns a new connData.
func NewConnData(co *gorillaWebsocket.Conn,
	connecLostChan chan string,
	receivePacketChan chan *ReceivePacketData,
	aliveTimeout time.Duration) (ConnData, error) {

	conUUID, _ := uuid.GetV4()
	aesKey, err := security.GenerateAES256RandomBytes()

	if err != nil {
		aesKey = []byte("qwertyuiasdfghjk12345678!@#$%^&*")
	}

	c := &connData{
		active:            false,
		aesKey:            aesKey,
		conn:              co,
		sendChan:          make(chan []byte, 300),
		mutex:             &sync.Mutex{},
		sendDone:          false,
		receiveDone:       false,
		connID:            conUUID,
		connectLostChan:   connecLostChan,
		receivePacketChan: receivePacketChan,
		aliveTimeout:      aliveTimeout,
	}

	return c, nil
}

func (cd *connData) Start() {
	cd.active = true
	go cd.receive()
	go cd.send() //send from sendChan
}

func (cd *connData) Stop() {
	cd.active = false
	//set read timeout 0
	if cd.conn != nil {
		cd.conn.SetReadDeadline(time.Now().Add(time.Millisecond * 1))
	}
}
func (cd *connData) GetConnID() string {
	return cd.connID
}

func (cd *connData) closeSendChan() {
	close(cd.sendChan)
}

//after accept and after add user
//Close conn 1. aliveTimeout 2.readPacket error
func (cd *connData) receive() {

	//recover panic
	defer func(connID string) {
		if r := recover(); r != nil {
			common.PrintError(cd.receive, "panic", r)
		}
		cd.closeSendChan() //只在這裡close
		cd.receiveDone = true
		cd.active = false
		if nil != cd.conn && cd.sendDone {
			cd.connectLostChan <- connID // 只在這裡加入
			cd.conn.Close()
			cd.conn = nil
		}

	}(cd.connID)

	for cd != nil && cd.active && cd.connID != "" {

		//Set aliveTimeout, add duration from now
		if err := cd.conn.SetReadDeadline(time.Now().Add(cd.aliveTimeout)); err != nil {
			//if err := cd.conn.SetReadDeadline(time.Now().Add(time.Second * 10)); err != nil {
			common.PrintInfo(cd.receive, "conn alive timeout ...")

			cd.active = false
			//close timeout
			cd.conn.SetReadDeadline(time.Time{})
			return
		}

		//read from gorilla websocket conn
		msgType, resBuff, err2 := cd.conn.ReadMessage()


		//After readPacket,maybe timeout already, so return
		if !cd.active {
			common.PrintInfo(cd.receive, "After readPacket but aliveTimeout")
			cd.conn.SetReadDeadline(time.Time{})
			return
		}

		//disable readDeadline time.Time{} =zero
		//After readPacket, diable aliveTimeout
		if deadLineErr := cd.conn.SetReadDeadline(time.Time{}); deadLineErr != nil {
			common.PrintInfo(cd.receive, "Disable conn aliveTimeout error")
			cd.active = false
			return
		}

		//fmt.Printf("receive resBuff=%v \n",resBuff)

		//Check if connection lost or error
		if err2 != nil {
			fmt.Printf("receive error, resBuff=%v \n",resBuff)

			//check close error
			if e, ok := err2.(*gorillaWebsocket.CloseError); ok {
				common.PrintColor(common.ColorGreen, cd.receive, "read close error userID :", strconv.Itoa(-1), "websocket closeError code:", e.Code)
				cd.active = false
				return
			}
			// common.PrintColor(common.ColorGreen, cd.receive, "read error userID :", strconv.Itoa(cd.userID), err.Error())
			cd.active = false
			return
		}

		//text message
		if msgType == gorillaWebsocket.TextMessage {
			cd.receivePacketChan <- &ReceivePacketData{Packet: resBuff, ConnID: cd.connID, MessageType: TextMessage}
		} else {
			cd.receivePacketChan <- &ReceivePacketData{Packet: resBuff, ConnID: cd.connID, MessageType: BinaryMessage}
		}
	}
}

//Send write packet body into sendChan
//and go routine do real send
func (cd *connData) Send(sentPacket []byte) {
	if cd.active && cd.sendChan != nil {
		cd.sendChan <- sentPacket
	}
}

func (cd *connData) send() {
	defer func(connID string) {
		if r := recover(); r != nil {
			common.PrintError(common.ColorGreen, cd.send, "panic :", r)
		}
		cd.active = false
		cd.sendDone = true
		if nil != cd.conn && cd.receiveDone {
			cd.connectLostChan <- connID // 只在這裡加入
			cd.conn.Close()
			cd.conn = nil

		}
	}(cd.connID)

	for nil != cd && cd.active && nil != cd.conn {
		b, ok := <-cd.sendChan

		if ok && len(b) > 0 {

			//check
			if nil != cd && cd.active && nil != cd.conn {

				fmt.Println(string(b))

				//err := cd.conn.WriteMessage(gorillaWebsocket.BinaryMessage, b)
				err := cd.conn.WriteMessage(gorillaWebsocket.TextMessage, b)

				if err != nil {
					cd.active = false
					common.PrintColor(common.ColorGreen, cd.send, "send error userID :", strconv.Itoa(-1), err.Error())
					return
				}
			} else {
				common.PrintColor(common.ColorGreen, cd.send, "sendChan close userID:", strconv.Itoa(-1))
				return
			}

		}

	}
}
