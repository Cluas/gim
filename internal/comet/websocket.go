package comet

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Cluas/gim/internal/comet/conf"
	"github.com/Cluas/gim/pkg/log"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

func InitWebsocket(s *Server, c *conf.WebsocketConfig) (err error) {
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
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Error(zap.Error(err))
		return
	}

	ch := NewChannel(s.c.BroadcastSize)
	ch.conn = conn

	go s.writePump(ch)
	go s.readPump(ch)
}

func (s *Server) readPump(ch *Channel) {
	defer func() {
		dArg := new(DisConnArg)

		dArg.RoomID = ch.Room.ID
		if ch.uid != "" {
			dArg.Uid = ch.uid
		}

		s.Bucket(ch.uid).delCh(ch)
		ch.conn.Close()
	}()

	ch.conn.SetReadLimit(s.c.MaxMessageSize)
	_ = ch.conn.SetReadDeadline(time.Now().Add(s.c.PongWait))
	ch.conn.SetPongHandler(func(string) error {
		_ = ch.conn.SetReadDeadline(time.Now().Add(s.c.PongWait))
		return nil
	})

	for {
		_, message, err := ch.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Errorf("readPump ReadMessage err:%v", err)
				return
			}
		}
		if message == nil {
			return
		}
		var (
			connArg *ConnArg
		)

		log.Infof("message :%s", message)
		if err := json.Unmarshal([]byte(message), &connArg); err != nil {
			log.Errorf("message struct %b", connArg)
		}
		connArg.ServerId = conf.Conf.Base.ServerId

		if err != nil {
			log.Errorf("s.operator.Connect error %s", err)
			return
		}

	}
}

func (s *Server) writePump(ch *Channel) {
	ticker := time.NewTicker(s.c.PingPeriod)
	log.Infof("ticker :%v", ticker)

	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		case message, ok := <-ch.broadcast:
			_ = ch.conn.SetWriteDeadline(time.Now().Add(s.c.WriteWait))
			if !ok {
				// The hub closed the channel.
				log.Warn("SetWriteDeadline not ok ")
				_ = ch.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := ch.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Warn(" ch.conn.NextWriter err :%s  ", err)
				return
			}
			log.Infof("message write body:%s", message.Body)
			_, _ = w.Write(message.Body)

			// Add queued chat messages to the current websocket message.
			// n := len(ch.broadcast)
			// for i := 0; i < n; i++ {
			// 	w.Write(newline)
			// 	w.Write(<-ch.broadcast)
			// }

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			_ = ch.conn.SetWriteDeadline(time.Now().Add(s.c.WriteWait))
			log.Infof("websocket.PingMessage :%v", websocket.PingMessage)
			if err := ch.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
