package comet

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"

	"github.com/Cluas/gim/internal/comet/conf"
	"github.com/Cluas/gim/pkg/log"
)

// InitWebsocket is func to initial Websocket
func InitWebsocket(s *Server, c *conf.WebsocketConf) (err error) {
	mux := s.createServeMux()
	go func() {
		if err = http.ListenAndServe(c.Bind, mux); err != nil {
			log.Bg().Panic("启动失败:", zap.Error(err))

		}
	}()

	return err

}

// serveWS handles websocket requests from the peer.
func (s *Server) serveWS(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	herder := http.Header{}
	wsProto := strings.Split(r.Header.Get("Sec-WebSocket-Protocol"), ",")
	if len(wsProto) < 2 {
		return
	}
	token := wsProto[0]
	roomID := strings.TrimSpace(wsProto[1])
	args := &ConnectArg{
		Auth:     token,
		RoomID:   roomID,
		ServerID: conf.Conf.Base.ServerID,
	}
	uid, err := s.operator.Connect(ctx, args)

	if err != nil {
		log.Bg().Error("调用logic Connect 方法失败", zap.Error(err))
	}

	herder.Add("Sec-WebSocket-Protocol", roomID)

	upgrades := websocket.Upgrader{
		ReadBufferSize:    s.c.Websocket.ReadBufferSize,
		WriteBufferSize:   s.c.Websocket.WriteBufferSize,
		EnableCompression: true,
	}
	// CORS
	upgrades.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrades.Upgrade(w, r, herder)

	if err != nil {
		log.Bg().Error("", zap.Error(err))
		return
	}
	if uid == "" {
		_ = conn.WriteJSON(map[string]string{"code": "401", "msg": "token error!"})
		_ = conn.Close()
		return
	}

	ch := NewChannel(s.c.Bucket.BroadcastSize)
	ch.conn = conn

	b := s.Bucket(ctx, uid)
	err = b.Put(uid, roomID, ch)
	if err != nil {
		log.Bg().Error("bucket Put err: ", zap.Error(err))
		_ = ch.conn.Close()
	}

	go s.writePump(ctx, ch)
	go s.readPump(ctx, ch)
}

func (s *Server) readPump(ctx context.Context, ch *Channel) {
	defer func() {
		if ch.uid != "" {
			s.Bucket(ctx, ch.uid).delCh(ch)
			args := new(DisconnectArg)
			args.UID = ch.uid
			args.RoomID = ch.Room.ID
			if err := s.operator.Disconnect(ctx, args); err != nil {
				log.Bg().Error("Disconnect err :%s", zap.Error(err))
			}
		}
		_ = ch.conn.Close()
	}()

	ch.conn.SetReadLimit(s.c.Websocket.MaxMessageSize)
	_ = ch.conn.SetReadDeadline(time.Now().Add(s.c.Websocket.PongWait))
	ch.conn.SetPongHandler(func(string) error {
		_ = ch.conn.SetReadDeadline(time.Now().Add(s.c.Websocket.PongWait))
		return nil
	})

	for {
		_, message, err := ch.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Bg().Error("readPump ReadMessage err: ", zap.Error(err))
			}
		}
		if message == nil {
			return
		}

	}
}

func (s *Server) writePump(ctx context.Context, ch *Channel) {
	ticker := time.NewTicker(s.c.Websocket.PingPeriod)

	defer func() {
		ticker.Stop()
		_ = ch.conn.Close()
	}()
	for {
		select {
		case message, ok := <-ch.broadcast:
			_ = ch.conn.SetWriteDeadline(time.Now().Add(s.c.Websocket.WriteWait))
			if !ok {
				// The hub closed the channel.
				log.Bg().Warn("SetWriteDeadline is not ok ")

				_ = ch.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := ch.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Bg().Info("ch.conn.NextWriter err: ", zap.Error(err))
				return
			}
			log.Bg().Info("received message: ", zap.Binary("message", message.Body))
			_, _ = w.Write(message.Body)

			if e := w.Close(); e != nil {
				return
			}
		// Heartbeat
		case <-ticker.C:
			_ = ch.conn.SetWriteDeadline(time.Now().Add(s.c.Websocket.WriteWait))
			if err := ch.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Bg().Info("use close connection", zap.Error(err))
				return
			}
		}
	}
}
