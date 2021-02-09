// This file handles which event should be executed depending on the event name sent within a Message.

package websocket

type eventExecutor interface {
	execute() (Message, bool)
}

func eventHandler(event string, message Message, room *Room) (Message, bool) {
	var executor eventExecutor
	switch event {
	case "ping":
		executor = ping{message, room}
	// host only events
	case "play", "pause", "seekTo":
		executor = playback{message, room, event}
	case "addVideoQueue", "popVideoQueue", "emptyVideoQueue":
		executor = videoQueue{message, room, event}
	}

	return executor.execute()
}
