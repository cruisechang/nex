package entity

import (
	"errors"
	"sync"
)

type Room interface {
	ID() int //get userID when create user
	Name() string
	Active() int
	SetActive(value int)
	Type() int
	SetType(value int)

	Status() int
	SetStatus(v int)

	StatusStart()int64
	SetStatusStart(v int64)

	HallID() int
	setHallID(hallID int)

	AddUser(user User) error
	RemoveUser(user User)
	GetUsers() []User

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

type room struct {
	rwMutex sync.RWMutex
	roomID  int
	name    string
	//active indicates whether hall is active, 0=not, 1=active
	active int
	//typ is room type
	tpy    int
	hallID int
	//users contains user ids in this room
	//status is indicates the status of the room,
	//user can define the value of you own
	//default value is -1
	status int
	statusStart int64
	users  []User

	//string varialbe table
	strVarTable map[string]string
	//float32 varialbe table
	floatVarTable map[string]float32
	//int variable table
	intVarTable   map[string]int
	int64VarTable map[string]int64
	//int variable table
	interfaceVarTable map[string]interface{}
}

func NewRoom(roomID int, name string) Room {

	return &room{
		rwMutex:           sync.RWMutex{},
		roomID:            roomID,
		name:              name,
		active:            1,
		tpy:               -1,
		hallID:            -1,
		status:            -1,
		strVarTable:       make(map[string]string),
		floatVarTable:     make(map[string]float32),
		intVarTable:       make(map[string]int),
		interfaceVarTable: make(map[string]interface{}),

		users: []User{},
	}
}

func (r *room) Name() string {
	return r.name
}
func (r *room) ID() int {
	return r.roomID
}

func (r *room) AddUser(user User) error {
	r.rwMutex.Lock()
	defer r.rwMutex.Unlock()

	for _, v := range r.users {
		if v.UserID() == user.UserID() {
			return errors.New("user already in table")
		}
	}

	user.setRoomID(r.roomID)
	r.users = append(r.users, user)

	return nil
}
func (r *room) RemoveUser(user User) {
	r.rwMutex.Lock()
	defer r.rwMutex.Unlock()

	for i, v := range r.users {
		if v.UserID() == user.UserID() {
			user.setRoomID(-1)
			r.users = append(r.users[:i], r.users[i+1:]...)
			break
		}
	}
}
func (r *room) GetUsers() []User {
	r.rwMutex.RLock()
	defer r.rwMutex.RUnlock()

	//tmp := make([]int, len(r.users))
	//copy(tmp, r.users)
	//return tmp
	return r.users
}

func (r *room) Active() int {
	return r.active
}
func (r *room) SetActive(v int) {
	r.active = v
}
func (r *room) Type() int {
	return r.tpy
}
func (r *room) SetType(v int) {
	r.tpy = v
}

func (r *room) HallID() int {
	return r.hallID
}

func (r *room) setHallID(v int) {
	r.hallID = v
}

func (r *room) Status() int {
	return r.status
}
func (r *room) SetStatus(v int) {
	r.status = v
}

func (r *room) StatusStart() int64 {
	return r.statusStart
}
func (r *room) SetStatusStart(v int64) {
	r.statusStart = v
}

func (r *room) SetStringVariable(key, value string) {
	r.strVarTable[key] = value
}
func (r *room) GetStringVariable(key string) (string, error) {
	v, ok := r.strVarTable[key]
	if ok {
		return v, nil
	}
	return "", errors.New("variable not found")
}

func (r *room) SetFloatVariable(key string, value float32) {
	r.floatVarTable[key] = value
}
func (r *room) GetFloatVariable(key string) (float32, error) {
	v, ok := r.floatVarTable[key]
	if ok {
		return v, nil
	}
	return 0, errors.New("variable not found")
}

func (r *room) SetIntVariable(key string, value int) {
	r.intVarTable[key] = value
}
func (r *room) GetIntVariable(key string) (int, error) {
	v, ok := r.intVarTable[key]
	if ok {
		return v, nil
	}
	return 0, errors.New("variable not found")
}

func (r *room) SetInt64Variable(key string, value int64) {
	r.int64VarTable[key] = value
}

func (r *room) GetInt64Variable(key string) (int64, error) {
	v, ok := r.int64VarTable[key]
	if ok {
		return v, nil
	}
	return 0, errors.New("variable not found")
}

func (r *room) SetInterfaceVariable(key string, value interface{}) {
	r.interfaceVarTable[key] = value
}
func (r *room) GetInterfaceVariable(key string) (interface{}, error) {
	v, ok := r.interfaceVarTable[key]
	if ok {
		return v, nil
	}
	return nil, errors.New("variable not found")
}
