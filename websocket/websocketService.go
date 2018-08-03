package websocket

/*
Packet Structure=head+type+body
head 4 bytes
type 1 bytes
body n bytes
minium packet length = 6 bytes
*/

import (
	"errors"
	"net/http"
	"sync"
	"time"

	gorillaWebsocket "github.com/gorilla/websocket"

	"github.com/cruisechang/nex/common"
)

var (
	gorillaWebsocketUpgrader gorillaWebsocket.Upgrader
)

//Service this is a socket manager used for followings purpose:
//Hold connection
//Read packet
//Write packet
//Send packet to chan
type WebsocketService interface {
	Start() error
	Stop() error
	Destroy()

	Send(receiverIDs []string, sentPacket []byte) error
	Disconnect(id string) error
	GetReceivePacketChan() <-chan *ReceivePacketData
	GetLostConnIDChan() <-chan string
}

type websocketService struct {
	active  bool
	destroy bool

	address string
	port    string

	aesKey []byte

	//外部
	//receivePachetChan 收到packet 放這裡， gameServer 會來取
	receivePacketChan chan *ReceivePacketData
	//connID in this chan is for others to get
	lostConnIDChan chan string

	connectDataTable map[string]ConnData
	mutex            *sync.Mutex

	acceptTimeout time.Duration
	aliveTimeout  time.Duration
	httpServer    *http.Server
}

//NewService make a new SocketManager
func NewWebsocketService(addresss, portt string, aesKey []byte,
	connectPoolCapacity int,
	packetChanCapacity int,
	connectionLostChanCapacity int,
	acceptTimeout time.Duration,
	aliveTimeout time.Duration) (WebsocketService, error) {

	// address = addresss
	// port = portt

	svc := &websocketService{
		address:      addresss,
		port:         portt,
		aesKey:       aesKey,
		mutex:        &sync.Mutex{},
		active:       true,
		destroy:      false,
	}

	err := svc.initial(connectPoolCapacity,
		packetChanCapacity,
		connectionLostChanCapacity,
		acceptTimeout,
		aliveTimeout)

	return svc, err
}

//Initial to init package variable
//acceptTimeout is socket accept timout
//aliveTimeout is used to check if connection is alive. If timeout, connection is closed.
func (s *websocketService) initial(connectPoolCapacity int,
	packetChanCapacity int,
	connectionLostChanCapacity int,
	acceptTimeout time.Duration,
	aliveTimeout time.Duration) (reErr error) {

	defer func() {
		if r := recover(); r != nil {
			common.PrintError(s.initial, "panic:", r)
			reErr = errors.New("")
		}

	}()

	s.connectDataTable = make(map[string]ConnData, connectPoolCapacity)
	s.receivePacketChan = make(chan *ReceivePacketData, packetChanCapacity)
	s.lostConnIDChan = make(chan string, connectionLostChanCapacity)
	s.acceptTimeout = acceptTimeout
	s.aliveTimeout = aliveTimeout

	gorillaWebsocketUpgrader = gorillaWebsocket.Upgrader{
		// HandshakeTimeout specifies the duration for the handshake to complete.
		//HandshakeTimeout time.Duration

		// ReadBufferSize and WriteBufferSize specify I/O buffer sizes. If a buffer
		// size is zero, then a default value of 4096 is used. The I/O buffer sizes
		// do not limit the size of the messages that can be sent or received.
		ReadBufferSize:  10240,
		WriteBufferSize: 10240,

		// Subprotocols specifies the server's supported protocols in order of
		// preference. If this field is set, then the Upgrade method negotiates a
		// subprotocol by selecting the first match in this list with a protocol
		// requested by the client.
		//Subprotocols []string

		// Error specifies the function for generating HTTP error responses. If Error
		// is nil, then http.Error is used to generate the HTTP response.
		//Error func(w http.ResponseWriter, r *http.Request, status int, reason error)

		// CheckOrigin returns true if the request Origin header is acceptable. If
		// CheckOrigin is nil, the host in the Origin header must not be set or
		// must match the host of the request.
		//CheckOrigin func(r *http.Request) bool
		//不檢查origin
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		// EnableCompression specify if the server should attempt to negotiate per
		// message compression (RFC 7692). Setting this value to true does not
		// guarantee that compression will be supported. Currently only "no context
		// takeover" modes are supported.
		EnableCompression: false,
	}

	reErr = nil
	return reErr
}

//Start start read packet action and others
func (s *websocketService) Start() error {
	s.active = true

	go s.listen()
	return nil
}

//Stop stop read packet action and others
func (s *websocketService) Stop() error {
	common.PrintInfo(s.Stop, "Websocket stopping...")
	s.active = false
	s.stopConns()
	return nil
}

//Destroy delete all things.
func (s *websocketService) Destroy() {
	s.Stop()
	s.destroy = true
}

//Add add new conn data to pool
func (s *websocketService) addConnData(cd ConnData) error {

	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.connectDataTable[cd.GetConnID()] = cd
	return nil
}

//Update resets conn which is already in conn map
/*
func (s *websocketService) Update(id int, aesKey []byte, conn *gorillaWebsocket.Conn) error {
	//conn.SetLinger(0)
	cd := NewConnData(id, aesKey, conn)

	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.connectDataTable[id] = cd
	return nil
}
*/

//Remove remove conn data in pool
func (s *websocketService) removeConnData(id string) (isFound bool) {
	if ok := s.containConnData(id); !ok {
		return false
	}
	s.mutex.Lock()
	delete(s.connectDataTable, id)
	s.mutex.Unlock()
	return true
}

//StartConn starts conntion action
func (s *websocketService) startConnData(id string) {
	if cd, ok := s.getConnData(id); ok {
		cd.Start()
	}
}

//StopConn stops action of connection
func (s *websocketService) stopConnData(id string) {
	if cd, ok := s.getConnData(id); ok {
		cd.Stop()
	}
}

func (s *websocketService) stopConns() {
	defer s.mutex.Unlock()
	s.mutex.Lock()
	for _, cd := range s.connectDataTable {
		cd.Stop()
	}
}

func (s *websocketService) getConnData(id string) (ConnData, bool) {
	s.mutex.Lock()
	c, ok := s.connectDataTable[id]
	s.mutex.Unlock()
	if ok {
		return c, true
	}
	return nil, false
}

//Contains return if this id is already have
func (s *websocketService) containConnData(id string) bool {
	s.mutex.Lock()
	_, ok := s.connectDataTable[id]
	s.mutex.Unlock()
	return ok
}

//Send send packet to client.
func (s *websocketService) Send(receiverConnIDs []string, sentPacket []byte) error {

	defer func() {
		if r := recover(); r != nil {
			common.PrintError(s.Send, "panic :", r)
		}
	}()

	if s.active && len(receiverConnIDs) > 0 {

		idLen := len(receiverConnIDs)
		for i := 0; i < idLen; i++ {
			cd, ok := s.getConnData(receiverConnIDs[i])
			if ok {
				cd.Send(sentPacket)
			}
		}

		return nil
	}
	return errors.New("Send error")
}

//Disconnect disconnects target user's socket connection.
func (s *websocketService) Disconnect(id string) error {
	if cd, ok := s.getConnData(id); ok {
		cd.Stop()
		// s.disconnectChan <- id
		//不能先close，因為這時conn可能在waiting,且有sendChan
		//只有在最後close的地方才能close
		return nil
	}
	return errors.New("websocket.disconnect errors")
}

//GetPacketChan return packetChan
func (s *websocketService) GetReceivePacketChan() <-chan *ReceivePacketData {
	return s.receivePacketChan
}

//GetConnectionLostChan returns conntionLostChan containing connID
//Others using connid to do his work
func (s *websocketService) GetLostConnIDChan() <-chan string {
	return s.lostConnIDChan
}

//listen wait and listen connect
//accept handshake begin
func (s *websocketService) listen() error {

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handlerAccept)

	s.httpServer = &http.Server{
		Addr:           s.address + ":" + s.port,
		Handler:        mux,
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   20 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	s.httpServer.ListenAndServe()

	return nil
}
func (s *websocketService) handlerAccept(w http.ResponseWriter, r *http.Request) {

	if s.active && !s.destroy {

		common.PrintInfo(s.handlerAccept, "accept connecting")

		//websocket
		conn, err := gorillaWebsocketUpgrader.Upgrade(w, r, nil)

		if err != nil {

			common.PrintInfo(s.handlerAccept, "error", err.Error())

			if conn != nil {
				conn.Close()
			}

			//if HandshakeError
			if _, ok := err.(gorillaWebsocket.HandshakeError); ok {
				//http.Error(w, "Not a websocket handshake", 400)
				common.PrintInfo(s.handlerAccept, "handshakeError", err.Error())
				//此時未有connection
				// if conn != nil {
				// 	conn.Close()
				// }
			}
			return
		}

		//timeout seconds
		common.PrintInfo(s.handlerAccept, "connection pass, newConnData")

		//create new connData for receive and send
		cd, _ := NewConnData(conn, s.lostConnIDChan, s.receivePacketChan, s.aliveTimeout)
		s.addConnData(cd)
		cd.Start()

	}
}
