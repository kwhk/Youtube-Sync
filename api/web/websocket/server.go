package websocket

import (
	"fmt"
	"encoding/json"
	"log"

	wp "github.com/kwhk/sync/api/utils/workerPool"

	"github.com/kwhk/sync/api/config"
	models "github.com/kwhk/sync/api/models/redis"
	repo "github.com/kwhk/sync/api/repository/redis"
)

const PubSubGeneralChannel = "general"

type WsServer struct {
	users []models.User
	clients map[string]*Client
	register chan *Client
	unregister chan *Client
	emit chan Message
	rooms map[string]*Room
	workerPool *wp.Pool
	roomRepository repo.RoomRepository
	userRepository repo.UserRepository
	playerRepository repo.PlayerRepository
}

func NewWebsocketServer(roomRepo repo.RoomRepository, userRepo repo.UserRepository, playerRepo repo.PlayerRepository) *WsServer {
	wsServer := &WsServer{
		users: make([]models.User, 0),
		clients: make(map[string]*Client),
		register: make(chan *Client),
		unregister: make(chan *Client),
		emit: make(chan Message),
		rooms: make(map[string]*Room),
		roomRepository: roomRepo,
		userRepository: userRepo,
		playerRepository: playerRepo,
		workerPool: wp.NewPool(10),
	}

	return wsServer
}

func (server *WsServer) Run() {
	go server.listenPubSubChannel()
	go server.workerPool.Run()

	for {
		select {
		case client := <-server.register:
			server.registerClient(client)
		case client := <-server.unregister:
			server.unregisterClient(client)
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

func (server *WsServer) findRoomByID(id string) *Room {
	// Check if room already exists in instance.
	if elem, ok := server.rooms[id]; ok {
		return elem
	}

	// Fetch room from database and then store in instance.
	if r, ok := server.roomRepository.FindRoomByID(id); ok {
		room := newRoomFromRedis(r, server)
		go room.run()
		server.rooms[room.ID] = room
		return room
	}
	
	return nil
}

func (server *WsServer) findUserByID(id string) models.User {
	var foundUser models.User

	for _, client := range server.users {
		if client.GetID() == id {
			foundUser = client
			break
		}
	}

	return foundUser
}

func (server *WsServer) createRoom() *Room {
	// Create new instance of room.
	room := newRoom(server)

	// Add instance to database.
	if ok := server.roomRepository.AddRoom(room); ok {
		fmt.Printf("Room %s created\n", room.ID)
		go room.run()
		server.rooms[room.ID] = room
		return room
	}

	return nil
}

func (server *WsServer) deleteRoom(room *Room) {
	// Delete from database
	server.roomRepository.DeleteRoom(room)
	delete(server.rooms, room.ID)
}