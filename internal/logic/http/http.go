package http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"

	"github.com/Cluas/gim/internal/logic/conf"
	"github.com/Cluas/gim/pkg/log"
	"github.com/Cluas/gim/pkg/net/http/middleware"
	"github.com/Cluas/gim/pkg/net/http/router"
)

// ping test
func ping(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.For(r.Context()).Info("test ping")
	_, _ = fmt.Fprint(w, "pong")
}

// ParseNetwork is func to parse network
func parseNetwork(str string) (network, addr string, err error) {
	if idx := strings.Index(str, "@"); idx == -1 {
		err = fmt.Errorf("addr: \"%s\" error, must be network@tcp:port or network@unixsocket", str)
	} else {
		network = str[:idx]
		addr = str[idx+1:]
	}
	return
}

// Init 初始化 http server
func Init(conf *conf.Config) (err error) {
	//router := tracing.NewRouter()
	r := router.New()
	r.Use(middleware.CORS, middleware.Tracer, middleware.Logger, middleware.Recover)

	r.GlobalOPTIONS = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Access-Control-Request-Method") != "" {
			// Set CORS headers
			header := w.Header()
			header.Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,PATCH,OPTIONS")
			header.Set("Access-Control-Allow-Origin", "*")
			header.Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
		}
		// Adjust status code to 204
		w.WriteHeader(http.StatusNoContent)
	})
	// 中间件
	//mid := NewMiddleware(router)
	// ping test
	r.GET("/ping", ping)
	r.HEAD("/ping", ping) // 阿里云服务健康监控

	_, addr, err := parseNetwork(conf.HTTP.Address[0])
	if err != nil {
		return err
	}
	log.Bg().Info("start server: " + conf.HTTP.Address[0])
	server := &http.Server{
		Addr:              addr,
		Handler:           r,
		ReadTimeout:       conf.HTTP.ReadTimeout,
		ReadHeaderTimeout: conf.HTTP.ReadHeaderTimeout,
		WriteTimeout:      conf.HTTP.WriteTimeout,
		IdleTimeout:       conf.HTTP.IdleTimeout,
	}

	go func() {
		err = server.ListenAndServe()
		log.Bg().Fatal("start server err", zap.Error(err))
	}()
	return nil
}
