package main

import (
	"fmt"
	"log"
	"net/http"
)

type Order struct {
	l *log.Logger
}

func NewOrder(l *log.Logger) *Order {
	return &Order{l}
}

func (h *Order) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.l.Println("Handle order")

	fmt.Fprint(w, "Hello world, it's ORDER service")
}
