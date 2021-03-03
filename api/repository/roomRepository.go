package repository

import (
	"log"
	"github.com/kwhk/sync/api/models"
	"github.com/go-pg/pg/v10"
)

type Room struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Private bool `json:"private"`
	tableName struct{} `pg:"rooms"`
}

func (room *Room) GetID() string {
	return room.ID
}

func (room *Room) GetName() string {
	return room.Name
}

func (room *Room) GetPrivate() bool {
	return room.Private
}

type RoomRepository struct {
	DB *pg.DB
}

func (repo *RoomRepository) AddRoom(room models.Room) {
	temp := &Room{
		ID: room.GetID(),
		Name: room.GetName(),
		Private: room.GetPrivate(),
	}
	_, err := repo.DB.Model(temp).Insert()
	if err != nil {
		log.Printf("Failed to add room %s to db.\n", room.GetID())
	}
}

func (repo *RoomRepository) DeleteRoom(room models.Room) {
	temp := &Room{
		ID: room.GetID(),
	}
	_, err := repo.DB.Model(temp).WherePK().Delete()
	if err != nil {
		log.Printf("Failed to delete room %s from db.\n", room.GetID())
	}
}

func (repo *RoomRepository) FindRoomByID(ID string) models.Room {
	temp := &Room{
		ID: ID,
	}

	err := repo.DB.Model(temp).WherePK().Select()
	if err != nil {
		log.Printf("Failed to find room id %s from db.\n", ID)
	}

	return temp
}

func (repo *RoomRepository) FindRoomByName(name string) models.Room {
	temp := &Room{
		Name: name,
	}

	err := repo.DB.Model(temp).Where("name = ?", temp.Name).Limit(1).Select()
	if err != nil {
		log.Printf("Failed to find room name %s from db.\n", name)
	}

	return temp
}

func (repo *RoomRepository) GetAllRooms() []models.Room {
	var rooms []Room
	err := repo.DB.Model(&rooms).Select()
	if err != nil {
		log.Printf("Failed to get all rooms from db.\n")
	}

	data := make([]models.Room, len(rooms))
	for i, elem := range(rooms) {
		data[i] = &Room{ID: elem.ID, Name: elem.Name, Private: elem.Private}
	}

	return data
}