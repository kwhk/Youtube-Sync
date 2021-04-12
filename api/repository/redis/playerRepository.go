package redis

import (
	"log"
	models "github.com/kwhk/sync/api/models/redis"
	"github.com/go-redis/redis/v8"
)

// type Player struct {
// 	RoomID string
// 	Queue []models.Video
// 	CurrVideo models.Video
// 	Clock models.VideoClock
// }

type Video struct {
	Encoding []byte
}

type VideoClock struct {
	Encoding []byte
}

func (clock *VideoClock) GetEncoding() []byte{
	return clock.Encoding
}

func (video *Video) GetEncoding() []byte {
	return video.Encoding
}

// func (player *Player) GetRoomID() string {
// 	return player.RoomID
// }

// func (player *Player) GetQueue() []models.Video { 
// 	return player.Queue
// }

// func (player *Player) GetCurrVideo() models.Video {
// 	return player.CurrVideo
// }

// func (player *Player) GetTimer() models.VideoClock {
// 	return player.Clock
// }

type PlayerRepository struct {
	Redis *redis.Client
}

func (repo *PlayerRepository) AddToVideoQueue(roomID string, video models.Video) bool {
	err := repo.Redis.RPush(ctx, newKey(roomKey, roomID, queueKey), video.GetEncoding())
	if err != nil {
		log.Printf("Failed to add video to queue in room %s.\n", roomID)
		return false
	}

	return true
}

func (repo *PlayerRepository) RemoveFromVideoQueue(roomID string, video models.Video) bool {
	err := repo.Redis.LRem(ctx, newKey(roomKey, roomID, queueKey), 1, video.GetEncoding())
	if err != nil {
		log.Printf("Failed to remove video from queue in room %s.\n", roomID)
		return false
	}
	return true
}

func (repo *PlayerRepository) EmptyVideoQueue(roomID string, video models.Video) bool {
	err := repo.Redis.Del(ctx, newKey(roomKey, roomID, queueKey))
	if err != nil {
		log.Printf("Failed to empty video queue in room %s.\n", roomID)
		return false
	}
	return true
}

func (repo *PlayerRepository) GetCurrVideo(roomID string) (models.Video, bool) {
	video, err := repo.Redis.HGet(ctx, newKey(roomKey, roomID), currVideoKey).Bytes()
	if err != nil || err == redis.Nil {
		log.Printf("Failed to get current video in room %s.\n", roomID)
		return &Video{}, false
	}

	return &Video{
		Encoding: video,
	}, true
}

func (repo *PlayerRepository) SetCurrVideo(roomID string, video models.Video) bool {
	err := repo.Redis.HSet(ctx, newKey(roomKey, roomID), currVideoKey, video.GetEncoding())
	if err != nil {
		log.Printf("Failed to set current video in room %s.\n", roomID)
		return false
	}
	return true
}

func (repo *PlayerRepository) GetClock(roomID string) (models.Clock, bool) {
	val, err := repo.Redis.HGet(ctx, newKey(roomKey, roomID), clockKey).Bytes()
	if err != nil || err == redis.Nil {
		log.Printf("Failed to get clock from room %s.\n", roomID)
		return &VideoClock{}, false
	}
	return &VideoClock{
		Encoding: val,
	}, true
}

func (repo *PlayerRepository) SetClock(roomID string, clock models.Clock) bool {
	err := repo.Redis.HSet(ctx, newKey(roomKey, roomID), clockKey, clock.GetEncoding())
	if err != nil {
		log.Printf("Failed to set clock in room %s.\n", roomID)
		return false
	}
	return true
}