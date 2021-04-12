package websocket

import (
	"fmt"
	"time"
	"sync"
	"context"
	"log"
	"encoding/json"

	"github.com/google/uuid"
	models "github.com/kwhk/sync/api/models/redis"
	"github.com/kwhk/sync/api/utils/clock"
	"github.com/kwhk/sync/api/config"
)

var ctx = context.Background()

type Room struct {
	ID string `json:"id"`

	wsServer *WsServer `json:"-"`
	register chan *Client `json:"-"`
	unregister chan *Client `json:"-"`
	clients map[string]*Client `json:"-"`
	// All messages that need to broadcasted to other servers.
	broadcast chan Message `json:"-"`
	// All messages that need to be sent from local server to client

	// Video
	Video videoDetails `json:"video"`
	Clock *clock.Clock `json:"clock"`
}

type videoDetails struct {
	Curr Video `json:"curr"`
	Queue []Video `json:"queue"`
	mu sync.Mutex `json:"-"`
}

func newRoom(server *WsServer) *Room {
	return &Room {
		ID: uuid.New().String(),
		wsServer: server,
		register: make(chan *Client),
		unregister: make(chan *Client),
		clients: make(map[string]*Client),
		broadcast: make(chan Message),
		Clock: &clock.Clock{ Start: time.Now(), Progress: 0}, 
		Video: videoDetails{
			Curr: Video{},
			Queue: make([]Video, 0),
		},
	}
}

func newRoomFromRedis(room models.Room, server *WsServer) *Room {
	return &Room{
		ID: room.GetID(),
		wsServer: server,
		register: make(chan *Client),
		unregister: make(chan *Client),
		clients: make(map[string]*Client),
		broadcast: make(chan Message),
		Clock: clock.DecodeClock(room.GetClock().GetEncoding()),
		Video: videoDetails{
			Curr: decodeVideo(room.GetCurrVideo().GetEncoding()),
			Queue: readQueue(room.GetQueue()),
		},
	}
}

func (room *Room) run() {
	go room.subscribeToRoomMessages()

	for {
		select {
		case client := <-room.register:
			room.registerClient(client)
		case client := <-room.unregister:
			room.unregisterClient(client)
		case message := <-room.broadcast:
			room.wsServer.workerPool.AddJob(func() {room.publishRoomMessage(message.encode())})
		}
	}
}

func (room *Room) unregisterClient(client *Client) {
	if _, ok := room.clients[client.GetID()]; ok {
		delete(room.clients, client.GetID())
		fmt.Println("Size of connection room: ", len(room.clients))
		room.notifyClientLeft(client)
	} else {
		log.Printf("Failed to unregister client %s from room %s\n", client.GetID(), room.GetID())
		return
	}
}

func (room *Room) registerClient(client *Client) {
	room.clients[client.GetID()] = client
	fmt.Println("Size of connection room: ", len(room.clients))
	room.notifyClientJoined(client)
}

// Notify all clients in room that new client has joined
func (room *Room) notifyClientJoined(client *Client) {
	message := Message{
		Action: JoinRoomAction,
		Target: room,
		Data: client.GetID(),
	}

	room.publishRoomMessage(message.encode())
}

func (room *Room) notifyClientLeft(client *Client) {
	message := Message{
		Action: LeaveRoomAction,
		Target: room,
		Data: client.GetID(),
	}

	room.publishRoomMessage(message.encode())
}

func (room *Room) broadcastToClients(message Message) {
	for id, client := range room.clients {
		select {
		case client.send <- message.encode():
		default:
			// Done, no more messages to send.
			delete(room.clients, id)
		}
	}
}

func (room *Room) publishRoomMessage(message []byte) {
	err := config.Redis.Publish(ctx, room.ID, message).Err()

	if err != nil {
		log.Println(err)
	}
}

func (room *Room) subscribeToRoomMessages() {
	pubsub := config.Redis.Subscribe(ctx, room.GetID())

	ch := pubsub.Channel()

	for message := range ch {
		var msg Message
		if err := json.Unmarshal([]byte(message.Payload), &msg); err != nil {
			log.Printf("Error on unmarshal JSON message %s", err)
			return
		}

		room.wsServer.workerPool.AddJob(func() {
			if newMsg, ok := room.eventHandler(msg); ok {
				room.broadcastToClients(newMsg)
			}
		})
	}
}

func (room *Room) eventHandler(message Message) (Message, bool) {
	var h handler
	action := message.Action
	switch action {
	// host only events
	case PlayVideoAction, PauseVideoAction, SeekToVideoAction:
		h = newPlayback(message, room, action)
	case PlayVideoQueueAction, AddVideoQueueAction, RemoveVideoQueueAction, EmptyVideoQueueAction:
		h = videoQueue{message, room, action}
	case JoinRoomAction, LeaveRoomAction:
		return message, true
	default:
		log.Printf("Room eventHandler does not recognize event %s.\n", action)
		return Message{}, false
	}

	return h.handle()
}

func (room *Room) GetID() string {
	return room.ID
}

func (room *Room) GetCurrVideo() models.Video {
	return room.Video.Curr
}

func (room *Room) GetQueue() []models.Video {
	var queue []models.Video = make([]models.Video, len(room.Video.Queue))

	for index, val := range room.Video.Queue {
		queue[index] = val
	}

	return queue
}

// readQueue converts slice of models.Video to slice of video
// to match repository type
func readQueue(queue []models.Video) []Video {
	var newQueue []Video = make([]Video, len(queue))

	for index, val := range queue {
		newQueue[index] = decodeVideo(val.GetEncoding())
	}

	return newQueue
}

func (room *Room) GetClock() models.Clock {
	return room.Clock
}