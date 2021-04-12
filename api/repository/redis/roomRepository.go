package redis

import (
	"log"
	models "github.com/kwhk/sync/api/models/redis"
	"github.com/go-redis/redis/v8"
	"context"
)

var ctx context.Context = context.Background()

type Room struct {
	ID string
	CurrVideo models.Video
	Queue []models.Video
	Clock models.Clock
}

func (room *Room) GetID() string {
	return room.ID
}

func (room *Room) GetCurrVideo() models.Video {
	return room.CurrVideo
}

func (room *Room) GetQueue() []models.Video {
	return room.Queue
}

func (room *Room) GetClock() models.Clock {
	return room.Clock
}

type RoomRepository struct {
	Redis *redis.Client
}

func (repo *RoomRepository) checkRoomExists(roomID string) bool {
	val, err := repo.Redis.HGetAll(ctx, newKey(roomKey, roomID)).Result()
	if err != nil || err == redis.Nil {
		log.Printf("Failed to check if room %s exists.\n", roomID)
		return false
	}

	return len(val) != 0
}

func (repo *RoomRepository) AddRoom(room models.Room) bool {
	if repo.checkRoomExists(room.GetID()) {
		return false
	}

	pipe := repo.Redis.Pipeline()
	pipe.HSet(ctx, newKey(roomKey, room.GetID()), currVideoKey, room.GetCurrVideo().GetEncoding())
	pipe.HSet(ctx, newKey(roomKey, room.GetID()), clockKey, room.GetClock().GetEncoding())
	for _, video := range room.GetQueue() {
		pipe.RPush(ctx, newKey(roomKey, room.GetID(), queueKey), video.GetEncoding())
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Printf("Failed to add room %s.\n", room.GetID())
		return false
	}

	return true
}

func (repo *RoomRepository) DeleteRoom(room models.Room) bool {
	err := repo.Redis.HDel(ctx, roomKey, room.GetID()).Err()
	if err != nil {
		log.Printf("Failed to delete room %s from redis.\n", room.GetID())
		return false
	}
	return true
}

func (repo *RoomRepository) FindRoomByID(ID string) (models.Room, bool) {
	// Room does not exist
	if !repo.checkRoomExists(ID) {
		return &Room{}, false
	}

	pipe := repo.Redis.Pipeline()
	queue := pipe.SMembers(ctx, newKey(roomKey, ID, queueKey))
	currVideo := pipe.HGet(ctx, newKey(roomKey, ID), currVideoKey)
	clock := pipe.HGet(ctx, newKey(roomKey, ID), clockKey)
	
	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Printf("Failed to get room %s details.\n", ID)
		return &Room{}, false
	}

	queueRes, _ := queue.Result()
	currVideoRes, _ := currVideo.Bytes()
	clockRes, _ := clock.Bytes()

	var newQueue []models.Video = make([]models.Video, len(queueRes))

	for index, video := range queueRes {
		newQueue[index] = &Video{Encoding: []byte(video)}
	}

	return &Room{
		ID: ID,
		CurrVideo: &Video{Encoding: currVideoRes},
		Queue: newQueue,
		Clock: &VideoClock{Encoding: clockRes},
	}, true
}

func (repo *RoomRepository) GetUsers(ID string) ([]models.User, bool) {
	val, err := repo.Redis.SMembers(ctx, newKey(roomKey, ID, userKey + "s")).Result()
	if err != nil {
		log.Printf("Failed to get all users from room %s.\n", ID)
		return nil, false
	}

	var users []models.User = make([]models.User, len(val))
	for i, userID := range val {
		users[i] = &User{ID: userID}
	}

	return users, true
}