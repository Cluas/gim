package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

type requestTimeKey struct{}

//CORS is middleware to CORS
func CORS(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		now := time.Now()
		ctx := context.WithValue(r.Context(), requestTimeKey{}, now)
		r = r.WithContext(ctx)
		fn(w, r, p)
	}
}
