package websocket

type handler interface {
	handle() (Message, bool)
}