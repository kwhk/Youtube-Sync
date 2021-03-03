package timer 

import "time"

// VideoTimer tracks the timing of the current playing video
type VideoTimer struct {
	Start time.Time
	Progress time.Duration
	stop bool
}

// Pause pauses the timer for the currently playing video
func (t *VideoTimer) Pause() *VideoTimer {
	t.Progress = time.Since(t.Start) + t.Progress
	t.Start = time.Now()
	t.stop = true
	return t
}

// Play starts the timer for the currently playing video
func (t *VideoTimer) Play() *VideoTimer {
	t.Start = time.Now()
	t.stop = false
	return t
}

// Elapsed returns time elapsed since video start in milliseconds
func (t *VideoTimer) Elapsed() int64 {

	if t.stop {
		return t.Progress.Milliseconds()
	} else {
		elapsed := t.Progress + time.Since(t.Start)
		return elapsed.Milliseconds()
	}
}

// SeekTo sets time elapsed to the provided millisecond
func (t *VideoTimer) SeekTo(ms int64) *VideoTimer {
	t.Progress = time.Duration(ms) * time.Millisecond
	t.Start = time.Now()
	return t
}