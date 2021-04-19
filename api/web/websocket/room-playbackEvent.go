package websocket

import (
	"fmt"

	"github.com/kwhk/sync/api/repository/redis"
	"github.com/kwhk/sync/api/utils/clock"
)

type playback struct {
	timestamp int64
	clock *clock.Clock
	message Message
	action string
	playerRepo redis.PlayerRepository
	roomID string
}

func newPlayback(message Message, room *Room, action string) playback {
	new := playback{
		timestamp: int64(message.Data.(float64)),
		clock: room.Clock,
		message: message,
		action: action,
		playerRepo: room.wsServer.playerRepository,
		roomID: room.ID,
	}
	return new
}

func (p playback) handle() (Message, bool) {
	switch p.action {
	case PlayVideoAction:
		return p.play()
	case PauseVideoAction:
		return p.pause()
	case SeekToVideoAction:
		return p.seekTo()
	default:
		return p.message, true
	}
}

func (p playback) play() (Message, bool)  {
	p.clock.SeekTo(p.timestamp).Play()
	fmt.Printf("Play(), seconds elapsed: %2f\n", float64(p.clock.Elapsed()) / 1000.0)
	p.playerRepo.SetClock(p.roomID, p.clock)
	return p.message, true
}

func (p playback) pause() (Message, bool) {
	p.clock.SeekTo(p.timestamp).Pause()
	fmt.Printf("Pause(), seconds elapsed: %2f\n", float64(p.clock.Elapsed()) / 1000.0)
	p.playerRepo.SetClock(p.roomID, p.clock)
	return p.message, true
}

func (p playback) seekTo() (Message, bool) {
	p.clock.SeekTo(p.timestamp)
	fmt.Printf("SeekTo(), seek to second: %2f\n", float64(p.clock.Elapsed()) / 1000.0)
	p.playerRepo.SetClock(p.roomID, p.clock)
	return p.message, true
}