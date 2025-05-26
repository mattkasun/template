// sample program
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/mattkasun/template/templates"

	"github.com/a-h/templ"
)

func favicon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "files/favicon.svg")
}

func styles(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "files/styles.css")
}

func extra(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("extra logging")
		next.ServeHTTP(w, r)
	})
}

func logs(w http.ResponseWriter, r *http.Request) { //nolint:revive,varnamelen
	logs, err := os.ReadFile(os.Args[0] + ".log")
	if err != nil {
		log.Println("get logs", err)
		http.Error(w, "unable to retrieve logs", http.StatusInternalServerError)
		return
	}
	data := []string{}
	lines := strings.Split(string(logs), "\n")
	for i := len(lines) - 1; i > len(lines)-200; i-- {
		if i < 0 {
			break
		}
		data = append(data, lines[i])
	}
	if err := templates.ShowLogs(data).Render(context.Background(), w); err != nil { //nolint:contextcheck
		log.Println("render error", err)
	}
}

func mainPage(w http.ResponseWriter, r *http.Request) { //nolint:revive
	components := []templ.Component{
		templates.Hello(),
	}
	if err := templates.Layout("Template Example", components).Render(context.Background(), w); err != nil { //nolint:contextcheck,lll
		log.Println("render error", err)
	}
}
