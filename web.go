package main

import (
	"context"
	"crypto/rand"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/mattkasun/template/middleware"

	"github.com/gorilla/sessions"
)

const (
	sessionBytes = 32
	cookieAge    = 300
)

var store *sessions.CookieStore //nolint:gochecknoglobals

func web(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	store = sessions.NewCookieStore(randBytes(sessionBytes))
	store.MaxAge(cookieAge)
	store.Options.HttpOnly = true
	store.Options.SameSite = http.SameSiteStrictMode
	log.Println("starting web server")

	logger := middleware.New(http.DefaultServeMux)
	logger.Use(middleware.Logger)
	http.HandleFunc("/favicon.ico", favicon)
	http.HandleFunc("/styles.css", styles)

	plain := middleware.Router("", extra)
	plain.HandleFunc("/{$}", mainPage)
	plain.HandleFunc("/logs", logs)

	server := http.Server{Addr: ":8090", ReadHeaderTimeout: time.Second, Handler: logger} //nolint:exhaustruct
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Println("web server", err)
		}
	}()
	log.Println("web server running :8090")
	<-ctx.Done()
	log.Println("shutdown web ...")
	if err := server.Shutdown(context.Background()); err != nil { //nolint:contextcheck
		log.Println("web server shutdown", err)
	}
}

func randBytes(l int) []byte {
	bytes := make([]byte, l)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}
	return bytes
}
