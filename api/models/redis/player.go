package redis

type Encodable interface {
	Encode() []byte
}

type PlayerRepository interface {
	AddToVideoQueue(roomID string, video Encodable) bool
	RemoveFromVideoQueue(roomID string, video Encodable) bool
	EmptyVideoQueue(roomID string) bool
	GetVideoQueue(roomID string) []Encodable
	ReorderVideoQueue(roomID string, newIndex int) bool
	GetCurrVideo(roomID string) (Encodable, bool)
	SetCurrVideo(roomID string, video Encodable) bool
	SetClock(roomID string, clock Encodable) bool
	GetClock(roomID string) (Encodable, bool)
}