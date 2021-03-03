package repository

import (
	"log"
	"github.com/kwhk/sync/api/models"
	"github.com/go-pg/pg/v10"
)

type User struct {
	ID string
	Name string
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

func (repo *UserRepository) RemoveUser(user models.User) {
	temp := &User{
		ID: user.GetID(),
	}
	_, err := repo.DB.Model(temp).WherePK().Delete()
	if err != nil {
		log.Printf("Failed to delete user %s from db.\n", user.GetID())
	}
}

func (repo *UserRepository) FindUserById(ID string) models.User {
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