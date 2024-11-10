package main

import (
	"fmt"
	"log"
	"net/http"
)

type Backend struct {
	l *log.Logger
}

func NewBackend(l *log.Logger) *Backend {
	return &Backend{l}
}

func (h *Backend) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.l.Println("Handle Backend")

	fmt.Fprint(w, "Hello world, it's Backend service")
}
