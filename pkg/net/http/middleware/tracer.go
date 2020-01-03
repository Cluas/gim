package middleware

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
)

//RequestIDContextKey is context key
type RequestIDContextKey struct{}

//Tracer is middleware for tracing HTTP Server
func Tracer(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

		if !opentracing.IsGlobalTracerRegistered() {
			fn(w, r, p)
			return
		}
		tracer := opentracing.GlobalTracer()
		ctx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		sp := tracer.StartSpan("HTTP "+r.Method+" "+r.URL.Path, ext.RPCServerOption(ctx))
		ext.HTTPMethod.Set(sp, r.Method)
		ext.HTTPUrl.Set(sp, r.URL.String())
		componentName := "httprouter"
		ext.Component.Set(sp, componentName)

		sct := &statusCodeTracker{ResponseWriter: w}
		rc := opentracing.ContextWithSpan(r.Context(), sp)
		// Add requestID
		if sc, ok := sp.Context().(jaeger.SpanContext); ok {
			rc = context.WithValue(rc, RequestIDContextKey{}, sc.TraceID().String())
		}
		fn(sct.wrappedResponseWriter(), r.WithContext(rc), p)
		defer func() {
			ext.HTTPStatusCode.Set(sp, uint16(sct.status))
			if sct.status >= http.StatusInternalServerError || !sct.wroteheader {
				ext.Error.Set(sp, true)
			}
			sp.Finish()
		}()

	}
}
