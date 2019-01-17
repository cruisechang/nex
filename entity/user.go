package entity

import "errors"

type User interface {
	ConnID() string
	UserID() int
	Account() string
	SetAccount(account string)
	//玩家名
	Name() string
	SetName(name string)
	Status()int
	SetStatus(int)

	SessionID() string
	SetSessionID(sessionID string)

	GameID() string
	SetGameID(gameID string)

	HallID()int
	//used by package
	setHallID(int)

	RoomID()int
	//used by package
	setRoomID(int)




	IsTest() int
	SetIsTest(is int)

	Credit() float32
	SetCredit(credit float32)

	AccessID() string
	SetAccessID(v string)

	//user defined variable
	SetStringVariable(key, value string)
	GetStringVariable(key string) (string, error)

	SetFloatVariable(key string, value float32)
	GetFloatVariable(key string) (float32, error)

	SetIntVariable(key string, value int)

	GetIntVariable(key string) (int, error)

	SetInt64Variable(key string, value int64)
	GetInt64Variable(key string) (int64, error)

	SetInterfaceVariable(key string, value interface{})
	GetInterfaceVariable(key string) (interface{}, error)
}

type user struct {
	connID    string //connection uuid for websocket
	userID    int
	account string
	name      string
	sessionID string
	accessID  string
	isTest    int
	credit    float32
	gameID    string
	//hallID indicates the hall user in,
	//user can define the value of you own
	//default value is -1
	hallID int
	//roomID indicates the room user in,
	//user can define the value of you own
	//default value is -1
	roomID int
	//status indicates the status of the room,
	//user can define the value of you own
	//default value is -1
	status int

	strVarTable   map[string]string  //string varialbe table
	floatVarTable map[string]float32 //float32 varialbe table
	intVarTable   map[string]int     //int variable table
	int64VarTable   map[string]int64     //int variable table
	interfaceVarTable map[string]interface{} //int variable table
}

func NewUser(userID int, connID string) User {

	return &user{
		userID:        userID,
		connID:        connID,
		strVarTable:   make(map[string]string),
		floatVarTable: make(map[string]float32),
		intVarTable:   make(map[string]int),
		int64VarTable:   make(map[string]int64),
		interfaceVarTable: make(map[string]interface{}),
		hallID:-1,
		roomID:-1,
		status:-1,
		name:"default",
	}
}

func (u *user) UserID() int {
	return u.userID
}

func (u *user) ConnID() string {
	return u.connID
}

func (u *user) GameID() string {
	return u.gameID
}
func (u *user) SetGameID(v string) {
	u.gameID = v
}
func (u *user) Account() string {
	return u.account
}
func (u *user) SetAccount(account string) {
	u.account = account
}

func (u *user) Name() string {
	return u.name
}
func (u *user) SetName(name string) {
	u.name = name
}
func (u *user) SessionID() string {
	return u.sessionID
}
func (u *user) SetSessionID(v string) {
	u.sessionID = v
}
func (u *user) HallID() int {
	return u.hallID
}
func (u *user) setHallID(id int) {
	u.hallID = id
}
func (u *user) RoomID() int {
	return u.roomID
}
func (u *user) setRoomID(id int) {
	u.roomID = id
}

func (u *user) Status() int {
	return u.status
}
func (u *user) SetStatus(v int) {
	u.status = v
}

func (u *user) IsTest() int {
	return u.isTest
}
func (u *user) SetIsTest(v int) {
	u.isTest = v
}
func (u *user) Credit() float32 {
	return u.credit;
}
func (u *user) SetCredit(v float32) {
	u.credit = v
}
func (u *user) AccessID() string {
	return u.accessID
}
func (u *user) SetAccessID(v string) {
	u.accessID = v
}

func (u *user) SetStringVariable(key, value string) {
	u.strVarTable[key] = value
}
func (u *user) GetStringVariable(key string) (string, error) {
	v, ok := u.strVarTable[key]
	if ok {
		return v, nil
	}
	return "", errors.New("variable not found")
}

func (u *user) SetFloatVariable(key string, value float32) {
	u.floatVarTable[key] = value
}
func (u *user) GetFloatVariable(key string) (float32, error) {
	v, ok := u.floatVarTable[key]
	if ok {
		return v, nil
	}
	return 0, errors.New("variable not found")
}

func (u *user) SetIntVariable(key string, value int) {
	u.intVarTable[key] = value
}
func (u *user) GetIntVariable(key string) (int, error) {
	v, ok := u.intVarTable[key]
	if ok {
		return v, nil
	}
	return 0, errors.New("variable not found")
}

func (u *user) SetInt64Variable(key string, value int64) {
	u.int64VarTable[key] = value
}
func (u *user) GetInt64Variable(key string) (int64, error) {
	v, ok := u.int64VarTable[key]
	if ok {
		return v, nil
	}
	return 0, errors.New("variable not found")
}
func (u *user) SetInterfaceVariable(key string, value interface{}) {
	u.interfaceVarTable[key] = value
}
func (u *user) GetInterfaceVariable(key string) (interface{}, error) {
	v, ok := u.interfaceVarTable[key]
	if ok {
		return v, nil
	}
	return nil, errors.New("variable not found")
}
