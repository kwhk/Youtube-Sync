package websocket 

import (
	"fmt"
	"log"
	"github.com/kwhk/sync/api/repository/redis"
)

type videoQueue struct {
	message Message
	room *Room
	playerRepo redis.PlayerRepository
	action string
}

func newVideoQueue(message Message, room *Room, action string) videoQueue {
	new := videoQueue{
		message: message,
		room: room,
		playerRepo: room.wsServer.playerRepository,
		action: action,
	}
	return new
}

func (q videoQueue) createVideo(data map[string]interface{}) (Video, bool) {
	url, ok1 := data["url"]
	duration, ok2 := data["duration"]	

	if !ok1 || !ok2 {
		return Video{}, false
	}
	
	video := Video{URL: url.(string), Duration: int64(duration.(float64))}
	return video, true
}

func (q *videoQueue) getVideo(url string, index int) (Video, bool) {
	for i, video := range q.room.Video.Queue {
		if video.URL == url && index == i {
			return video, true
		}
	}

	return Video{}, false
}

func (q videoQueue) handle() (Message, bool) {
	switch q.action {
	case AddVideoQueueAction:
		return q.add()
	case RemoveVideoQueueAction:
		return q.remove()
	case EmptyVideoQueueAction:
		return q.empty()
	case PlayVideoQueueAction:
		return q.play()
	default:
		log.Println("No such video queue action exists")
		return Message{}, false
	}
}

func (q videoQueue) play() (Message, bool) {
	var msg map[string]interface{} = q.message.Data.(map[string]interface{})
	var url string = msg["url"].(string)
	var index int = int(msg["index"].(float64))

	if video, ok := q.getVideo(url, index); ok {
		q.playerRepo.SetCurrVideo(q.room.ID, video)
		q.room.Video.Curr.Details = q.room.Video.Queue[index]
		q.room.Video.Curr.Index = index
		q.room.Clock.Reset()
		q.playerRepo.SetClock(q.room.ID, q.room.Clock)
		return q.message, true
	}

	fmt.Printf("Video %s does not exist in queue.\n", url)
	return Message{}, false
}

func (q videoQueue) add() (Message, bool) {
	q.room.Video.mu.Lock()
	defer q.room.Video.mu.Unlock()
	video, ok := q.createVideo(q.message.Data.(map[string]interface{}))

	if !ok {
		fmt.Printf("Failed to add \"%s\" because not all video details were available.\n", video.URL)
		return q.message, false
	}

	q.room.Video.Queue = append(q.room.Video.Queue, video)
	q.playerRepo.AddToVideoQueue(q.room.ID, video)
	fmt.Printf("Added \"%s\" to video queue\n", video.URL)
	return q.message, true
}

func (q videoQueue) remove() (Message, bool) {
	q.room.Video.mu.Lock()
	defer q.room.Video.mu.Unlock()
	var msg map[string]interface{} = q.message.Data.(map[string]interface{})
	var url string = msg["url"].(string)
	var index int = int(msg["index"].(float64))

	if _, ok := q.getVideo(url, index); ok {
		q.playerRepo.RemoveFromVideoQueue(q.room.ID, q.room.Video.Queue[index], index)
		q.room.Video.Queue = append(q.room.Video.Queue[:index], q.room.Video.Queue[index+1:]...)
		fmt.Printf("Removed \"%s\" from video queue\n", url)
		return q.message, true
	}

	fmt.Printf("Failed to remove \"%s\" from video queue since it was not found.\n", url)
	return Message{}, false
}

func (q videoQueue) empty() (Message, bool) {
	q.room.Video.mu.Lock()
	defer q.room.Video.mu.Unlock()
	q.room.Video.Queue = nil
	q.playerRepo.EmptyVideoQueue(q.room.ID)
	fmt.Printf("Emptied video queue\n")
	return q.message, true
}