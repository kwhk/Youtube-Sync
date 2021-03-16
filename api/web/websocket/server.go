package websocket

import (
	"encoding/json"
	"log"

	wp "github.com/kwhk/sync/api/utils/workerPool"

	"github.com/google/uuid"
	"github.com/kwhk/sync/api/config"
	"github.com/kwhk/sync/api/models"
)

const PubSubGeneralChannel = "general"

type WsServer struct {
	users []models.User
	clients map[string]*Client
	register chan *Client
	unregister chan *Client
	emit chan Message
	rooms map[uuid.UUID]*Room
	workerPool *wp.Pool
	roomRepository models.RoomRepository
	userRepository models.UserRepository
}

func NewWebsocketServer(roomRepo models.RoomRepository, userRepo models.UserRepository) *WsServer {
	wsServer := &WsServer{
		clients: make(map[string]*Client),
		register: make(chan *Client),
		unregister: make(chan *Client),
		emit: make(chan Message),
		rooms: make(map[uuid.UUID]*Room),
		roomRepository: roomRepo,
		userRepository: userRepo,
		workerPool: wp.NewPool(10),
	}

	wsServer.users = userRepo.GetAllUsers()

	return wsServer
}

func (server *WsServer) Run() {
	go server.listenPubSubChannel()
	go server.workerPool.Run()

	for {
		select {
		case client := <-server.register:
			server.workerPool.AddJob(func() {server.registerClient(client)})
		case client := <-server.unregister:
			server.workerPool.AddJob(func() {server.unregisterClient(client)})
		case message := <-server.emit:
			server.workerPool.AddJob(func() {server.eventHandler(message)})
		}
	}
}

func (server *WsServer) registerClient(client *Client) {
	// Add user to the repo
	server.userRepository.AddUser(client)

	// Publish user in PubSub
	server.publishClientJoined(client)

	server.clients[client.GetID()] = client
}

func (server *WsServer) unregisterClient(client *Client) {
	if _, ok := server.clients[client.GetID()]; ok {
		delete(server.clients, client.GetID())

		// Remove user from repo
		server.userRepository.DeleteUser(client)

		// Publish user left in pub/sub
		server.publishClientLeft(client)
	}
}

func (server *WsServer) publishClientJoined(client *Client) {
	message := &Message{
		Action: UserJoinAction,
		Sender: client,
	}

	if err := config.Redis.Publish(ctx, PubSubGeneralChannel, message.encode()).Err(); err != nil {
		log.Println(err)
	}
}

func (server *WsServer) publishClientLeft(client *Client) {
	message := &Message{
		Action: UserLeaveAction,
		Sender: client,
	}

	if err := config.Redis.Publish(ctx, PubSubGeneralChannel, message.encode()).Err(); err != nil {
		log.Println(err)
	}
}

func (server *WsServer) listenPubSubChannel() {
	pubsub := config.Redis.Subscribe(ctx, PubSubGeneralChannel)
	
	ch := pubsub.Channel()

	for message := range ch {
		var msg Message
		if err := json.Unmarshal([]byte(message.Payload), &msg); err != nil {
			log.Printf("Error on umarshalJSON message %s\n", err)
			return
		}
		server.eventHandler(msg)
	}
}

func (server *WsServer) eventHandler(message Message) {
	switch message.Action {
	case UserJoinAction:
		server.handleUserJoined(message)
	case UserLeaveAction:
		server.handleUserLeft(message)
	default:
		log.Printf("WsServer eventHandler does not recognize event '%s'\n", message.Action)
	}
}

func (server *WsServer) handleUserJoined(message Message) {
	// Add the user to the slice
	server.users = append(server.users, message.Sender)
	server.broadcastToClients(message)
}

func (server *WsServer) handleUserLeft(message Message) {
	// Remove the user from the slice
	for i, user := range server.users {
		if user.GetID() == message.Sender.GetID() {
			server.users[i] = server.users[len(server.users)-1]
			server.users = server.users[:len(server.users)-1]
		}
	}

	server.broadcastToClients(message)
}

func (server *WsServer) emitToClient(message Message) {
	server.clients[message.Sender.GetID()].send <- message.encode()
}

func (server *WsServer) broadcastToClients(message Message) {
	msg := message.encode()
	for _, client := range server.clients {
		client.send <- msg
	}
}

func (server *WsServer) findRoomByID(ID uuid.UUID) *Room {
	if elem, ok := server.rooms[ID]; ok {
		return elem
	}
	
	return nil
}

func (server *WsServer) findUserByID(ID uuid.UUID) models.User {
	var foundUser models.User

	for _, client := range server.users {
		if client.GetID() == ID.String() {
			foundUser = client
			break
		}
	}

	return foundUser
}

func (server *WsServer) createRoom(name string, private bool) *Room {
	room := newRoom(name, private, server)
	server.roomRepository.AddRoom(room)

	go room.run()
	server.rooms[room.ID] = room

	return room
}

func (server *WsServer) deleteRoom(room *Room) {
	server.roomRepository.DeleteRoom(room)
	delete(server.rooms, room.ID)
}