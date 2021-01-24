package websocket

import (
	"fmt"
	"time"
	"log"
	"github.com/gorilla/websocket"
	"github.com/google/uuid"
)

// Client struct for identifying individual socket connection
type Client struct {
	ID uuid.UUID
	// the websocket connection
	Conn *websocket.Conn
	Room *Room
	// buffered channel of outbound messages
	Send chan Message
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


func (c *Client) Read() {
	defer func() {
		c.Room.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		// convert from JSON to Message
		fmt.Println(string(p))
		msg, err := UnmarshalJSONMessage([]byte(string(p)))

		if err != nil {
			fmt.Printf("Message received but unable to unmarshal to JSON\n")
		}
		
		c.Room.Broadcast <- *msg
	}
}

func (c *Client) readPump() {
	defer func() {
		c.Room.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, p, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		
		// convert from JSON to Message struct.
		fmt.Println(string(p))
		msg, err := UnmarshalJSONMessage([]byte(string(p)))

		if err != nil {
			fmt.Printf("Message received but unable to unmarshal to JSON\n")
		}
		
		c.Room.Broadcast <- *msg
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The room has closed the channel.
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.Conn.WriteJSON(message); err != nil {
				log.Printf("error: %v", err)
			}

			// Add queued chat messages to the current websocket message.
			n := len(c.Send)
			for i := 0; i < n; i++ {
				c.Conn.WriteJSON(<-c.Send)
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}

	}
}