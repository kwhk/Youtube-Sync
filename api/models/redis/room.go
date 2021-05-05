package redis

type Room interface {
	GetID() string
	GetCurrVideo() Encodable
	GetQueue() []Encodable
	GetClock() Encodable 
}

type RoomRepository interface {
	AddRoom(room Room) bool
	DeleteRoom(room Room) bool
	FindRoomByID(id string) (Room, bool)
	GetUsers(room Room) ([]User, bool)
}