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
	id uuid.UUID
	// the websocket connection
	conn *websocket.Conn
	room *Room
	// buffered channel of outbound messages
	send chan Message
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

func (c *Client) readPump() {
	defer func() {
		c.room.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, p, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		
		// convert from JSON to Message struct.
		msg, err := UnmarshalJSONMessage([]byte(string(p)))

		if err != nil {
			fmt.Printf("Message received but unable to unmarshal to JSON\n")
		}
		
		c.room.messageQueue <- *msg
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The room has closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteJSON(message); err != nil {
				log.Printf("error: %v", err)
				log.Printf("message: %+v\n", message)
			}

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				c.conn.WriteJSON(<-c.send)
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}

	}
}

func (c *Client) measurePing(msg Message) Message {
	// MAX_PING must match clientPingMeasure array size
	MAX_PING_NUM := 10

	type Handshake struct {
		// All caps represents whether ACK or SYN has been set
		ACK int `json:"ACK"`
		SYN int `json:"SYN"`
		FIN int `json:"FIN"`
		Seq int `json:"seq"`
		Ack int `json:"ack"`
	}

	var newMsg Message
	var hs map[string]interface{} = msg.Event.Data.(map[string]interface{})

	// Convert all values in handshake from float to int.
	for key, val := range hs {
		hs[key] = int(val.(float64))
	}

	// If SYN bit is 1 and ACK bit has not been set been client, then this is the request for handshake with server.
	if hs["SYN"] == 1 && hs["ACK"] == 0 {
		newMsg = Message{ Action: "event", Source: &c.id, Event: Event{ Name: "ping", Data: Handshake{ACK: 1, SYN: 1, Seq: 0, Ack: hs["seq"].(int) + 1}}}
		c.room.clientLastPing[c.id] = time.Now()
	// If SYN bit and ACK bit have been set by client, then handshake is complete and is ready for ping measurements.
	} else if hs["ACK"] == 1 {
		currTime := time.Now()
		ping := currTime.Sub(c.room.clientLastPing[c.id]).Milliseconds()
		c.room.clientLastPing[c.id] = currTime

		if arr, ok := c.room.clientPingMeasure[c.id]; ok {
			// -1 because initial seq sent (seq = 0) is to request to start ping measurement.
			arr[hs["seq"].(int) - 1] = int(ping)
		}

		// Calculate average after MAX_PING tries.
		if hs["seq"].(int) >= MAX_PING_NUM{
			sum := 0
			for i := 0; i < MAX_PING_NUM; i++ {
				sum += c.room.clientPingMeasure[c.id][i]
				// Reset array.
				c.room.clientPingMeasure[c.id][i] = 0
			}

			// Divide sum by 2 because ping stored in clientPingMeasure measures latency for round trip and not one way.
			c.room.clientPing[c.id] = (sum / 2) / MAX_PING_NUM
			fmt.Printf("Ping is %d ms\n", c.room.clientPing[c.id])
			newMsg = Message{ Action: "event", Source: &c.id, Event: Event{ Name: "ping", Data: Handshake{FIN: 1}}}
		} else {
			newMsg = Message{ Action: "event", Source: &c.id, Event: Event{ Name: "ping", Data: Handshake{ACK: 1, Seq: hs["ack"].(int), Ack: hs["seq"].(int) + 1}}}
		}
	}

	return newMsg
}