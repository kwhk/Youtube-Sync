package clock 

import (
	"log"
	"time"
	"encoding/json"
)

// Clock tracks the timing of the current playing video
type Clock struct {
	Start time.Time `json:"start"`
	Progress time.Duration `json:"progress"`
	Stop bool `json:"stop"`
}

// Pause pauses the timer for the currently playing video
func (t *Clock) Pause() *Clock {
	t.Progress = time.Since(t.Start) + t.Progress
	t.Start = time.Now()
	t.Stop = true
	return t
}

// Play starts the timer for the currently playing video
func (t *Clock) Play() *Clock {
	t.Start = time.Now()
	t.Stop = false
	return t
}

// Elapsed returns time elapsed since video start in milliseconds
func (t *Clock) Elapsed() int64 {

	if t.Stop {
		return t.Progress.Milliseconds()
	} else {
		elapsed := t.Progress + time.Since(t.Start)
		return elapsed.Milliseconds()
	}
}

// SeekTo sets time elapsed to the provided millisecond
func (t *Clock) SeekTo(ms int64) *Clock {
	t.Progress = time.Duration(ms) * time.Millisecond
	t.Start = time.Now()
	return t
}

func (t *Clock) encode() []byte {
	json, err := json.Marshal(t)
	if err != nil {
		log.Println(err)
	}
	
	return json
}

func DecodeClock(p []byte) *Clock {
	var clock Clock
	if err := json.Unmarshal(p, &clock); err != nil {
		log.Printf("Error on unmarshal clock: %s\n", err)
	}

	return &clock
}

func (t *Clock) GetEncoding() []byte {
	return t.encode()
}