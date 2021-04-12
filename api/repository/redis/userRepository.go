package redis

import (
	"log"
	models "github.com/kwhk/sync/api/models/redis"
	"github.com/go-redis/redis/v8"
)

type User struct {
	ID string
}

func (user *User) GetID() string {
	return user.ID
}

type UserRepository struct {
	Redis *redis.Client
}

func (repo *UserRepository) JoinRoom(user models.User, room models.Room) {
	
	err := repo.Redis.SAdd(ctx, newKey(roomKey, room.GetID(), userKey + "s"), user.GetID()).Err()
	if err != nil {
		log.Printf("Failed to add user %s to room %s on redis.\n", user.GetID(), room.GetID())
	}
}

func (repo *UserRepository) LeaveRoom(user models.User, room models.Room) {
	err := repo.Redis.SRem(ctx, newKey(roomKey, room.GetID(), userKey + "s"), user.GetID()).Err()
	if err != nil {
		log.Printf("Failed to remove user %s from room %s.\n", user.GetID(), room.GetID())
	}
}

func (repo *UserRepository) AddUser(user models.User) {
	err := repo.Redis.HSetNX(ctx, userKey, user.GetID(), 1).Err()
	if err != nil {
		log.Printf("Failed to add user %s to redis.\n", user.GetID())
	}
}

func (repo *UserRepository) DeleteUser(user models.User) {
	err := repo.Redis.HDel(ctx, userKey, user.GetID()).Err()
	if err != nil {
		log.Printf("Failed to delete user %s from redis.\n", user.GetID())
	}
}

func (repo *UserRepository) FindUserByID(ID string) models.User {
	_, err := repo.Redis.HGet(ctx, userKey, ID).Result()
	if err != nil || err == redis.Nil {
		log.Printf("Failed to find user %s from redis.\n", ID)
	}

	return &User{
		ID: ID,
	}
}