package nex

import (
	"testing"
)

func TestUserManager(t *testing.T) {


	connID:="connID"

	um,_:= NewUserManager()

	user,err:=um.CreateUser(connID)
	if err!=nil{
		t.Error("TestUserManager createUser error:%s\n",err.Error())

	} else if user==nil{
		t.Error("TestUserManager createUser user==nil\n")
	}

	t.Logf("TestUserManager createUser user:%#v\n",user)

	getUser,ok:=um.GetUser(user.UserID())

	if !ok{
		t.Error("TestUserManager createUser error:%s\n",err.Error())

	}

	t.Logf("TestUserManager GetUser user:%#v\n",getUser)

	um.RemoveUser(user.UserID())

	um.GetUsers()
}