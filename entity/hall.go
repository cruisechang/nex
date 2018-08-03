package entity

import (
	"errors"
	"sync"
)

type Hall interface {
	ID() int
	Name() string
	SetName(name string)
	Active() int
	SetActive(value int)
	Type() int
	SetType(value int)
	Status() int
	SetStatus(value int)

	AddRoom(roomID int) error
	RemoveRoom(roomD int)
	GetRooms() []int

	SetStringVariable(key, value string)
	GetStringVariable(key string) (string, error)

	SetFloatVariable(key string, value float32)
	GetFloatVariable(key string) (float32, error)

	SetIntVariable(key string, value int)
	GetIntVariable(key string) (int, error)

	SetInterfaceVariable(key string, value interface{})
	GetInterfaceVariable(key string) (interface{}, error)

	SetIntSliceVariable(key string, value int)
	GetIntSliceVariable(key string) ([]int, error)
}

type hall struct {
	rwMutex sync.RWMutex
	hallID  int
	name    string
	//active indicates whether hall is active
	//0=not, 1=active
	//default value is 1
	active int
	//tpy is indicates the type of the room,
	//user can define the value of you own
	//default value is -1
	tpy     int
	//status is indicates the status of the room,
	//user can define the value of you own
	//default value is -1
	status int
	rooms   []int

	//string varialbe table
	strVarTable map[string]string
	//float32 varialbe table
	floatVarTable map[string]float32
	//int variable table
	intVarTable map[string]int
	//int variable table
	interfaceVarTable map[string]interface{}
	intSliceVarTable  map[string][]int
}

func NewHall(hallID int, name string) Hall {

	return &hall{
		rwMutex:           sync.RWMutex{},
		hallID:            hallID,
		name:              name,
		active:            1,
		tpy:               -1,
		status:            -1,
		rooms:             []int{},
		strVarTable:       make(map[string]string),
		floatVarTable:     make(map[string]float32),
		intVarTable:       make(map[string]int),
		interfaceVarTable: make(map[string]interface{}),
		intSliceVarTable:  make(map[string][]int),
	}
}

func (h *hall) ID() int {
	return h.hallID
}
func (h *hall) Name() string {
	return h.name
}
func (h *hall) SetName(name string) {
	h.name = name
}
func (h *hall) Active() int {
	return h.active
}
func (h *hall) SetActive(v int) {
	h.active = v
}
func (h *hall) Type() int {
	return h.tpy
}
func (h *hall) SetType(v int) {
	h.tpy = v
}
func (h *hall) Status() int {
	return h.status
}
func (h *hall) SetStatus(v int) {
	h.status = v
}

func (h *hall) AddRoom(userID int) error {
	h.rwMutex.Lock()
	defer h.rwMutex.Unlock()

	for _, v := range h.rooms {
		if v == userID {
			return errors.New("user already in table")
		}
	}

	h.rooms = append(h.rooms, userID)
	return nil
}
func (h *hall) RemoveRoom(userID int) {
	h.rwMutex.Lock()
	defer h.rwMutex.Unlock()

	for i, v := range h.rooms {
		if v == userID {
			h.rooms = append(h.rooms[:i], h.rooms[i+1:]...)
			break
		}
	}
}
func (h *hall) GetRooms() []int {
	h.rwMutex.RLock()
	defer h.rwMutex.RUnlock()

	tmp := make([]int, len(h.rooms))
	copy(tmp, h.rooms)
	return tmp
}

func (h *hall) SetStringVariable(key, value string) {
	h.strVarTable[key] = value
}
func (h *hall) GetStringVariable(key string) (string, error) {
	v, ok := h.strVarTable[key]
	if ok {
		return v, nil
	}
	return "", errors.New("variable not found")
}

func (h *hall) SetFloatVariable(key string, value float32) {
	h.floatVarTable[key] = value
}
func (h *hall) GetFloatVariable(key string) (float32, error) {
	v, ok := h.floatVarTable[key]
	if ok {
		return v, nil
	}
	return 0, errors.New("variable not found")
}

func (h *hall) SetIntVariable(key string, value int) {
	h.intVarTable[key] = value
}
func (h *hall) GetIntVariable(key string) (int, error) {
	v, ok := h.intVarTable[key]
	if ok {
		return v, nil
	}
	return 0, errors.New("variable not found")
}

func (h *hall) SetInterfaceVariable(key string, value interface{}) {
	h.interfaceVarTable[key] = value
}
func (h *hall) GetInterfaceVariable(key string) (interface{}, error) {
	v, ok := h.interfaceVarTable[key]
	if ok {
		return v, nil
	}
	return nil, errors.New("variable not found")
}

func (h *hall) SetIntSliceVariable(key string, value int) {
	s := h.intSliceVarTable[key]
	s = append(s, value)
}
func (h *hall) GetIntSliceVariable(key string) ([]int, error) {
	s, ok := h.intSliceVarTable[key]

	if ok {
		return s, nil
	}
	return nil, errors.New("variable not found")
}
