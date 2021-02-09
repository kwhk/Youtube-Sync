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

func (q videoQueue) execute() (Message, bool) {
	switch q.action {
	case "addVideoQueue":
		return q.add()
	case "popVideoQueue":
		return q.pop()
	case "emptyVideoQueue":
		return q.empty()
	default:
		fmt.Println("No such video queue operation exists")
		return q.message, true
	}
}

func (q videoQueue) add() (Message, bool) {
	// get video url
	var videoInfo map[string]interface{} = q.message.Event.Data.(map[string]interface{})
	video := video{url: videoInfo["url"].(string), duration: int64(videoInfo["duration"].(float64)), timer: &timer.VideoTimer{ Start: time.Now(), Progress: 0}, isPlaying: false}
	queue := append(q.room.videoQueue, video)
	q.room.videoQueue = queue
	fmt.Printf("Added \"%s\" to video queue\n", videoInfo["url"])
	fmt.Printf("%v\n", q.room.videoQueue)
	return q.message, true
}

func (q videoQueue) pop() (Message, bool) {
	var videoInfo map[string]interface{} = q.message.Event.Data.(map[string]interface{})
	popped, queue := q.room.videoQueue[0], q.room.videoQueue[1:]

	if popped.url == videoInfo["url"] {
		fmt.Printf("Popped \"%s\" from video queue\n", videoInfo["url"])
		q.room.videoQueue = queue
		return q.message, true
	}

	fmt.Printf("Failed to pop \"%s\" from video queue, popped != videoURL\n", videoInfo["url"])
	return Message{}, false
}

func (q videoQueue) empty() (Message, bool) {
	q.room.videoQueue = nil
	fmt.Printf("Emptied video queue\n")
	return q.message, true
}