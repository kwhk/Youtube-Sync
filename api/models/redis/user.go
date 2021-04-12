package redis

type User interface {
	GetID() string
}

type UserRepository interface {
	AddUser(user User)
	DeleteUser(user User)
	FindUserByID(ID string) User
	JoinRoom(user User, room Room)
}