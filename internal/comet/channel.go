package comet

import "github.com/gorilla/websocket"

// Channel is struct of channel
type Channel struct {
	Room      *Room
	broadcast chan []byte
	uid       string
	conn      *websocket.Conn
	Next      *Channel
	Prev      *Channel
}

// NewChannel is constructor of Channel
func NewChannel(svr int) *Channel {
	c := new(Channel)
	c.broadcast = make(chan []byte, svr)
	c.Next = nil
	c.Prev = nil
	return c
}
