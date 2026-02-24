package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/devilcove/cookie"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

// WriteHeader overrides std WriteHeader to save response code.
func (rec *statusRecorder) WriteHeader(code int) {
	rec.status = code
	rec.ResponseWriter.WriteHeader(code)
}

// logger is a logging middleware that logs useragent, RemoteAddr, Method, Host, Path and response.Status to stdlib log.
func logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		rec := statusRecorder{w, http.StatusOK}
		next.ServeHTTP(&rec, r)
		// remote := strings.Split(r.RemoteAddr, ":")[0]
		remote := r.RemoteAddr
		if r.Header.Get("X-Forwarded-For") != "" {
			remote = r.Header.Get("X-Forwarded-For")
		}
		details := fmt.Sprintf(
			"%s %s%s %d %s %s %s",
			r.Method,
			r.Host,
			r.URL.Path,
			rec.status,
			remote,
			time.Since(now).String(),
			r.UserAgent(),
		)
		log.Println(details)
	})
}

func methods() []string {
	return []string{
		http.MethodGet,
		http.MethodHead,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodConnect,
		http.MethodOptions,
		http.MethodTrace,
	}
}

func notFound(mux *http.ServeMux) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var allowed []string
		method := r.Method
		_, current := mux.Handler(r)
		for _, method := range methods() {
			r.Method = method
			if _, pattern := mux.Handler(r); pattern != current {
				allowed = append(allowed, method)
			}
		}
		r.Method = method
		if len(allowed) != 0 {
			w.Header().Set("Allow", strings.Join(allowed, ", "))
			w.WriteHeader(http.StatusMethodNotAllowed)
			io.WriteString(w, "this is not the method you are looking for...")
			return
		}
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w,
			"this is not the page you are loooking for...\n\ngo about your business\nmove along")
	})
}

func auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := cookie.Get(r, cookieName); err != nil {
			handleError(w, http.StatusUnauthorized, "<a href='/'>unauthorized</a>")
			return
		}
		next.ServeHTTP(w, r)
	})
}
