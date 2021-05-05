package websocket

import (
	"encoding/json"
	"log"

	models "github.com/kwhk/sync/api/models/redis"
)

const (
	// RoomWelcomeAction used to send data for new users that joined a room
	RoomWelcomeAction = "room-welcome"
	UserJoinAction = "user-join"
	UserLeaveAction = "user-left"
	JoinRoomAction = "join-room"
	LeaveRoomAction = "leave-room"
	CreateRoomAction = "create-room"
	DeleteRoomAction = "delete-room"
	SendMessageAction = "send-message"
	UserPingAction = "user-ping"
	PlayVideoAction = "play-video"
	PauseVideoAction = "pause-video"
	SeekToVideoAction = "seekto-video"
	AddVideoQueueAction = "add-video-queue"
	RemoveVideoQueueAction = "remove-video-queue"
	PlayVideoQueueAction = "play-video-queue"
	EmptyVideoQueueAction = "empty-video-queue"
)

type Message struct {
	Action string `json:"action"`
	Target *Room `json:"target"`
	Sender models.User `json:"sender"`
	Data interface{} `json:"data,omitempty"`
}

func (message *Message) UnmarshalJSON(data []byte) error {
	type Alias Message
	msg := &struct {
		Sender Client `json:"sender"`
		*Alias
	}{
		Alias: (*Alias)(message),
	}

	if err := json.Unmarshal(data, &msg); err != nil {
		return err
	}

	message.Sender = &msg.Sender
	return nil
}

func (message *Message) encode() []byte {
	json, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
	}

	return json
}