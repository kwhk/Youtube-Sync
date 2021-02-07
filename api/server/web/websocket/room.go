package websocket

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/kwhk/sync/api/server/utils/timer"
)

type Room struct {
	id uuid.UUID
	host *Client
	register chan *Client
	unregister chan *Client
	clients map[uuid.UUID]*Client
	messageQueue chan Message

	// Video
	video *video

	// store ping for each client in ms
	clientPing map[uuid.UUID]int
	clientPingMeasure map[uuid.UUID][]int
	clientLastPing map[uuid.UUID]time.Time
}

type video struct {
	url string
	name string
	timer *timer.VideoTimer
	isPlaying bool
}

// FOR TESTING
func (v *video) elapsed() {
	fmt.Println(v.timer.Elapsed())
}

type JoinData struct {
	RoomID uuid.UUID `json:"roomID"`
	ClientID uuid.UUID `json:"clientID"`
	VideoURL string `json:"videoUrl"`
	VideoElapsed int64 `json:"videoElapsed"`
	VideoIsPlaying bool `json:"videoIsPlaying"`
}

func NewRoom() *Room {
	return &Room {
		id: uuid.New(),
		host: nil,
		register: make(chan *Client),
		unregister: make(chan *Client),
		clients: make(map[uuid.UUID]*Client),
		messageQueue: make(chan Message),
		// for testing purposes
		video: &video{
			url: "0-q1KafFCLU", 
			name: "IU Celebrity", 
			timer: &timer.VideoTimer{ Start: time.Now(), Progress: 0}, 
			isPlaying: false,
		},
		clientPing: make(map[uuid.UUID]int),
		clientPingMeasure: make(map[uuid.UUID][]int),
		clientLastPing: make(map[uuid.UUID]time.Time),
		
	}
}

func (room *Room) Start() {
	//// FOR TESTING SYNCING TIMES
	ticker := time.NewTicker(1 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <- ticker.C:
				room.video.elapsed()
			case <- quit:
				ticker.Stop()
				return
			}
		}
	}()
	////
	for {
		select {
		case client := <-room.register:
			// Set first client to register to room as host.
			if len(room.clients) == 0 {
				room.host = client
			}
			room.clients[client.id] = client
			room.clientPingMeasure[client.id] = make([]int, 10)
			fmt.Println("Size of connection room: ", len(room.clients))
			client.conn.WriteJSON(Message{ Action: "event", Event: Event{ Name: "join", Data: JoinData{client.room.id, client.id, room.video.url, room.video.timer.Elapsed(), room.video.isPlaying}}})
			for _, other := range room.clients {
				if other != client {
					other.conn.WriteJSON(Message{ Action: "message", Event: Event{Name: "Connect", Data: "New User Joined, ID: " + client.id.String()} })
				}
			}

		case client := <-room.unregister:
			if _, ok := room.clients[client.id]; ok {
				delete(room.clients, client.id)
				close(client.send)
			}
			fmt.Println("Size of connection room: ", len(room.clients))
			for _, client := range room.clients {
				client.conn.WriteJSON(Message{ Action: "message", Event: Event{Name: "Disconnect", Data: "User Disconnected"} })
			}
		case message := <-room.messageQueue:
			room.messageController(message);
		}
	}
}

// messageController following front controller pattern.
func (room *Room) messageController(message Message) {
	eventName := message.Event.Name
	newMessage := eventHandler(eventName, message, room)

	if newMessage != nil {
		room.dispatcher(*newMessage)
	}
}

func (room *Room) dispatcher(message Message) {
	// Dispatch message.
	if message.Target != nil {
		room.broadcast(message)
	} else {
		client := room.clients[*message.Source]
		room.emit(message, client)
	}
}

// Broadcast sends a message to all users in the same room
// with the option of including / excluding the sender.
func (room *Room) broadcast(message Message) {
	if message.Target.IncludeSender == false {
		for id, client := range room.clients {
			if id != *message.Source {
				select {
				case client.send <- message:
				default:
					// Done, no more messages to send.
					close(client.send)
					delete(room.clients, id)
				}
			}
		}
	} else {
		for id, client := range room.clients {
			select {
			case client.send <- message:
			default:
				// Done, no more messages to send.
				close(client.send)
				delete(room.clients, id)
			}
		}
	}
}

// Emit onlys ends 
func (room *Room) emit(message Message, client *Client) {
	client.send <- message
}