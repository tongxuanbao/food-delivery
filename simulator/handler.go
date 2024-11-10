package main

import (
	"fmt"
	"log"
	"net/http"
)

type Simulator struct {
	l *log.Logger
}

func NewSimulator(l *log.Logger) *Simulator {
	return &Simulator{l}
}

func (h *Simulator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.l.Println("Handle Simulator")

	fmt.Fprint(w, "Hello world, it's Simulator service")
}
