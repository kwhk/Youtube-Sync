package websocket

import (
	"encoding/json"
	"log"
)

type Video struct {
	// URL of video.
	URL string `json:"url"`
	// Duration of video in ms.
	Duration int64 `json:"duration"`
	// Status to notify joining users playback state.
	IsPlaying bool `json:"isPlaying"`
}

func (v Video) encode() []byte {
	json, err := json.Marshal(v)
	if err != nil {
		log.Println(err)
	}

	return json
}

func decodeVideo(p []byte) Video {
	var video Video
	if err := json.Unmarshal(p, &video); err != nil {
		log.Printf("Error on unmarshal video: %s\n", err)
	}

	return video
}

func (v Video) GetEncoding() []byte {
	return v.encode()
}