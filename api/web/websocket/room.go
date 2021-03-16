package websocket

import (
	"fmt"
	"time"
	"sync"
	"context"
	"log"
	"encoding/json"

	"github.com/google/uuid"

	"github.com/kwhk/sync/api/utils/timer"
	"github.com/kwhk/sync/api/config"
)

var ctx = context.Background()

type Room struct {
	ID uuid.UUID `json:"id"`
	Name string `json:"name"`
	Private bool `json:"private"`

	wsServer *WsServer
	register chan *Client `json:"-"`
	unregister chan *Client `json:"-"`
	Clients map[string]*Client `json:"-"`
	// All messages that need to broadcasted to other servers.
	broadcast chan Message `json:"-"`
	// All messages that need to be sent from local server to client

	// Video
	Video videoDetails `json:"-"`
}

type videoDetails struct {
	Curr Video
	Queue []Video
	Mu sync.Mutex
}

type Video struct {
	// URL of video.
	URL string `json:"url"`
	// Duration of video in ms.
	Duration int64 `json:"-"`
	// Timer to record how much time elapsed since video start.
	Timer *timer.VideoTimer `json:"-"`
	// Status to notify joining users playback state.
	IsPlaying bool `json:"isPlaying"`
}

// FOR TESTING
func (v *Video) elapsed() {
	fmt.Println(v.Timer.Elapsed())
}


func newRoom(name string, private bool, server *WsServer) *Room {
	return &Room {
		ID: uuid.New(),
		Name: name,
		Private: private,
		wsServer: server,
		register: make(chan *Client),
		unregister: make(chan *Client),
		Clients: make(map[string]*Client),
		broadcast: make(chan Message),
		Video: videoDetails{
			Curr: Video{
				Timer: &timer.VideoTimer{ Start: time.Now(), Progress: 0}, 
			},
			Queue: make([]Video, 0),
		},
	}
}

func (room *Room) run() {
	go room.subscribeToRoomMessages()

	for {
		select {
		case client := <-room.register:
			room.wsServer.workerPool.AddJob(func() {room.registerClientInRoom(client)})
		case client := <-room.unregister:
			room.wsServer.workerPool.AddJob(func() {room.unregisterClientInRoom(client)})
		case message := <-room.broadcast:
			room.wsServer.workerPool.AddJob(func() {room.eventHandler(message)})
		}
	}
}

func (room *Room) unregisterClientInRoom(client *Client) {
	if _, ok := room.Clients[client.GetID()]; ok {
		delete(room.Clients, client.GetID())
	}
	fmt.Println("Size of connection room: ", len(room.Clients))

	room.notifyClientLeft(client)
}

func (room *Room) registerClientInRoom(client *Client) {
	room.Clients[client.GetID()] = client
	fmt.Println("Size of connection room: ", len(room.Clients))
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
	for id, client := range room.Clients {
		select {
		case client.send <- message.encode():
		default:
			// Done, no more messages to send.
			delete(room.Clients, id)
		}
	}
}

func (room *Room) publishRoomMessage(message []byte) {
	err := config.Redis.Publish(ctx, room.GetID(), message).Err()

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

		room.broadcastToClients(msg)
	}
}

func (room *Room) eventHandler(message Message) {
	var h handler
	action := message.Action
	switch action {
	// host only events
	case PlayVideoAction, PauseVideoAction, SeekToVideoAction:
		h = newPlayback(message, room, action)
	case PlayVideoQueueAction, AddVideoQueueAction, RemoveVideoQueueAction, EmptyVideoQueueAction:
		h = videoQueue{message, room, action}
	case JoinRoomAction, LeaveRoomAction:
		break
	default:
		log.Printf("Room eventHandler does not recognize event %s.\n", action)
		return
	}

	if msg, ok := h.handle(); ok {
		room.publishRoomMessage(msg.encode())
	}
}

func (room *Room) GetID() string {
	return room.ID.String()
}

func (room *Room) GetName() string {
	return room.Name
}

func (room *Room) GetPrivate() bool {
	return room.Private
}

