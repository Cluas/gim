package comet

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Cluas/gim/internal/comet/conf"
	"github.com/Cluas/gim/internal/comet/rpc"
	"github.com/Cluas/gim/pkg/log"
	"github.com/gorilla/websocket"
)

func InitWebsocket(s *Server, c *conf.WebsocketConf) (err error) {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(s, w, r)
	})
	err = http.ListenAndServe(c.Bind, nil)

	return err

}

// serveWs handles websocket requests from the peer.
func serveWs(s *Server, w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  s.c.ReadBufferSize,
		WriteBufferSize: s.c.WriteBufferSize,
	}
	// CORS
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Error(err)
		return
	}

	ch := NewChannel(s.c.BroadcastSize)
	ch.conn = conn

	go s.writePump(ch)
	go s.readPump(ch)
}

func (s *Server) readPump(ch *Channel) {
	defer func() {
		ch.conn.Close()
	}()

	ch.conn.SetReadLimit(s.c.MaxMessageSize)
	_ = ch.conn.SetReadDeadline(time.Now().Add(s.c.PongWait))
	ch.conn.SetPongHandler(func(string) error {
		_ = ch.conn.SetReadDeadline(time.Now().Add(s.c.PongWait))
		log.Infof("websocket pong...")
		return nil
	})

	for {
		_, message, err := ch.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Errorf("readPump ReadMessage err:%v", err)
			}
		}
		if message == nil {
			return
		}
		var (
			connArg *rpc.ConnectArg
		)

		log.Infof("message :%s", message)
		if err := json.Unmarshal([]byte(message), &connArg); err != nil {
			log.Errorf("message struct %b", connArg)
			connArg = &rpc.ConnectArg{
				Auth:     "123",
				RoomID:   100,
				ServerID: 100,
			}
		}
		uid, err := s.operator.Connect(connArg)
		log.Infof("websocket uid:%s", uid)

		if err != nil {
			log.Errorf("s.operator.Connect error %s", err)
		}

		b := s.Bucket(uid)
		// TODO rpc 操作获取uid 存入ch 存入Server

		// b.broadcast <- message
		err = b.Put(uid, connArg.RoomID, ch)
		if err != nil {
			log.Errorf("conn close err: %s", err)
			ch.conn.Close()
		}
		log.Infof("message  333 :%s", message)
		ch.broadcast <- message

	}
}

func (s *Server) writePump(ch *Channel) {
	ticker := time.NewTicker(s.c.PingPeriod)
	log.Infof("ticker :%v", ticker)

	defer func() {
		ticker.Stop()
		ch.conn.Close()
	}()
	for {
		select {
		case message, ok := <-ch.broadcast:
			_ = ch.conn.SetWriteDeadline(time.Now().Add(s.c.WriteWait))
			if !ok {
				// The hub closed the channel.
				log.Warn("SetWriteDeadline is not ok ")
				_ = ch.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			log.Infof("TextMessage :%v", websocket.TextMessage)
			w, err := ch.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Warn(" ch.conn.NextWriter err :%s  ", err)
				return
			}
			log.Infof("message: %v", message)
			_, _ = w.Write(message)

			// Add queued chat messages to the current websocket message.
			l := len(ch.broadcast)
			for i := 0; i < l; i++ {
				_, _ = w.Write([]byte{'\n'})
				_, _ = w.Write(<-ch.broadcast)
			}

			if err := w.Close(); err != nil {
				return
			}
		// Heartbeat
		case <-ticker.C:
			_ = ch.conn.SetWriteDeadline(time.Now().Add(s.c.WriteWait))
			log.Infof("websocket ping... :%v", websocket.PingMessage)
			if err := ch.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Error(err)
				return
			}
		}
	}
}
