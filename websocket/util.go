package websocket

import (
	"encoding/binary"
	"errors"

)


func checkPacketLength(readBuff []byte,packetHeadLen, packetTypeLen int)error{
	//read from gorilla websocket conn


	//make headBuff
	headBuff := make([]byte, packetHeadLen)


	//get head
	copy(headBuff[:], readBuff[0:packetHeadLen])

	if len(headBuff) < 4 {
		return  errors.New("byteToInt32 data length not enougth")
	}

	dataLen := int(binary.LittleEndian.Uint32(headBuff))

	//get data len from head
	//dataLen, err := binaryUtil.ByteToIntByLittleEndian(headBuff)

	//if err != nil {
	//	return  errors.New("headBuff convert to data len error")
	//}

	//check if data len correct
	if len(readBuff) != dataLen+packetTypeLen+packetHeadLen {
		return errors.New("readBuff len not equal headBuff to data len")
	}

	return  nil
}
//ReadPacket from conn
//packetHeadLen is packet head byte number usually 4
//packetTypeLen is packet type byte number usuall 1
//return head + type + body ([]byte)
/*
func readPacket(conn *gorillaWebsocket.Conn, packetHeadLen, packetTypeLen int) ([]byte, error) {

	//read from gorilla websocket conn
	_, readBuff, err := conn.ReadMessage()

	if err != nil {
		common.PrintColor(common.ColorGreen, readPacket, "read error:"+err.Error())
		return nil, err
	}

	//make headBuff
	headBuff := make([]byte, packetHeadLen)


	//get head
	copy(headBuff[:], readBuff[0:packetHeadLen])

	if len(headBuff) < 4 {
		return nil, errors.New("byteToInt32 data length not enougth")
	}

	dataLen:=int(binary.LittleEndian.Uint32(headBuff))

	//check if data len correct
	if len(readBuff) != dataLen+packetTypeLen+packetHeadLen {
		common.PrintColor(common.ColorGreen, readPacket, "head data num !=read data num:"+strconv.Itoa(dataLen)+"/"+strconv.Itoa(len(readBuff)-+packetTypeLen-packetHeadLen))
		return nil, errors.New("readBuff len not equal headBuff to data len")
	}

	return readBuff, nil

}
*/

//ParsePacket parse packet from reader
//Return data length in byte,dataType,data,error
func parsePacket(bytes []byte, packetAllDataLengthMin, packetAllDataLengthMax int) (int, int, []byte, error) {

	if len(bytes) < (packetAllDataLengthMin) {
		return 0, 0, bytes, errors.New("packet length not enough")
	}

	bodyLen:=int(binary.LittleEndian.Uint32(bytes[0:4]))


	//Check data length
	//bodyLen, _ := binaryUtil.ByteToIntByLittleEndian(bytes[0:4])

	if bodyLen > packetAllDataLengthMax {
		return 0, 0, bytes, errors.New("packet data length exceed max")
	} else if bodyLen != len(bytes)-5 {
		return 0, 0, nil, errors.New("Packet data length not match")
	}

	//Check data type
	typeBytes := make([]byte, 2)
	//高位元補0
	typeBytes[0] = 0
	typeBytes[1] = bytes[4]
	dataType:= int(binary.BigEndian.Uint16(typeBytes))

	bodyByte:=bytes[5:]
	//dataType, _ := binaryUtil.SingleByteToIntByBigEndian(bytes[4])
	return bodyLen, dataType, bodyByte, nil
}

//ComposePacket add head and type to packet
//return a complete packet
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
