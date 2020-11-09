package myroni

import (
	"net/http"
)

type HandlerFunc func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)

func (f HandlerFunc) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	f(rw, r, next)
}

type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request, http.HandlerFunc)
}

type middleware struct {
	handler Handler
	nextfn  http.HandlerFunc // middleware.ServeHTTP
}

func (m *middleware) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	m.handler.ServeHTTP(rw, r, m.nextfn)
}

type Myroni struct {
	middleware middleware
}

func (m *Myroni) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	m.middleware.ServeHTTP(rw, r)
}

func New(handlers ...Handler) *Myroni {
	return &Myroni{
		middleware: build(handlers),
	}
}

func Wrap(handler http.Handler) HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		handler.ServeHTTP(rw, r)
		next(rw, r)
	}
}

func build(handlers []Handler) middleware {
	var next middleware
	switch {
	case len(handlers) == 0:
		return voidMiddleware()
	case len(handlers) > 1:
		next = build(handlers[1:])
	default:
		next = voidMiddleware()
	}
	return newMiddleware(handlers[0], &next)
}

func voidMiddleware() middleware {
	return newMiddleware(
		HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {}),
		&middleware{},
	)
}

func newMiddleware(handler Handler, next *middleware) middleware {
	return middleware{
		handler: handler,
		nextfn:  next.ServeHTTP,
	}
}
