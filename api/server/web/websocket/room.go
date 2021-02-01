package websocket 

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

type Room struct {
	id uuid.UUID
	host *Client
	register chan *Client
	unregister chan *Client
	clients map[uuid.UUID]*Client
	messageQueue chan Message
	// store ping for each client in ms
	clientPing map[uuid.UUID]int
	clientPingMeasure map[uuid.UUID][]int
	clientLastPing map[uuid.UUID]time.Time
}

func NewRoom() *Room {
	return &Room {
		id: uuid.New(),
		host: nil,
		register: make(chan *Client),
		unregister: make(chan *Client),
		clients: make(map[uuid.UUID]*Client),
		messageQueue: make(chan Message),
		clientPing: make(map[uuid.UUID]int),
		clientPingMeasure: make(map[uuid.UUID][]int),
		clientLastPing: make(map[uuid.UUID]time.Time),
	}
}

func (room *Room) Start() {
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
			client.conn.WriteJSON(Message{ Action: "event", Event: Event{ Name: "welcome", Data: struct { RoomID uuid.UUID `json:"roomID"`; ClientID uuid.UUID `json:"clientID"`}{client.room.id, client.id}}})
			for _, other := range room.clients {
				if other != client {
					other.conn.WriteJSON(Message{ Action: "message", Event: Event{Name: "Connect", Data: "New User Joined, ID: " + client.id.String()} })
				}
			}
			break
		case client := <-room.unregister:
			if _, ok := room.clients[client.id]; ok {
				delete(room.clients, client.id)
				close(client.send)
			}
			fmt.Println("Size of connection room: ", len(room.clients))
			for _, client := range room.clients {
				client.conn.WriteJSON(Message{ Action: "message", Event: Event{Name: "Disconnect", Data: "User Disconnected"} })
			}
			break
		case message := <-room.messageQueue:
			room.messageController(message);
		}
	}
}

// messageController following front controller pattern.
func (room *Room) messageController(message Message) {
	eventName := message.Event.Name
	newMessage := EventHandler(eventName, message, room)

	if newMessage != nil {
		room.dispatcher(*newMessage)
	}
}

func (room *Room) dispatcher(message Message) {
	fmt.Printf("%+v\n", message)
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