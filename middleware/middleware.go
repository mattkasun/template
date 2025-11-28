// Package middleware provides middleware framework and logging middleware
package middleware

import (
	"log"
	"net/http"
)

// Middleware provides a chain of middlewares.
type Middleware struct {
	chain http.Handler
}

// Default creates a default middleware using DefaultServeMux.
func Default() *Middleware {
	return &Middleware{http.DefaultServeMux}
}

// New returns a new middleware.
func New(handler http.Handler) *Middleware {
	return &Middleware{handler}
}

// Router creates a severMux for given prefix and middleware chain.
func Router(prefix string, middleware ...func(http.Handler) http.Handler) *http.ServeMux {
	mux := http.NewServeMux()
	middle := New(mux)
	middle.UseGroup(middleware...)
	http.Handle(prefix+"/", http.StripPrefix(prefix, middle))
	return mux
}

// Use adds middleware to the chain.
func (a *Middleware) Use(handler func(http.Handler) http.Handler) {
	a.chain = handler(a.chain)
}

// UseGroup adds a chain of middlewares.
func (a *Middleware) UseGroup(handlers ...func(http.Handler) http.Handler) {
	for _, handler := range handlers {
		a.Use(handler)
	}
}

// ServeHTTP implements http.Handler interface.
func (a *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.chain.ServeHTTP(w, r)
}

// Logger is a logging middleware that logs useragent, RemoteAddr, Method, Host, Path and response.Status to stdlib log.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := statusRecorder{w, http.StatusOK}
		next.ServeHTTP(&rec, r)
		log.Println(r.UserAgent(), r.RemoteAddr, r.Method, r.Host, r.URL.Path, rec.status)
	})
}

type statusRecorder struct {
	http.ResponseWriter

	status int
}

// WriteHeader overrides std WriteHeader func to save response code.
func (rec *statusRecorder) WriteHeader(code int) {
	rec.status = code
	rec.ResponseWriter.WriteHeader(code)
}
