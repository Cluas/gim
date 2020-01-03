package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"

	"github.com/Cluas/gim/pkg/log"
)

const (
	//RequestIDHeaderKey is header key
	RequestIDHeaderKey = "X-Request-ID"
)

//Logger is middleware to record std log
func Logger(fn httprouter.Handle) httprouter.Handle {

	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		sct := &statusCodeTracker{ResponseWriter: w}
		fn(sct.wrappedResponseWriter(), r, p)
		logging(sct, r)
	}
}

func logging(sct *statusCodeTracker, r *http.Request) {
	duration := time.Now().Sub(r.Context().Value(requestTimeKey{}).(time.Time))
	if sct.status >= http.StatusOK && sct.status < 400 {
		log.Bg().Info(
			"",
			zap.String("request_ip", getRequestIP(r)),
			zap.String("request_id", getRequestID(r)),
			zap.String("status_code", strconv.Itoa(sct.status)),
			zap.String("duration", duration.String()),
			zap.String("method", r.Method),
			zap.String("uri", r.RequestURI),
			zap.String("agent", r.UserAgent()),
		)
	} else {
		log.Bg().Error("",
			zap.String("request_ip", getRequestIP(r)),
			zap.String("request_id", getRequestID(r)),
			zap.String("status_code", strconv.Itoa(sct.status)),
			zap.String("duration", duration.String()),
			zap.String("method", r.Method),
			zap.String("uri", r.RequestURI),
			zap.String("agent", r.UserAgent()),
		)
	}
}

func getRequestID(r *http.Request) string {
	reqID, ok := r.Context().Value(RequestIDContextKey{}).(string)
	if ok {
		return reqID
	}
	reqID = r.Header.Get(RequestIDHeaderKey)
	if reqID == "" {
		reqID = uuid.New().String()
	}
	return reqID
}

func getRequestIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}
