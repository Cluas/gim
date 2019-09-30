package comet

import (
	"sync"
)

// BucketConfig is struct of Bucket Config
type BucketConfig struct {
	ChannelSize int
	RoomSize    int
}

// Bucket is struct of Bucket
type Bucket struct {
	cLock       sync.RWMutex        // protect the channels for chs
	chs         map[string]*Channel // map sub key to a channel
	bConfig     *BucketConfig
	rooms       map[int32]*Room // bucket room channel
	routinesNum uint64
	broadcast   chan []byte
}

// NewBucket is constructor of Bucket
func NewBucket(bConfig *BucketConfig) *Bucket {
	return &Bucket{
		chs: make(map[string]*Channel, bConfig.ChannelSize),
	}
}

// Put is func to add channel
func (b *Bucket) Put(key string, rid int32, ch *Channel) error {
	var (
		room *Room
		ok   bool
	)
	b.cLock.Lock()

	if rid != NOROOM {
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
