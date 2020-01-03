package perf

import (
	"net/http"
	"net/http/pprof"

	"go.uber.org/zap"

	"github.com/Cluas/gim/pkg/log"
)

// Init start http pprof.
func Init(pBind []string) {
	pprofServeMux := http.NewServeMux()
	pprofServeMux.HandleFunc("/debug/pprof/", pprof.Index)
	pprofServeMux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	pprofServeMux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	pprofServeMux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	for _, addr := range pBind {
		go func(addr string) {
			if err := http.ListenAndServe(addr, pprofServeMux); err != nil {
				log.Bg().Info("http.ListenAndServe err", zap.String("addr", addr), zap.Error(err))
			}
		}(addr)
	}
}
