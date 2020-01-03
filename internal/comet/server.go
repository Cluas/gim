package comet

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Cluas/gim/internal/comet/conf"
	"github.com/Cluas/gim/pkg/cityhash"
	"github.com/Cluas/gim/pkg/log"
)

// CurrentServer is var of current server
var CurrentServer *Server

// Server is struct for Comet Server
type Server struct {
	Buckets   []*Bucket // subKey bucket
	c         *conf.Config
	bucketIdx uint32
	operator  Operator
}

// NewServer returns a new Server.
func NewServer(c *conf.Config) *Server {
	s := new(Server)
	s.operator = new(DefaultOperator)
	s.Buckets = make([]*Bucket, c.Bucket.Size)
	s.bucketIdx = uint32(len(s.Buckets))
	s.c = c
	for i := 0; i < conf.Conf.Bucket.Size; i++ {
		s.Buckets[i] = NewBucket(&BucketOptions{
			ChannelSize:   conf.Conf.Bucket.Channel,
			RoomSize:      conf.Conf.Bucket.Room,
			RoutineAmount: conf.Conf.Bucket.RoutineAmount,
			RoutineSize:   conf.Conf.Bucket.RoutineSize,
		})
	}
	CurrentServer = s
	return s
}

// Bucket is func to location bucket use cityhash
func (s *Server) Bucket(ctx context.Context, subKey string) *Bucket {
	idx := cityhash.CityHash32([]byte(subKey), uint32(len(subKey))) % s.bucketIdx
	log.Bg().Info(fmt.Sprintf("\"%s\" hit channel bucket index: %d use cityhash", subKey, idx))
	return s.Buckets[idx]
}

func (s *Server) createServeMux() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/ws", http.HandlerFunc(s.serveWS))
	return mux
}
