package websocket

import (
	"fmt"
	"time"

	"github.com/kwhk/sync/api/server/utils/timer"
)

type videoQueue struct {
	message Message
	room *Room
	action string
}

func (q videoQueue) createVideo(data map[string]interface{}) (video, bool) {
	url, ok1 := data["url"]
	duration, ok2 := data["duration"]	

	if !ok1 || !ok2 {
		return video{}, false
	}
	
	video := video{url: url.(string), duration: int64(duration.(float64)), timer: &timer.VideoTimer{ Start: time.Now(), Progress: 0}, isPlaying: false}
	return video, true
}

func (q videoQueue) getVideo(url string, index int) (int, bool) {
	for i, video := range q.room.video.queue {
		if video.url == url && index == i {
			return index, true
		}
	}

	return -1, false
}

func (q videoQueue) execute() (Message, bool) {
	switch q.action {
	case "addVideoQueue":
		return q.add()
	case "removeVideoQueue":
		return q.remove()
	case "emptyVideoQueue":
		return q.empty()
	case "playVideoQueue":
		return q.play()
	default:
		fmt.Println("No such video queue operation exists")
		return Message{}, false
	}
}

func (q videoQueue) play() (Message, bool) {
	var video map[string]interface{} = q.message.Event.Data.(map[string]interface{})
	var url string = video["url"].(string)
	var index int = int(video["index"].(float64))

	if index, ok := q.getVideo(url, index); ok {
		fmt.Printf("Currently playing video is %s\n", url)
		q.room.video.curr = q.room.video.queue[index]
		return q.message, true
	}

	fmt.Printf("Video %s does not exist in queue.\n", url)
	return Message{}, false
}

func (q videoQueue) add() (Message, bool) {
	q.room.video.mu.Lock()
	defer q.room.video.mu.Unlock()
	video, ok := q.createVideo(q.message.Event.Data.(map[string]interface{}))

	if !ok {
		fmt.Printf("Failed to add \"%s\" because not all video details were available.\n", video.url)
		return q.message, false
	}

	q.room.video.queue = append(q.room.video.queue, video)
	fmt.Printf("Added \"%s\" to video queue\n", video.url)
	return q.message, true
}

func (q videoQueue) remove() (Message, bool) {
	q.room.video.mu.Lock()
	defer q.room.video.mu.Unlock()
	var video map[string]interface{} = q.message.Event.Data.(map[string]interface{})
	var url string = video["url"].(string)
	var index int = int(video["index"].(float64))

	if index, ok := q.getVideo(url, index); ok {
		q.room.video.queue = append(q.room.video.queue[:index], q.room.video.queue[index+1:]...)
		fmt.Printf("Removed \"%s\" from video queue\n", url)
		return q.message, true
	}

	fmt.Printf("Failed to remove \"%s\" from video queue since it was not found.\n", url)
	return Message{}, false
}

func (q videoQueue) empty() (Message, bool) {
	q.room.video.mu.Lock()
	defer q.room.video.mu.Unlock()
	q.room.video.queue = nil
	fmt.Printf("Emptied video queue\n")
	return q.message, true
}