package web

import (
	"net/http"
)

type CorsHandler struct {
	internalHandler http.Handler
}

func NewCorsHandler(handler http.Handler) *CorsHandler {
	return &CorsHandler{internalHandler: handler}
}

func (l *CorsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	l.internalHandler.ServeHTTP(w, r)
}
