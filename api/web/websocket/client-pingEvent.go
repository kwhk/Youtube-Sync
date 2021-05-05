package websocket 
import (
	"time"
	"fmt"
)

type ping struct {
	client *Client
	data map[string]int
}

type handshakeJSON struct {
	// All caps represents whether ACK or SYN has been set
	ACK int `json:"ACK"`
	SYN int `json:"SYN"`
	FIN int `json:"FIN"`
	Seq int `json:"seq"`
	Ack int `json:"ack"`
}

func newPing(message Message, client *Client) ping {
	data := message.Data.(map[string]interface{})
	var d map[string]int = make(map[string]int)
	for key, val := range data {
		d[key] = int(val.(float64))
	}

	new := ping{
		client: client,
		data: d,
	}
	return new
}

func (p ping) handle() (Message, bool) {
	// MAX_PING must match clientPingMeasure array size
	MAX_PING_NUM := 10

	var newMsg Message = Message{ Action: UserPingAction, Sender: p.client}
	hs := p.data

	// If SYN bit is 1 and ACK bit has not been set been client, then this is the request for handshake with server.
	if hs["SYN"] == 1 && hs["ACK"] == 0 {
		newMsg.Data = handshakeJSON{ACK: 1, SYN: 1, Seq: 0, Ack: hs["seq"] + 1}
		p.client.lastPing = time.Now()
	// If SYN bit and ACK bit have been set by client, then handshake is complete and is ready for ping measurements.
	} else if hs["ACK"] == 1 {
		currTime := time.Now()
		ping := currTime.Sub(p.client.lastPing).Milliseconds()
		p.client.lastPing = currTime

		// -1 because initial seq sent (seq = 0) is to request to start ping measurement.
		// p.client.pingMeasure[hs["seq"] - 1] = int(ping)
		arr := append(p.client.pingMeasure, int(ping))
		p.client.pingMeasure = arr

		// Calculate average after MAX_PING tries.
		if hs["seq"] >= MAX_PING_NUM{
			sum := 0

			for _, ping := range p.client.pingMeasure {
				sum += ping
			}

			// Reset measuring array.
			p.client.pingMeasure = p.client.pingMeasure[:0]

			// Divide sum by 2 because ping stored in clientPingMeasure measures latency for round trip and not one way.
			p.client.ping = (sum / 2) / MAX_PING_NUM
			fmt.Printf("Ping is %d ms\n", p.client.ping)
			newMsg.Data = handshakeJSON{FIN: 1}
		} else {
			newMsg.Data = handshakeJSON{ACK: 1, Seq: hs["ack"], Ack: hs["seq"] + 1}
		}
	}

	return newMsg, true
}