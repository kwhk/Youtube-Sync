// This file handles which event should be executed depending on the event name sent within a Message.

package websocket

type eventExecutor interface {
	execute() *Message
}

type ping struct {
	message Message
	room *Room
}


func eventHandler(event string, message Message, room *Room) *Message {
	var executor eventExecutor
	switch event {
	case "ping":
		executor = ping{message, room}
	// host only events
	case "play", "pause", "seekTo":
		executor = Playback{message, room, event}
	}

	return executor.execute()
}

func (p ping) execute() *Message {
	client := p.room.clients[*p.message.Source]
	var message Message = client.measurePing(p.message)
	return &message
}