package nex

import (
	"encoding/binary"
	"errors"
	"github.com/cruisechang/goutil/security"
	"github.com/cruisechang/nex/websocket"
)

const (
	packetLengthMin  int = 4 + 1 + 1
	packetLengthMax  int = 40960
	packetHeadLength int = 4
	packetTypeLength int = 1
)

type PacketHandler interface {
	Handle(packetData *websocket.ReceivePacketData) ([]byte, string, error)
	//CheckPacketLength(readBuff []byte)error
	//ParsePacket(bytes []byte) (int, int, []byte, error)
}
type packetHandler struct {
}

func NewPacketHandler() (PacketHandler, error) {

	ph := &packetHandler{
	}

	return ph, nil
}

//Handle checks packet length and parse packet into bodyLen, dataType, body
func (ph *packetHandler) Handle(packetData *websocket.ReceivePacketData) (packetBody []byte, connUUID string, err error) {

	//TextMessage
	if packetData.MessageType == websocket.TextMessage {
		packetBody = packetData.Packet
		return packetBody, packetData.ConnID, nil
	}

	if err = ph.checkPacketLength(packetData.Packet); err != nil {
		return packetData.Packet, packetData.ConnID, err
	}

	_, _, packetBody, err = ph.parsePacket(packetData.Packet)

	if err != nil {
		return packetBody, packetData.ConnID, err
	}

	//decrypted
	//encData, err := gs.packetHandler.DecryptPacket(bodyByte, cd.aesKey)
	//if err != nil {
	//	common.PrintError(cd.receive, "encrypt packet error userID :", strconv.Itoa(-1), err.Error())
	//	common.LogFileError(-1, "websocket encrypt packet error:"+err.Error())
	//	return
	//}

	return packetBody, packetData.ConnID, nil
}
func (ph *packetHandler) checkPacketLength(readBuff []byte) error {
	//read from gorilla websocket conn

	//make headBuff
	headBuff := make([]byte, packetHeadLength)

	//get head
	copy(headBuff[:], readBuff[0:packetHeadLength])

	if len(headBuff) < 4 {
		return errors.New("byteToInt32 data length not enougth")
	}

	dataLen := int(binary.LittleEndian.Uint32(headBuff))

	//check if data len correct
	if len(readBuff) != dataLen+packetTypeLength+packetHeadLength {
		return errors.New("readBuff len not equal headBuff to data len")
	}

	return nil
}

//ParsePacket parse packet from reader
//Return data length in byte,dataType,data,error
func (ph *packetHandler) parsePacket(readBuff []byte) (int, int, []byte, error) {

	if len(readBuff) < (packetLengthMin) {
		return 0, 0, readBuff, errors.New("packet length not enough")
	}

	bodyLen := int(binary.LittleEndian.Uint32(readBuff[0:4]))

	//Check data length
	//bodyLen, _ := binaryUtil.ByteToIntByLittleEndian(bytes[0:4])

	if bodyLen > packetLengthMax {
		return 0, 0, readBuff, errors.New("packet data length exceed max")
	} else if bodyLen != len(readBuff)-5 {
		return 0, 0, nil, errors.New("Packet data length not match")
	}

	//Check data type
	typeBytes := make([]byte, 2)
	//高位元補0
	typeBytes[0] = 0
	typeBytes[1] = readBuff[4]
	dataType := int(binary.BigEndian.Uint16(typeBytes))

	bodyByte := readBuff[5:]
	//dataType, _ := binaryUtil.SingleByteToIntByBigEndian(bytes[4])
	return bodyLen, dataType, bodyByte, nil
}

//ComposePacket add head and type to packet
//return a complete packet
func (ph *packetHandler) composePacket(packetType uint16, packetBody []byte) []byte {
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

//encrypt read packet data
func (ph *packetHandler) encryptPacket(data []byte, key []byte) (resData []byte, resErr error) {
	defer func() {
		if r := recover(); r != nil {
			resData = nil
			resErr = errors.New("encryptPacket panic")
		}
	}()
	resData, resErr = security.AESCbcPkcs7PaddingEncrypt(data, key)
	return
}

//decrypt write packet data
func (ph *packetHandler) decryptPacket(data []byte, key []byte) (resData []byte, resErr error) {

	defer func() {
		if r := recover(); r != nil {
			resData = nil
			resErr = errors.New("encrypt panic")
		}
	}()
	resData, resErr = security.AESCbcPkcs7PaddingDecrypt(data, key)

	return resData, resErr
}
