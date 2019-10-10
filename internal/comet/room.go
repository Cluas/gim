package comet

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/Cluas/gim/pkg/log"
)

// ErrRoomIsDropped is err for room is dropped
var ErrRoomIsDropped = errors.New("room is dropped")

// NoRoom is flag
var NoRoom = "NoRoom"

// Proto is struct for msg protocol
type Proto struct {
	Ver       int16           `json:"ver"`  // protocol version
	Operation int32           `json:"op"`   // operation for request
	Body      json.RawMessage `json:"body"` // binary body bytes(json.RawMessage is []byte)

}

// RoomMsgArg is struct of room msg
type RoomMsgArg struct {
	RoomID string
	P      Proto
}

// Room is struct of Room
type Room struct {
	ID        string // Room ID
	rLock     sync.RWMutex
	next      *Channel
	isDropped bool
	Online    int
}

// NewRoom is constructor of Room
func NewRoom(ID string) *Room {
	return &Room{
		ID:        ID,
		next:      nil,
		isDropped: false,
		Online:    0,
	}
}

// Put is func for add channel
func (r *Room) Put(ch *Channel) error {
	if !r.isDropped {
		if r.next != nil {
			r.next.Prev = ch
		}
		ch.Next = r.next
		ch.Prev = nil
		r.next = ch
		r.Online++
		return nil
	}
	return ErrRoomIsDropped
}

// Push is func for room push msg
func (r *Room) Push(p *Proto) {
	r.rLock.RLock()

	for ch := r.next; ch != nil; ch = ch.Next {

		// log.Infof("Room Push info %v", p)
		err := ch.Push(p)
		if err != nil {
			log.Errorf("Room Channel Push err: %v", err)
		}
	}

	r.rLock.RUnlock()
	return
}

// Del is func to del channel from room
func (r *Room) Del(ch *Channel) bool {
	r.rLock.RLock()
	if ch.Next != nil {
		//if not footer
		ch.Next.Prev = ch.Prev
	}

	if ch.Prev != nil {
		// if not header
		ch.Prev.Next = ch.Next
	} else {
		r.next = ch.Next
	}
	r.Online--
	r.isDropped = r.Online == 0
	r.rLock.RUnlock()

	return r.isDropped
}
