package websocket 

import (
	"fmt"
)

type playback struct {
	timestamp int64
	currVideo Video
	message Message
	action string
}

func newPlayback(message Message, room *Room, action string) playback {
	new := playback{
		timestamp: int64(message.Data.(float64)),
		currVideo: room.Video.Curr,
		message: message,
		action: action,
	}
	return new
}

func (p playback) handle() (Message, bool) {
	switch p.action {
	case "play":
		return p.play()
	case "pause":
		return p.pause()
	case "seekTo":
		return p.seekTo()
	default:
		return p.message, true
	}
}

func (p playback) play() (Message, bool)  {
	p.currVideo.Timer.SeekTo(p.timestamp).Play()
	fmt.Printf("Play(), seconds elapsed: %2f\n", float64(p.currVideo.Timer.Elapsed()) / 1000.0)
	p.currVideo.IsPlaying = true
	return p.message, true
}

func (p playback) pause() (Message, bool) {
	p.currVideo.Timer.SeekTo(p.timestamp).Pause()
	fmt.Printf("Pause(), seconds elapsed: %2f\n", float64(p.currVideo.Timer.Elapsed()) / 1000.0)
	p.currVideo.IsPlaying = false
	return p.message, true
}

func (p playback) seekTo() (Message, bool) {
	p.currVideo.Timer.SeekTo(p.timestamp)
	fmt.Printf("SeekTo(), seek to second: %2f\n", float64(p.currVideo.Timer.Elapsed()) / 1000.0)
	return p.message, true
}