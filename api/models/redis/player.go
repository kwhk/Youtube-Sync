package redis

type Clock interface {
	GetEncoding() []byte
}

type Video interface {
	GetEncoding() []byte
}

type PlayerRepository interface {
	AddToVideoQueue(video Video) bool
	RemoveFromVideoQueue(video Video) bool
	EmptyVideoQueue() bool
	GetCurrVideo() (Video, bool)
	SetCurrVideo(video Video) bool
	SetClock(clock Clock) bool
	GetClock() (Clock, bool)
}