package websocket

import (
	"encoding/json"
	"github.com/google/uuid"
)

// Event struct
type Event struct {
	Name string `json:"name"`
	Data interface{} `json:"data,omitempty"`
}

// Message struct storing message content and the message type
type Message struct {
	// Message type can either be a message or an event
	Action string `json:"action"`
	// Target can either be to an individual, to 
	Target *Target `json:"target,omitempty"`
	Source *uuid.UUID `json:"sourceClientID,omitempty"`
	Event Event `json:"event"`
}

// Target struct to identify client to direct message to
type Target struct {
	UserID *uuid.UUID `json:"userID,omitempty"`
	RoomID *uuid.UUID `json:"roomID,omitempty"`
	IncludeSender bool `json:"includeSender"`
}

// UnmarshalJSONMessage for message to discern different message types and unmarshal
func UnmarshalJSONMessage(message []byte) (*Message, error) {
	var msg Message
	
	if err := json.Unmarshal(message, &msg); err != nil {
		return nil, err
	}

	return &msg, nil
}