// This file handles playback in a room ensuring all clients are synced with host.

package websocket

import "fmt"

type playback struct {
	message Message
	room *Room
	action string
}

func (p playback) execute() *Message {
	switch p.action {
	case "play":
		return p.play()
	case "pause":
		return p.pause()
	case "seekTo":
		return p.seekTo()
	default:
		return &p.message
	}
}

func (p playback) play() *Message {
	p.room.video.timer.SeekTo(int64(p.message.Event.Data.(float64))).Play()
	fmt.Printf("Play(), seconds elapsed: %2f\n", float64(p.room.video.timer.Elapsed()) / 1000.0)
	p.room.video.isPlaying = true
	return &p.message
}

func (p playback) pause() *Message {
	p.room.video.timer.SeekTo(int64(p.message.Event.Data.(float64))).Pause()
	fmt.Printf("Pause(), seconds elapsed: %2f\n", float64(p.room.video.timer.Elapsed()) / 1000.0)
	p.room.video.isPlaying = false
	return &p.message
}

func (p playback) seekTo() *Message {
	p.room.video.timer.SeekTo(int64(p.message.Event.Data.(float64)))
	fmt.Printf("SeekTo(), seek to second: %2f\n", float64(p.room.video.timer.Elapsed()) / 1000.0)
	return &p.message
}