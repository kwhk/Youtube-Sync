package redis

import (
	"log"
	models "github.com/kwhk/sync/api/models/redis"
	"github.com/go-redis/redis/v8"
	"fmt"
	"strings"
)

const (
	/* 	Adding new videos in video queue will require this separator instead
		of the colon separator defined in newKey(). This is to easily split
		keys to separate the score and the video value in the queue for
		reordering. I'm specifically using this character because it will
		not conflict with other characters in the Video object (i.e. URL)
	*/
	videoQueueSep = ">"
)

type Video struct {
	Encoding []byte
}

type VideoClock struct {
	Encoding []byte
}

func (clock *VideoClock) Encode() []byte{
	return clock.Encoding
}

func (video *Video) Encode() []byte {
	return video.Encoding
}

type PlayerRepository struct {
	Redis *redis.Client
}

// resetVideoQueue resets scores to integer values when there is a collision of keys 
// due to scores being the same
func (repo *PlayerRepository) resetVideoQueue(roomID string) bool {
	resArr, err := repo.Redis.ZRangeByScoreWithScores(ctx, newKey(roomKey, roomID, queueKey), &redis.ZRangeBy{Min: "-inf", Max: "+inf"}).Result()
	if err != nil {
		log.Printf("Failed to remove video from queue in room %s.\n", roomID)
		return false
	}

	pipe := repo.Redis.Pipeline()
	for i, val := range resArr {
		video := strings.Split(val.Member.(string), videoQueueSep)[1]
		resArr[i].Member = fmt.Sprintf("%f", float64(i)) + videoQueueSep + video
		resArr[i].Score = float64(i)
		pipe.ZAdd(ctx, newKey(roomKey, roomID, queueKey), &resArr[i])
	}

	_, err2 := pipe.Exec(ctx)
	if err2 != nil {
		log.Printf("Failed to remove video from queue in room %s.\n", roomID)
		return false
	}

	return true
}

func (repo *PlayerRepository) AddToVideoQueue(roomID string, video models.Encodable) bool {
	resArr, err := repo.Redis.ZRevRangeByScoreWithScores(ctx, newKey(roomKey, roomID, queueKey), &redis.ZRangeBy{Min: "-inf", Max: "+inf", Offset: 0, Count: 1}).Result()
	if err != nil {
		log.Printf("Failed to add video to queue in room %s.\n", roomID)
		return false
	}
	
	// If queue is empty, then set score as 1
	var score float64
	if len(resArr) == 0 {
		score = 1
	// Set score of new video +1 of score of last element in queue
	} else {
		score = resArr[len(resArr)-1].Score + 1
	}
	
	_, err2 := repo.Redis.ZAdd(ctx, newKey(roomKey, roomID, queueKey), &redis.Z{Score: score, Member: fmt.Sprintf("%f", score) + videoQueueSep + string(video.Encode())}).Result()
	if err2 != nil {
		log.Printf("Failed to add video to queue in room %s.\n", roomID)
		return false
	}

	return true
}

func (repo *PlayerRepository) RemoveFromVideoQueue(roomID string, video models.Encodable, index int) bool {
	resArr, err := repo.Redis.ZRangeByScore(ctx, newKey(roomKey, roomID, queueKey), &redis.ZRangeBy{Min: "-inf", Max: "+inf"}).Result()
	if err != nil {
		log.Printf("Failed to remove video from queue in room %s.\n", roomID)
		return false
	}

	_, err2 := repo.Redis.ZRem(ctx, newKey(roomKey, roomID, queueKey), resArr[index]).Result()
	if err2 != nil {
		log.Printf("Failed to remove video from queue in room %s.\n", roomID)
		return false
	}

	return true
}

func (repo *PlayerRepository) ReorderVideoQueue(roomID string, newIndex int, oldIndex int) bool {
	resArr, err := repo.Redis.ZRangeByScoreWithScores(ctx, newKey(roomKey, roomID, queueKey), &redis.ZRangeBy{Min: "-inf", Max: "+inf"}).Result()
	if err != nil {
		log.Printf("Failed to reorder video queue in room %s.\n", roomID)
		return false
	}

	// Split video string into {score}>{video} to get video value.
	video := strings.Split(resArr[oldIndex].Member.(string), videoQueueSep)[1]

	// Remove video from old index in queue.
	_, err2 := repo.Redis.ZRem(ctx, newKey(roomKey, roomID, queueKey), resArr[oldIndex].Member).Result()
	if err2 != nil {
		log.Printf("Failed to reorder video queue in room %s.\n", roomID)
		return false
	}

	var score float64
	if newIndex == 0 {
		score = resArr[0].Score / 2
	} else if newIndex == len(resArr)-1 {
		score = resArr[newIndex].Score + 1
	}
	_, err3 := repo.Redis.ZAdd(ctx, newKey(roomKey, roomID, queueKey), &redis.Z{Score: score, Member: fmt.Sprintf("%f", score) + videoQueueSep + video,}).Result()
	if err3 != nil {
		log.Printf("Failed to reorder video queue in room %s.\n", roomID)
		return false
	}

	return true
}

func (repo *PlayerRepository) EmptyVideoQueue(roomID string) bool {
	_, err := repo.Redis.Del(ctx, newKey(roomKey, roomID, queueKey)).Result()
	if err != nil {
		log.Printf("Failed to empty video queue in room %s.\n", roomID)
		return false
	}
	return true
}

func (repo *PlayerRepository) GetCurrVideo(roomID string) (models.Encodable, bool) {
	video, err := repo.Redis.HGet(ctx, newKey(roomKey, roomID), currVideoKey).Bytes()
	if err != nil || err == redis.Nil {
		log.Printf("Failed to get current video in room %s.\n", roomID)
		return &Video{}, false
	}

	return &Video{
		Encoding: video,
	}, true
}

func (repo *PlayerRepository) SetCurrVideo(roomID string, video models.Encodable) bool {
	_, err := repo.Redis.HSet(ctx, newKey(roomKey, roomID), currVideoKey, video.Encode()).Result()
	if err != nil {
		log.Printf("Failed to set current video in room %s: %s\n", roomID, err)
		return false
	}
	return true
}

func (repo *PlayerRepository) GetClock(roomID string) (models.Encodable, bool) {
	val, err := repo.Redis.HGet(ctx, newKey(roomKey, roomID), clockKey).Bytes()
	if err != nil || err == redis.Nil {
		log.Printf("Failed to get clock from room %s.\n", roomID)
		return &VideoClock{}, false
	}
	return &VideoClock{
		Encoding: val,
	}, true
}

func (repo *PlayerRepository) SetClock(roomID string, clock models.Encodable) bool {
	_, err := repo.Redis.HSet(ctx, newKey(roomKey, roomID), clockKey, clock.Encode()).Result()
	if err != nil {
		log.Printf("Failed to set clock in room %s: %s\n", roomID, err)
		return false
	}
	return true
}