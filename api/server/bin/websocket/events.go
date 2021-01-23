package websocket

// import (
// 	"fmt"
// 	"encoding/json"
// )


// EventController reads message and outputs according to message value
// func EventController(event Event) Message {
// 	// assert that eventName exists and is of type string
// 	switch event.Name {
// 		case "seekTo":
// 			const secs = event.Data
// 			fmt.Printf(" to seek to %d.\n", secs)
// 			return Message{ Type: 1, Body:  Data{eventName string; data int}{"seekTo", secs} }
// 		case "pause":
// 			fmt.Printf(" to pause.\n")
// 			return Message{ Type: 1, Body: struct {eventName string}{"pause"}}
// 		case "play":
// 			fmt.Printf(" to play.\n")
// 			return Message{ Type: 1, Body: struct {eventName string}{"play"}}
// 	}

// 	return msg
// }