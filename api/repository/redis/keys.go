package redis

const (
	roomKey = "room"
	userKey = "user"
	playerKey = "player"
	queueKey = "queue"
	currVideoKey = "currVideo"
	clockKey = "clock"
)

func newKey(keys ...string) string {
	final := ""
	for index, key := range keys {
		if index < 1 {
			final = key
		} else {
			final += ":" + key
		}
	}
	return final
}