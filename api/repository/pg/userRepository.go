package pg 

import (
	"log"
	models "github.com/kwhk/sync/api/models/pg"
	"github.com/go-pg/pg/v10"
)

type User struct {
	ID string
	Name string
	RoomID string `pg:"room_id"`
	tableName struct{} `pg:"users"`
}

func (user *User) GetID() string {
	return user.ID
}

func (user *User) GetName() string {
	return user.Name
}

type UserRepository struct {
	DB *pg.DB
}

func (repo *UserRepository) AddUser(user models.User) {
	temp := &User{
		ID: user.GetID(),
		Name: user.GetName(),
	}
	_, err := repo.DB.Model(temp).Insert()
	if err != nil {
		log.Printf("Failed to add user %s to db.\n", user.GetID())
	}
}

func (repo *UserRepository) DeleteUser(user models.User) {
	temp := &User{
		ID: user.GetID(),
	}
	_, err := repo.DB.Model(temp).WherePK().Delete()
	if err != nil {
		log.Printf("Failed to delete user %s from db.\n", user.GetID())
	}
}

func (repo *UserRepository) FindUserByID(ID string) models.User {
	temp := &User{
		ID: ID,
	}

	err := repo.DB.Model(temp).WherePK().Select()
	if err != nil {
		log.Printf("Failed to find user id %s from db.\n", ID)
	}

	return temp
}

func (repo *UserRepository) GetAllUsers() []models.User {
	var users []User
	err := repo.DB.Model(&users).Select()
	if err != nil {
		log.Println(err)
		log.Printf("Failed to get all users from db.\n")
	}

	data := make([]models.User, 0)
	for _, elem := range(users) {
		data = append(data, &elem)
	}

	return data
}

func (repo *UserRepository) JoinRoom(user models.User, room models.Room) {
	temp := &User{
		ID: user.GetID(),
		RoomID: room.GetID(),
	}

	_, err := repo.DB.Model(temp).Set("room_id = ?room_id").Where("id = ?id").Update()

	if err != nil {
		log.Println(err)
		log.Printf("Failed to join user %s to room %s\n", user.GetID(), room.GetID())
	}
}	