package websocket

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
)

type Event struct {
	Name string `json:"eventName"`
	Data interface{} `json:"data,omitempty"`
}

// Message struct storing message content and the message type
type Message struct {
	// Message type can either be a message or an event
	Type string `json:"messageType"`
	// Target can either be to an individual, to 
	Target *Target `json:"target,omitempty"`
	Source *uuid.UUID `json:"sourceClientID,omitempty"`
	Body interface{} `json:"body"`
}

type Target struct {
	SourceClientID *uuid.UUID `json:"sourceClientID,omitempty"`
	UserID *uuid.UUID `json:"userID,omitempty"`
	RoomID *uuid.UUID `json:"roomID,omitempty"`
	IncludeSender bool `json:"includeSender"`
}

// UnmarshalJSONMessage for message to discern different message types and unmarshal
func UnmarshalJSONMessage(message []byte) (*Message, error) {
	var body json.RawMessage
	msg := Message{
		Body: &body,
	}
	
	if err := json.Unmarshal(message, &msg); err != nil {
		return nil, err
	}

	switch msg.Type {
		case "event":
			var e Event
			if err := json.Unmarshal(body, &e); err != nil {
				return nil, err
			}
		case "message":
			var s string
			if err := json.Unmarshal(body, &s); err != nil {
				return nil, err
			}
		default:
			return nil, errors.New("Invalid JSON object value, cannot convert")
	}

	// fmt.Printf("Message: %+v\n", msg)

	return &msg, nil
}