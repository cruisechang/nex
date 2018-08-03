package nex

import (
	"sync"
	"errors"
	"github.com/cruisechang/nex/entity"
)

type HallManager interface {
	CreateHall(hallID int, name string) (entity.Hall, error)
	RemoveHall(hallID int)
	GetHall(hallID int) (entity.Hall, bool)
	GetHalls() []entity.Hall
	ContainHall(hallID int)bool
}
type hallManager struct {
	hallTable map[int]entity.Hall
	mutex     sync.RWMutex
}

func NewHallManager() (HallManager, error) {
	return &hallManager{
		hallTable: make(map[int]entity.Hall),
		mutex:     sync.RWMutex{},
	}, nil
}

//AddUser returns user just created or error
func (hm *hallManager) CreateHall(id int, name string) (entity.Hall, error) {

	if _, ok := hm.hallTable[id]; ok {
		return nil, errors.New("user already in table")
	}

	hm.mutex.Lock()
	defer hm.mutex.Unlock()
	h := entity.NewHall(id, name)
	hm.hallTable[id] = h

	return h, nil
}

func (hm *hallManager) RemoveHall(hallID int) {
	hm.mutex.Lock()
	defer hm.mutex.Unlock()
	delete(hm.hallTable, hallID)
}

func (hm *hallManager) GetHall(hallID int) (entity.Hall, bool) {

	hm.mutex.RLock()
	defer hm.mutex.RUnlock()
	v, ok := hm.hallTable[hallID]
	return v, ok
}

func (hm *hallManager) GetHalls() []entity.Hall {

	hm.mutex.RLock()
	defer hm.mutex.RUnlock()

	halls := []entity.Hall{}
	for _, v := range hm.hallTable {
		halls = append(halls, v)
	}

	return halls
}

func (hm *hallManager)ContainHall(hallID int)bool{
	_, ok := hm.hallTable[hallID]
	return ok
}