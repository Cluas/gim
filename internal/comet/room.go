package comet

import (
	"encoding/json"
	"errors"
	"sync"
)

// ErrRoomIsDropped is err for room is dropped
var ErrRoomIsDropped = errors.New("room is dropped")

// NO_ROOM is flag
var NoRoom = int32(-1)

type Proto struct {
	Ver       int16 `json:"ver"` // protocol version
	Operation int32 `json:"op"`  // operation for request
	// SeqId     int32           `json:"seq"`  // sequence number chosen by client
	Body json.RawMessage `json:"body"` // binary body bytes(json.RawMessage is []byte)

}

type RoomMsgArg struct {
	RoomID int32
	P      Proto
}

// Room is struct of Room
type Room struct {
	ID        int32 // Room ID
	rLock     sync.RWMutex
	next      *Channel
	isDropped bool
	Online    int
}

// NewRoom is constructor of Room
func NewRoom(ID int32) *Room {
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
