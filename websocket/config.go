package websocket

const (
	//Packet structure 4bytes head(packet length) +1bytes type +body

	//PacketTypeJSON is const of packet data type json
	PacketTypeJSON int = 0
	//PacketTypeClass is const of packet data type class
	PacketTypeClass int = 1
	//PacketTypeHeartbeat is const of packet data type heartbeat
	PacketTypeHeartbeat int = 2

	packetLengthMin  int = 4 + 1 + 1
	packetLengthMax  int = 40960
	packetHeadLength int = 4
	packetTypeLength int = 1

	//TextMessage the same as gorilla message type
	TextMessage int =1
	//BinaryMessage the same as gorilla message type
	BinaryMessage int =2
)

var (
	//pingPongPacket for connData pingPong
	pingPongPacket = composePacket(uint16(PacketTypeHeartbeat), []byte{1})
)



//ReceivePacketData contain parsed date which is read from connection
//Without packet head
//With data type and packet body
type ReceivePacketData struct {
	Packet     []byte
	ConnID   string
	MessageType int
}

//SendPacketData store receiveIDs and packet to send
type SendPacketData struct {
	receiverIDs []string
	packetBody  []byte
}
