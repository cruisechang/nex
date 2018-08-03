# nex

# package interfaces
import "kumay/gameServer2/interfaces"

### type service

Service is a public api interface,
Plugin using service to call functions in gameServer
Including functions:

##### module

* Command
* Game
* Zone
* Room
* User

##### socket functions 

- SendCommand
- SendCommandByData

##### db functions 

* selectUser
* updateUser

```
func (s *Service) GetUser(userID int) (*module.User, error) 
func (s *Service) ContainsUser(userID int) bool 
func (s *Service) GetRoom(gameID, zoneID string, roomNum int) (*module.Room, error)

// DB SelectUser select  user by userID.
func (s *Service) SelectUser(userID int) ([]map[string]interface{}, error) 

// DB UpdateUser updates user data in DB
func (s *Service) UpdateUser(userID int, args map[string]interface{}) ([]map[string]interface{}, error) 

func (s *Service) SendCommand(gameID, zoneID string, senderID int, receiverIDs []int, cmd *module.Command) error
func (s *Service) SendCommand(gameID, zoneID string, senderID int, receiverIDs []int, cmd *module.Command) error
```

## package module
import "kumay/gameServer2/interfaces/module"

### type CommandInfo
commandInfo contains command name,*Command and *User

```
  type CommandInfo struct {
	   Command string
	   Data    *Command //command struct
	   User    *User    //nil for nouser
}
```

### type CommandContainer

CommandContainer contains conn and command and userID
socket receive command from client
wrap command int CommandContainer pass to gameServer


```
type CommandContainer struct {
	Conn    *websocket.Conn
	Command *CommandCom
	UserID  int
}

```

### type NotifyCommand 

NotifyCommand sent by plugin to gameServer
notify server what happend

```
type NotifyCommand struct {
	Code   int  //code is defined in interfaces.config
	UserID []int
}
```

### type Game

Game is the top module in the structure of game.
The structure of game is like following:

Game->Zone->Room->User
-- Game contains zone
-- Zone contains room
-- Romm contains user(id)

```
func (m *Game) GetZone(zoneID string) (*Zone, error)
func (m *Game) ContainsZone(zoneID string) bool
func (m *Game) GetZoneAll() map[string]*Zone
func (m *Game) CountZone() int

```

### type Zone

```
func (z *Zone) GetRoom(roomNum int) (*Room, error)
func (z *Zone) GetRoomAll() map[int]*Room
func (z *Zone) CountRoom() int
```

### type Room

Room is basic module in game which contains userIDs which are in room players ids.

##### Plugin與server溝通管道都在Room裡

> commandInfoChan

server pass commandInfo to plugin
GetCommandInfoChan()

>notifyCommandChan 

plugin pass notify to server
WriteNotifyCommandChan()

>HTTPParamChan 

server pass http param to plugin 
GetHTTParamChan()

```
var RoomNum int                  
fun (r *Room) IsActive()bol                
func ContainsUser(userID int)bool  
func GetUserIDs()([]int, error)

func GetCommandInfoChan() <-chan *CommandInfo  
//Get commandInfo chan to re

func WriteNotifyCommandChan(cmd *NotifyCommand)
//Plugin writes notify cmd into chan

func GetHTTPParamChan() <-chan map[string]string
//HTTP param chan is for plugin to read http param like grabRedRev....
```

### type User
Contains user data.

```
type User struct {
	GameID     string `json:"WZGameID"` //自己的gameID (serverID)
	ZoneID     string `json:"WZZoneID"`
	RoomNumber int    `json:"WZRoomNumber"`
	SessionID  string `json:"WZSessionID"`

	Step         int     `json:"Step"`
	UserID       int     `json:"UserID"`
	Name         string  `json:"Name"`
	OrderID      int     `json:"OrderID"`
	Credit       float64 `json:"Credit"`
	Coin         float64 `json:"Coin"`
	VIPLevel     int     `json:"VIPLevel"`
	Exp          int     `json:"Exp"`
	UserLevel    int     `json:"UserLevel"`
	Email        string  `json:"Email"`
	Gender       string  `json:"Gender"`
	Birthday     string  `json:"Birthday"`
	QQAccount    string  `json:"QQAccount"`
	IdentityCard string  `json:"IdentityCard"`
	CellPhone    string  `json:"CellPhone"`
	Portrait     string  `json:"Portrait"`
	//不再db query中
	Account  string `json:"Account"`
	GameIDEx int    `json:"GameID"` //外部 gameID from loginEx url
	IP       string `json:"IP"`     //login時client自行傳入的
	IsRobot  bool   `json:"IsRobot"`

	AESKey []byte `json:"AESKey"`
	//Player     *interfaces.Player
}
```