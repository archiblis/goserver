package main

type User struct {
	gs *GameSessioner
	id uint64
}

func NewUser() *User {
	user := &User{}
	return user
}
