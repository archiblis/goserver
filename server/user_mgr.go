package main

import (
//"aqua"
)

type UserMgr struct {
	users   map[uint64]*User
	gsUsers map[*GameSessioner]*User
}

func NewUserMgr() *UserMgr {
	um := &UserMgr{
		users:   make(map[uint64]*User),
		gsUsers: make(map[*GameSessioner]*User),
	}
	return um
}

func (this *UserMgr) AddUser(id uint64, gs *GameSessioner) {
	user := NewUser()
	user.id = id
	user.gs = gs

	this.users[id] = user
	this.gsUsers[gs] = user
}

func (this *UserMgr) GetUser(id uint64) *User {
	if user, ok := this.users[id]; ok == true {
		return user
	}

	return nil
}

/*
func (this *UserMgr) GetUserBySessioner(s *aqua.Sessioner) {
	gs := (*GameSessioner)(s)
	if user, ok := this.gsUsers[gs]; ok == true {
		return user
	}

	return nil
}
*/
