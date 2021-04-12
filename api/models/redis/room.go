package redis

type Room interface {
	GetID() string
	GetCurrVideo() Video
	GetQueue() []Video
	GetClock() Clock
}

type RoomRepository interface {
	AddRoom(room Room) bool
	DeleteRoom(room Room) bool
	FindRoomByID(id string) (Room, bool)
	GetUsers(room Room) ([]User, bool)
}