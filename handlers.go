// sample program
package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"path/filepath"
	"runtime"
)

func displayMain(w http.ResponseWriter, _ *http.Request) {
	render(w, "welcome", nil)
}

func render(w io.Writer, template string, data any) {
	if err := templates.ExecuteTemplate(w, template, data); err != nil {
		slog.Error("render template", "caller", caller(2), "name", template,
			"data", data, "error", err)
	}
}

func caller(depth int) string {
	pc, file, no, ok := runtime.Caller(depth)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		return fmt.Sprintf("%s %s:%d", details.Name(), filepath.Base(file), no)
	}
	return "unknown caller"
}

func handleError(w http.ResponseWriter, status int, message string) {
	buf := bytes.Buffer{}
	l := log.New(&buf, "ERROR: ", log.Lshortfile)
	_ = l.Output(2, message)
	slog.Error(buf.String())
	w.WriteHeader(status)
	render(w, "error", "unauthorized")
}
