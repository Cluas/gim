package comet

import (
	"net/http"
	"strings"
	"time"

	"github.com/Cluas/gim/internal/comet/conf"
	"github.com/Cluas/gim/pkg/log"
	"github.com/gorilla/websocket"
)

// InitWebsocket is func to initial Websocket
func InitWebsocket(s *Server, c *conf.WebsocketConf) (err error) {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(s, w, r)
	})
	err = http.ListenAndServe(c.Bind, nil)

	return err

}

// serveWs handles websocket requests from the peer.
func serveWs(s *Server, w http.ResponseWriter, r *http.Request) {
	herder := http.Header{}
	wsProto := strings.Split(r.Header.Get("Sec-WebSocket-Protocol"), ",")
	if len(wsProto) < 2 {
		return
	}
	token := wsProto[0]
	roomID := wsProto[1]
	args := &ConnectArg{
		Auth:     token,
		RoomID:   roomID,
		ServerID: conf.Conf.Base.ServerID,
	}
	uid, err := s.operator.Connect(args)

	if err != nil {
		log.Errorf("s.operator.Connect error %s", err)
	}

	herder.Add("Sec-WebSocket-Protocol", roomID)

	upgrades := websocket.Upgrader{
		ReadBufferSize:    s.c.ReadBufferSize,
		WriteBufferSize:   s.c.WriteBufferSize,
		EnableCompression: true,
	}
	// CORS
	upgrades.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrades.Upgrade(w, r, herder)

	if err != nil {
		log.Error(err)
		return
	}
	if uid == "" {
		_ = conn.WriteJSON(map[string]string{"code": "401", "msg": "token error!"})
		_ = conn.Close()
		return
	}

	ch := NewChannel(s.c.BroadcastSize)
	ch.conn = conn

	b := s.Bucket(uid)

	err = b.Put(uid, roomID, ch)
	if err != nil {
		log.Errorf("conn close err: %s", err)
		_ = ch.conn.Close()
	}

	go s.writePump(ch)
	go s.readPump(ch)
}

func (s *Server) readPump(ch *Channel) {
	defer func() {
		if ch.uid != "" {
			s.Bucket(ch.uid).delCh(ch)
		}
		_ = ch.conn.Close()
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
			}
		}
		if message == nil {
			return
		}

	}
}

func (s *Server) writePump(ch *Channel) {
	ticker := time.NewTicker(s.c.PingPeriod)
	log.Infof("create ticker...")

	defer func() {
		ticker.Stop()
		_ = ch.conn.Close()
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
			_, _ = w.Write(message.Body)

			if err := w.Close(); err != nil {
				return
			}
		// Heartbeat
		case <-ticker.C:
			_ = ch.conn.SetWriteDeadline(time.Now().Add(s.c.WriteWait))
			if err := ch.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Error(err)
				return
			}
		}
	}
}
