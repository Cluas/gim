package middleware

import (
	"io"
	"net/http"
)

type statusCodeTracker struct {
	http.ResponseWriter
	status      int
	wroteheader bool
}

func (w *statusCodeTracker) WriteHeader(status int) {
	w.status = status
	w.wroteheader = true
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusCodeTracker) Write(b []byte) (int, error) {
	if !w.wroteheader {
		w.wroteheader = true
		w.status = 200
	}
	return w.ResponseWriter.Write(b)
}

// wrappedResponseWriter returns a wrapped version of the original
// ResponseWriter and only implements the same combination of additional
// interfaces as the original.  This implementation is based on
// https://github.com/felixge/httpsnoop.
func (w *statusCodeTracker) wrappedResponseWriter() http.ResponseWriter {
	var (
		hj, i0 = w.ResponseWriter.(http.Hijacker)
		pu, i1 = w.ResponseWriter.(http.Pusher)
		fl, i2 = w.ResponseWriter.(http.Flusher)
		rf, i3 = w.ResponseWriter.(io.ReaderFrom)
	)

	switch {
	case !i0 && !i1 && !i2 && !i3:
		return struct {
			http.ResponseWriter
		}{w}
	case !i0 && !i1 && !i2 && i3:
		return struct {
			http.ResponseWriter
			io.ReaderFrom
		}{w, rf}
	case !i0 && !i1 && i2 && !i3:
		return struct {
			http.ResponseWriter
			http.Flusher
		}{w, fl}
	case !i0 && !i1 && i2 && i3:
		return struct {
			http.ResponseWriter
			http.Flusher
			io.ReaderFrom
		}{w, fl, rf}
	case !i0 && i1 && !i2 && !i3:
		return struct {
			http.ResponseWriter
			http.Pusher
		}{w, pu}
	case !i0 && i1 && !i2 && i3:
		return struct {
			http.ResponseWriter
			http.Pusher
			io.ReaderFrom
		}{w, pu, rf}
	case !i0 && i1 && i2 && !i3:
		return struct {
			http.ResponseWriter
			http.Pusher
			http.Flusher
		}{w, pu, fl}
	case false:
		return struct {
			http.ResponseWriter
			http.Pusher
			http.Flusher
			io.ReaderFrom
		}{w, pu, fl, rf}
	case !i0 && i1 && i2 && i3:
		return struct {
			http.ResponseWriter
			http.Pusher
			http.Flusher
			io.ReaderFrom
		}{w, pu, fl, rf}
	case i0 && !i1 && !i2 && !i3:
		return struct {
			http.ResponseWriter
			http.Hijacker
		}{w, hj}
	case i0 && !i1 && !i2 && i3:
		return struct {
			http.ResponseWriter
			http.Hijacker
			io.ReaderFrom
		}{w, hj, rf}
	case i0 && !i1 && i2 && !i3:
		return struct {
			http.ResponseWriter
			http.Hijacker
			http.Flusher
		}{w, hj, fl}
	case i0 && !i1 && i2 && i3:
		return struct {
			http.ResponseWriter
			http.Hijacker
			http.Flusher
			io.ReaderFrom
		}{w, hj, fl, rf}
	case i0 && i1 && !i2 && !i3:
		return struct {
			http.ResponseWriter
			http.Hijacker
			http.Pusher
		}{w, hj, pu}
	case i0 && i1 && !i2 && i3:
		return struct {
			http.ResponseWriter
			http.Hijacker
			http.Pusher
			io.ReaderFrom
		}{w, hj, pu, rf}
	case i0 && i1 && i2 && !i3:
		return struct {
			http.ResponseWriter
			http.Hijacker
			http.Pusher
			http.Flusher
		}{w, hj, pu, fl}
	case i0 && i1 && i2 && i3:
		return struct {
			http.ResponseWriter
			http.Hijacker
			http.Pusher
			http.Flusher
			io.ReaderFrom
		}{w, hj, pu, fl, rf}
	default:
		return struct {
			http.ResponseWriter
		}{w}
	}
}
