package pg

type Room interface {
	GetID() string
}

type RoomRepository interface {
	AddRoom(room Room)
	DeleteRoom(room Room)
	FindRoomByID(id string) Room
	GetAllRooms() []Room
}