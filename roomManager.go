package nex

import (
	"sync"
	"errors"
	"github.com/cruisechang/nex/entity"
)

type RoomManager interface {
	CreateRoom(roomID int, name string) (entity.Room, error)
	RemoveRoom(roomID int)
	GetRoom(roomID int) (entity.Room, bool)
	GetRooms() []entity.Room
}
type roomManager struct {
	roomTable map[int]entity.Room
	mutex     sync.Mutex
}

func NewRoomManager() (RoomManager, error) {
	return &roomManager{
		roomTable: make(map[int]entity.Room),
		mutex:     sync.Mutex{},
	}, nil
}

//AddUser returns user just created or error
func (rm *roomManager) CreateRoom(id int, name string) (entity.Room, error) {

	if _, ok := rm.roomTable[id]; ok {
		return nil, errors.New("already in table")
	}

	rm.mutex.Lock()
	h := entity.NewRoom(id, name)
	rm.roomTable[id] = h
	rm.mutex.Unlock()

	return h, nil
}

func (rm *roomManager) RemoveRoom(roomID int) {
	rm.mutex.Lock()
	delete(rm.roomTable, roomID)
	rm.mutex.Unlock()
}

func (rm *roomManager) GetRoom(hallID int) (entity.Room, bool) {

	v, ok := rm.roomTable[hallID]
	return v, ok
}

func (rm *roomManager) GetRooms() []entity.Room{

	rooms := []entity.Room{}
	for _, v := range rm.roomTable {
		rooms = append(rooms, v)
	}

	return rooms
}
