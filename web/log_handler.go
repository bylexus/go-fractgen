package web

import (
	"fmt"
	"net/http"
)

type LogHandler struct {
	internalHandler http.Handler
}

func NewLogHandler(handler http.Handler) *LogHandler {
	return &LogHandler{internalHandler: handler}
}

func (l *LogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	url := r.URL.String()
	fmt.Printf("%s %s\n", method, url)
	l.internalHandler.ServeHTTP(w, r)
}
