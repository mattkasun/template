package main

import (
	"context"
	"errors"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/devilcove/configuration"
	"github.com/devilcove/cookie"
)

type Config struct {
	Port int
}

const (
	cookieName = "templateCookie"
	cookieAge  = 300
)

var templates *template.Template //nolint:gochecknoglobals

func web(ctx context.Context) {
	if err := cookie.New(cookieName, cookieAge); err != nil {
		slog.Error("new cookie", "error", err)
		return
	}
	var config Config
	// create config dir if it doesn't exist
	xdg, _ := os.UserConfigDir()
	if err := os.MkdirAll(filepath.Join(xdg, filepath.Base(os.Args[0])), 0o700); err != nil {
		slog.Error("create config dir", "error", err)
		return
	}
	// create config file if it doesn't exit
	f, err := os.OpenFile(filepath.Join(xdg, filepath.Base(os.Args[0]), "config"), os.O_RDONLY|os.O_CREATE, 0o600)
	if err != nil {
		slog.Error("read config file", "error", err)
		return
	}
	f.Close()
	if err := configuration.Get(&config); err != nil {
		slog.Error("get configuration", "error", err)
		return
	}
	if config.Port == 0 {
		config.Port = 8080
	}
	log.Println("starting web server")
	server := http.Server{
		Addr:              "127.0.0.1:" + strconv.Itoa(config.Port),
		Handler:           setupRouter(),
		ReadHeaderTimeout: time.Second,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("web server start", "error", err)
			return
		}
	}()

	slog.Info("web server started", "port", config.Port)
	<-ctx.Done()
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("server shutdown", "error", err)
		return
	}
	slog.Info("web server shutdown")
}

func setupRouter() http.Handler {
	templates = template.Must(template.ParseGlob("templates/*html"))
	mux := http.NewServeMux()
	mux.Handle("/", notFound(mux))
	mux.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("assets"))))
	mux.HandleFunc("GET /{$}", displayMain)
	// mux.HandleFunc("GET /login", displayLogin)
	// mux.HandleFunc("POST /login", processLogin)
	// mux.HandleFunc("/logout", logout)

	group := http.NewServeMux()
	group.Handle("/", notFound(group))
	// group.HandleFunc("GET /{$}", groupMain)

	groupHandler := auth(group)
	mux.Handle("/group/", http.StripPrefix("/group", groupHandler))
	handler := logger(mux)
	return handler
}
