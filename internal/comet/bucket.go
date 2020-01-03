package comet

import (
	"sync"
	"sync/atomic"

	"github.com/Cluas/gim/pkg/log"
)

// BucketOptions is struct of Bucket Config
type BucketOptions struct {
	ChannelSize   int
	RoomSize      int
	RoutineAmount uint64
	RoutineSize   int
}

// Bucket is struct of Bucket
type Bucket struct {
	cLock       sync.RWMutex        // protect the channels for chs
	chs         map[string]*Channel // map sub key to a channel
	o           *BucketOptions      // bucket options
	rooms       map[string]*Room    // bucket room channel
	routines    []chan *RoomMsgArg
	routinesNum uint64
	broadcast   chan []byte
}

// NewBucket is constructor of Bucket
func NewBucket(o *BucketOptions) (b *Bucket) {
	b = new(Bucket)
	b.chs = make(map[string]*Channel, o.ChannelSize)
	b.o = o
	b.routines = make([]chan *RoomMsgArg, o.RoutineAmount)
	b.rooms = make(map[string]*Room, o.RoomSize)
	for i := uint64(0); i < b.o.RoutineAmount; i++ {
		c := make(chan *RoomMsgArg, o.RoutineSize)
		b.routines[i] = c
		go b.PushRoom(c)
	}
	return
}

// Put is func to add channel
func (b *Bucket) Put(uid string, rid string, ch *Channel) (err error) {
	var (
		room *Room
		ok   bool
	)
	b.cLock.Lock()

	if rid != NoRoom {
		if room, ok = b.rooms[rid]; !ok {
			room = NewRoom(rid)
			b.rooms[rid] = room
		}
		ch.Room = room
	}
	ch.uid = uid
	b.chs[uid] = ch
	b.cLock.Unlock()

	if room != nil {
		err = room.Put(ch)
	}
	return
}

// Channel is func to get Channel from Bucket by key
func (b *Bucket) Channel(key string) (ch *Channel) {
	// 读操作的锁定和解锁
	b.cLock.RLock()
	ch = b.chs[key]
	b.cLock.RUnlock()
	return
}

func (b *Bucket) delCh(ch *Channel) {
	var (
		ok   bool
		room *Room
	)
	b.cLock.RLock()

	if ch, ok = b.chs[ch.uid]; ok {
		room = b.chs[ch.uid].Room
		delete(b.chs, ch.uid)

	}
	if room != nil && room.Del(ch) {
		// if room empty delete
		room.Del(ch)
	}

	b.cLock.RUnlock()

}

// PushRoom is func to push room msg
func (b *Bucket) PushRoom(c chan *RoomMsgArg) {
	for {
		var (
			arg  *RoomMsgArg
			room *Room
		)
		arg = <-c

		if room = b.Room(arg.RoomID); room != nil {
			log.Bg().Info("开始推送消息")
			room.Push(&arg.P)
		}

	}

}

// Room get a room by rid.
func (b *Bucket) Room(rid string) (room *Room) {
	b.cLock.RLock()
	room, _ = b.rooms[rid]
	b.cLock.RUnlock()
	return
}

// BroadcastRoom is used to broadcast room
func (b *Bucket) BroadcastRoom(arg *RoomMsgArg) {
	// 广播消息递增id
	num := atomic.AddUint64(&b.routinesNum, 1) % b.o.RoutineAmount
	b.routines[num] <- arg

}
