package models

type User interface {
	GetID() string
	GetName() string
}

type UserRepository interface {
	AddUser(user User)
	DeleteUser(user User)
	FindUserByID(ID string) User
	GetAllUsers() []User
	JoinRoom(user User, room Room)
}