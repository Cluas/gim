package comet

import (
	"github.com/Cluas/gim/internal/comet/conf"
	"github.com/Cluas/gim/internal/comet/rpc"
	"github.com/Cluas/gim/pkg/cityhash"
	"github.com/Cluas/gim/pkg/log"
)

type Server struct {
	buckets   []*Bucket // subKey bucket
	c         *conf.Config
	bucketIdx uint32
	operator  rpc.Operator
}

// NewServer returns a new Server.
func NewServer(c *conf.Config) *Server {
	s := new(Server)
	s.operator = new(rpc.DefaultOperator)
	s.buckets = make([]*Bucket, c.Bucket.Size)
	s.bucketIdx = uint32(len(s.buckets))
	s.c = c
	for i := 0; i < conf.Conf.Bucket.Size; i++ {
		s.buckets[i] = NewBucket(&BucketOptions{
			ChannelSize: conf.Conf.Bucket.Channel,
			RoomSize:    conf.Conf.Bucket.Room,
		})
	}
	return s
}

func (s *Server) Bucket(subKey string) *Bucket {
	idx := cityhash.CityHash32([]byte(subKey), uint32(len(subKey))) % s.bucketIdx
	log.Infof("\"%s\" hit channel bucket index: %d use cityhash", subKey, idx)
	return s.buckets[idx]
}
