package websocket
import (
	"time"
	"fmt"
)

type ping struct {
	message Message
	room *Room
}

func (p ping) execute() *Message {
	c := p.room.clients[*p.message.Source]
	msg := p.message

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

	return &newMsg
}