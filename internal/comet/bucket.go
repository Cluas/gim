package comet

import (
	"sync"
	"sync/atomic"
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
	o           *BucketOptions
	rooms       map[int32]*Room // bucket room channel
	routines    []chan *RoomMsgArg
	routinesNum uint64
	broadcast   chan []byte
}

// NewBucket is constructor of Bucket
func NewBucket(o *BucketOptions) *Bucket {
	return &Bucket{
		chs:   make(map[string]*Channel, o.ChannelSize),
		o:     o,
		rooms: make(map[int32]*Room, o.RoomSize),
	}
}

// Put is func to add channel
func (b *Bucket) Put(key string, rid int32, ch *Channel) error {
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
	b.chs[key] = ch
	b.cLock.Unlock()

	if room != nil {
		err := room.Put(ch)
		return err
	}
	return nil
}
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
func (b *Bucket) PushRoom(c chan *RoomMsgArg) {
	for {
		var (
			arg  *RoomMsgArg
			room *Room
		)
		arg = <-c

		if room = b.Room(arg.RoomID); room != nil {
			room.Push(&arg.P)
		}

	}

}

// Room get a room by roomid.
func (b *Bucket) Room(rid int32) (room *Room) {
	b.cLock.RLock()
	room, _ = b.rooms[rid]
	b.cLock.RUnlock()
	return
}

func (b *Bucket) BroadcastRoom(arg *RoomMsgArg) {
	// 广播消息递增id
	num := atomic.AddUint64(&b.routinesNum, 1) % b.o.RoutineAmount
	// log.Infof("BroadcastRoom RoomMsgArg :%s", arg)
	// log.Infof("bucket routinesNum :%d", b.routinesNum)
	b.routines[num] <- arg

}
