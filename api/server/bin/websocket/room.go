package websocket

import (
	"fmt"
	"github.com/google/uuid"
)

func NewRoom() *Room {
	return &Room {
		ID: uuid.New(),
		Register: make(chan *Client),
		Unregister: make(chan *Client),
		Clients: make(map[*Client]bool),
		Broadcast: make(chan Message),
	}
}

func (room *Room) Start() {
	for {
		select {
		case client := <-room.Register:
			room.Clients[client] = true
			fmt.Println("Size of connection room: ", len(room.Clients))
			client.Conn.WriteJSON(Message{ Type: "message", Body: Event{ Name: "welcome", Data: struct { RoomID uuid.UUID `json:"roomID"`; ClientID uuid.UUID `json:"clientID"`}{client.Room.ID, client.ID}}})
			for others := range room.Clients {
				if others != client {
					others.Conn.WriteJSON(Message{ Type: "message", Body: "New User Joined, ID: " + client.ID.String() })
				}
			}
			break
		case client := <-room.Unregister:
			if _, ok := room.Clients[client]; ok {
				delete(room.Clients, client)
				close(client.Send)
			}
			fmt.Println("Size of connection room: ", len(room.Clients))
			for client := range room.Clients {
				client.Conn.WriteJSON(Message{ Type: "message", Body: "User Disconnected" })
			}
			break
		case message := <-room.Broadcast:			
			fmt.Printf("Sending message to all clients in room %s...\n", room.ID.String());
			room.messageController(message);
		}
	}
}

func (room *Room) messageController(message Message) {
	var target *Target = message.Target

	if target.IncludeSender == false {
		room.emit(message, *target.SourceClientID)
	} else {
		room.broadcast(message)
	}
}

func (room *Room) broadcast(message Message) {
	for client := range room.Clients {
		select {
		case client.Send <- message:
		default:
			// Done, no more messages to send.
			close(client.Send)
			delete(room.Clients, client)
		}
	}
}

func (room *Room) emit(message Message, srcClientID uuid.UUID) {
	for client := range room.Clients {
		if client.ID != srcClientID {
			select {
			case client.Send <- message:
			default:
				close(client.Send)
				delete(room.Clients, client)
			}
		}
	}
}