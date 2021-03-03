package websocket 

import (
	"fmt"
	"time"
	"log"

	"github.com/kwhk/sync/api/utils/timer"
)

type videoQueue struct {
	message Message
	room *Room
	action string
}


func (q videoQueue) createVideo(data map[string]interface{}) (Video, bool) {
	url, ok1 := data["url"]
	duration, ok2 := data["duration"]	

	if !ok1 || !ok2 {
		return Video{}, false
	}
	
	video := Video{URL: url.(string), Duration: int64(duration.(float64)), Timer: &timer.VideoTimer{ Start: time.Now(), Progress: 0}, IsPlaying: false}
	return video, true
}

func (q *videoQueue) getVideo(url string, index int) (int, bool) {
	for _, video := range q.room.Video.Queue {
		if video.URL == url /* && index == i */ {
			return index, true
		}
	}

	return -1, false
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
	var video map[string]interface{} = q.message.Data.(map[string]interface{})
	var url string = video["url"].(string)
	var index int = int(video["index"].(float64))

	if index, ok := q.getVideo(url, index); ok {
		fmt.Printf("Currently playing video is %s\n", url)
		q.room.Video.Curr = q.room.Video.Queue[index]
		return q.message, true
	}

	fmt.Printf("Video %s does not exist in queue.\n", url)
	return Message{}, false
}

func (q videoQueue) add() (Message, bool) {
	q.room.Video.Mu.Lock()
	defer q.room.Video.Mu.Unlock()
	video, ok := q.createVideo(q.message.Data.(map[string]interface{}))

	if !ok {
		fmt.Printf("Failed to add \"%s\" because not all video details were available.\n", video.URL)
		return q.message, false
	}

	q.room.Video.Queue = append(q.room.Video.Queue, video)
	fmt.Printf("Added \"%s\" to video queue\n", video.URL)
	return q.message, true
}

func (q videoQueue) remove() (Message, bool) {
	q.room.Video.Mu.Lock()
	defer q.room.Video.Mu.Unlock()
	var video map[string]interface{} = q.message.Data.(map[string]interface{})
	var url string = video["url"].(string)
	var index int = int(video["index"].(float64))

	if index, ok := q.getVideo(url, index); ok {
		q.room.Video.Queue = append(q.room.Video.Queue[:index], q.room.Video.Queue[index+1:]...)
		fmt.Printf("Removed \"%s\" from video queue\n", url)
		return q.message, true
	}

	fmt.Printf("Failed to remove \"%s\" from video queue since it was not found.\n", url)
	return Message{}, false
}

func (q videoQueue) empty() (Message, bool) {
	q.room.Video.Mu.Lock()
	defer q.room.Video.Mu.Unlock()
	q.room.Video.Queue = nil
	fmt.Printf("Emptied video queue\n")
	return q.message, true
}