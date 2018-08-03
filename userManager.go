package nex

import (
	"sync"
	"errors"
	"github.com/cruisechang/nex/entity"
)

var (
	userCount = 0
)

type UserManager interface {
	CreateUser(connID string) (entity.User, error)
	RemoveUser(userID int)
	GetUser(userID int) (entity.User, bool)
	GetUserByConnID(connID string) (entity.User, bool)
	GetUsers() ([]entity.User)
}

type userManager struct {
	userTable map[string]entity.User
	mutex     sync.Mutex
}

func NewUserManager() (UserManager, error) {
	return &userManager{
		userTable: make(map[string]entity.User),
		mutex:     sync.Mutex{},
	}, nil
}

//AddUser returns user just created or error
func (um *userManager) CreateUser(connID string) (entity.User, error) {

	if _, ok := um.userTable[connID]; ok {
		return nil, errors.New("user already in table")
	}

	um.mutex.Lock()
	userCount += 1
	u := entity.NewUser(userCount, connID)
	um.userTable[connID] = u
	um.mutex.Unlock()

	return u, nil
}

func (um *userManager) RemoveUser(userID int) {
	um.mutex.Lock()
	for _, u := range um.userTable {
		if u.UserID() == userID {
			delete(um.userTable, u.ConnID())
			break
		}
	}
	defer um.mutex.Unlock()

}

func (um *userManager) GetUser(userID int) (entity.User, bool) {

	for _, u := range um.userTable {
		if u.UserID() == userID {
			return u, true
		}
	}
	return nil, false
}

func (um *userManager) GetUserByConnID(connID string) (entity.User, bool) {

	v, ok := um.userTable[connID]
	return v, ok
}

func (um *userManager) GetUsers() []entity.User {

	users := []entity.User{}
	for _, v := range um.userTable {
		users = append(users, v)
	}

	return users
}
