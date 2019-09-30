package comet

import (
	"errors"
	"sync"
)

// ErrRoomIsDroped is err for room is droped
var ErrRoomIsDroped = errors.New("room is droped")

// NOROOM is flag
var NOROOM = int32(-1)

// Room is struct of Room
type Room struct {
	ID       int32 // Room ID
	rLock    sync.RWMutex
	next     *Channel
	isDroped bool
	Online   int
}

// NewRoom is constructor of Room
func NewRoom(ID int32) *Room {
	return &Room{
		ID:       ID,
		next:     nil,
		isDroped: false,
		Online:   0,
	}
}

// Put is func for add channel
func (r *Room) Put(ch *Channel) error {
	if !r.isDroped {
		if r.next != nil {
			r.next.Prev = ch
		}
		ch.Next = r.next
		ch.Prev = nil
		r.next = ch
		r.Online++
		return nil
	}
	return ErrRoomIsDroped
}
