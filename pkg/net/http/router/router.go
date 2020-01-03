package router

import (
	"github.com/julienschmidt/httprouter"
)

type middleware func(httprouter.Handle) httprouter.Handle

//Router is have middleware httprouter
type Router struct {
	middlewareChain []middleware
	*httprouter.Router
}

//New return a new Router
func New() *Router {
	return &Router{
		Router: httprouter.New(),
	}
}

//Use is func to add middleware
func (r *Router) Use(mm ...middleware) {
	for _, m := range mm {
		r.middlewareChain = append(r.middlewareChain, m)
	}

}

//GET is func with middleware
func (r *Router) GET(route string, h httprouter.Handle) {
	r.Router.GET(route, r.wrap(h))
}

//POST is func with middleware
func (r *Router) POST(route string, h httprouter.Handle) {
	r.Router.POST(route, r.wrap(h))
}

//PUT is func with middleware
func (r *Router) PUT(route string, h httprouter.Handle) {
	r.Router.PUT(route, r.wrap(h))
}

//PATCH is func with middleware
func (r *Router) PATCH(route string, h httprouter.Handle) {
	r.Router.PATCH(route, r.wrap(h))
}

//DELETE is func with middleware
func (r *Router) DELETE(route string, h httprouter.Handle) {
	r.Router.DELETE(route, r.wrap(h))
}

//HEAD is func with middleware
func (r *Router) HEAD(route string, h httprouter.Handle) {
	r.Router.HEAD(route, r.wrap(h))
}

//OPTIONS is func with middleware
func (r *Router) OPTIONS(route string, h httprouter.Handle) {
	r.Router.OPTIONS(route, r.wrap(h))
}

func (r *Router) wrap(h httprouter.Handle) httprouter.Handle {
	for i := len(r.middlewareChain) - 1; i >= 0; i-- {
		h = r.middlewareChain[i](h)
	}
	return h
}
