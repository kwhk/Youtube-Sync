package websocket

type eventExecutor interface {
	execute() *Message
}

type ping struct {
	message Message
	room *Room
}

type playbackControl struct {
	message Message
	room *Room
}

func EventHandler(event string, message Message, room *Room) *Message {
	var executor eventExecutor
	switch event {
	case "ping":
		executor = ping{message, room}
	// host only events
	case "play", "pause", "seekTo":
		executor = playbackControl{message, room}
	}

	return executor.execute()
}

func (p ping) execute() *Message {
	client := p.room.clients[*p.message.Source]
	var message Message = client.measurePing(p.message)
	return &message
}

func (p playbackControl) execute() *Message {
	client := p.room.clients[*p.message.Source]

	if client == p.room.host {
		return &p.message
	}

	return nil
}