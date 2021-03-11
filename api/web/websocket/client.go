package websocket

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"time"
	"encoding/json"
	"sync"

	"github.com/kwhk/sync/api/models"
)

var clientWg sync.WaitGroup

// Client struct for identifying individual socket connection
type Client struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	// the websocket connection
	conn *websocket.Conn
	room *Room
	// buffered channel of outbound messages
	send     chan []byte
	wsServer *WsServer

	ping int
	pingMeasure []int
	lastPing time.Time
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space = []byte{' '}
)

func initClient(conn *websocket.Conn, wsServer *WsServer) *Client {
	client := &Client{
		ID:       uuid.New(),
		conn:     conn,
		wsServer: wsServer,
		send:     make(chan []byte, 256),
		ping: 0,
		pingMeasure: make([]int, 0, 10),
		lastPing: time.Now(), 
	}

	return client
}

func (client *Client) disconnect() {
	client.wsServer.unregister <-client
	client.handleLeaveRoom()
	close(client.send)
	client.conn.Close()
}

func (client *Client) readPump() {
	defer func() {
		client.disconnect()
	}()

	client.conn.SetReadLimit(maxMessageSize)
	client.conn.SetReadDeadline(time.Now().Add(pongWait))
	client.conn.SetPongHandler(func(string) error { client.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, p, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		// convert from JSON to Message struct.
		var message Message
		if err := json.Unmarshal(p, &message); err != nil {
			log.Printf("Error on unmarshal JSON message %s", err)
			continue
		}

		message.Sender = client
		client.wsServer.workerPool.AddJob(func() {client.eventHandler(message)})
	}
}

func (client *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		client.conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.send:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The room has closed the channel.
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}

	}
}

func (client *Client) eventHandler(message Message) {
	message.Sender = client
	action := message.Action

	switch action {
	case JoinRoomAction:
		client.handleJoinRoom(message)
	case LeaveRoomAction:
		client.handleLeaveRoom()
	case CreateRoomAction:
		client.handleCreateRoom()
	case UserPingAction:
		client.handlePing(message)
	default:
		room := client.room
		if room != nil {
			room.broadcast <- message
		}
	}
}

func (client *Client) handlePing(message Message) {
	if msg, ok := newPing(message, client).handle(); ok {
		client.send <- msg.encode()
	}
}

func (client *Client) handleCreateRoom() {
	// private attribute set to false and no room name for the meantime
	room := client.wsServer.createRoom("", false)
	client.joinRoom(room.ID, client)
	
	message := Message{
		Action: CreateRoomAction,
		Data: room,
	}

	client.send <- message.encode()
}

func (client *Client) handleJoinRoom(message Message) {
	roomID := uuid.MustParse(message.Data.(string))
	room := client.joinRoom(roomID, nil)

	if room == nil {
		return
	}
	
	client.notifyRoomJoined(room)
}

func (client *Client) handleLeaveRoom() {
	room := client.room
	
	if room == nil {
		return
	}
	
	room.unregister <- client
	client.room = nil
}

func (client *Client) joinRoom(roomID uuid.UUID, sender models.User) *Room {
	room := client.wsServer.findRoomByID(roomID)

	if room == nil {
		log.Printf("RoomID %s not found.\n", roomID.String())
		return nil
	}

	if sender == nil && room.Private {
		return nil
	}

	// Check if user has already joined this room.
	if !client.isInRoom(room) {
		client.wsServer.userRepository.JoinRoom(client, room)
		client.room = room
		room.register <- client
		// client.notifyRoomJoined(room, sender)
	}

	return room
}

func (client *Client) isInRoom(room *Room) bool {
	return client.room == room
}

// Notify client that they have successfully joined room
func (client *Client) notifyRoomJoined(room *Room) {	
	clients := []string{}
	for id := range room.Clients {
		clients = append(clients, id)
	}

	var joinMsg Message = Message{ 
		Action: RoomWelcomeAction,
		Sender: client,
		Data: struct {
			VideoQueue []Video `json:"videoQueue"`
			ConnectedUsers []string `json:"connectedUsers"`
			CurrVideo struct {
				URL string `json:"url"`
				Elapsed int64 `json:"elapsed"`
				IsPlaying bool `json:"isPlaying"`
			} `json:"currVideo"`

		}{
			VideoQueue: room.Video.Queue,
			ConnectedUsers: clients,
			CurrVideo: struct {
				URL string `json:"url"`
				Elapsed int64 `json:"elapsed"`
				IsPlaying bool `json:"isPlaying"`
			}{
				room.Video.Curr.URL,
				room.Video.Curr.Timer.Elapsed(),
				room.Video.Curr.IsPlaying,
			},
		},
	}

	client.send <- joinMsg.encode()
}

func (client *Client) GetID() string {
	return client.ID.String()
}

func (client *Client) GetName() string {
	return client.Name
}
