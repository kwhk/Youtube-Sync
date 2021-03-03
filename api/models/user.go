package models

type User interface {
	GetID() string
	GetName() string
}

type UserRepository interface {
	AddUser(user User)
	RemoveUser(user User)
	FindUserById(ID string) User
	GetAllUsers() []User
}