package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Position struct {
	x, y int
}

func getPositionMap() map[Position]Position{
	m := make(map[Position]Position)

	m[Position{160, 427}] = Position{200, 370}
	m[Position{200, 370}] = Position{155, 300}
	m[Position{155, 300}] = Position{216, 207}
	m[Position{216, 207}] = Position{270, 270}
	m[Position{270, 270}] = Position{199, 365}
	m[Position{199, 365}] = Position{160, 427}

	return m
}

func getNextPosition(pos Position) Position {
	m := getPositionMap()
	if nextPosition, ok := m[pos]; ok {
		log.Printf("x: %q", nextPosition)
		return nextPosition
	}

	return Position{160, 427}
}

func main() {
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "pong\n")
	})

	http.HandleFunc("/square", func(w http.ResponseWriter, r *http.Request) {
		var x, y int

		if _, err := fmt.Sscanf(r.URL.Query().Get("x"), "%d", &x); err != nil {
			http.Error(w, "Invalid parameter (x)", http.StatusBadRequest)
			return
		}
		if _, err := fmt.Sscanf(r.URL.Query().Get("y"), "%d", &y); err != nil {
			http.Error(w, "Invalid paramete", http.StatusBadRequest)
			return
		}
		
		log.Printf("x: %d, y: %d", x, y)

		nextPosition := getNextPosition(Position{x, y})

		fmt.Fprintf(w, `
			<div id="dot-demo" class="dot smooth" style="top:%[1n]dpx;left:%[2]dpx;"
            hx-get="/square?x=%[1]d&y=%[2]d" hx-swap="outerHTML" hx-trigger="every 2s"></div>
		`, nextPosition.x, nextPosition.y)
	})

	dir := http.Dir("./static")
	fs := http.FileServer(dir)

	http.Handle("/", fs)

	server := &http.Server{Addr: ":8080"}

	go func() {
		log.Println("Starting SIMULATOR server on port 8080")

		err := server.ListenAndServe()
		if err != nil {
			log.Printf("Error starting SIMULATOR server: %s\n", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interrupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	// Block until a signal is received.
	sig := <-c
	log.Printf("Got signal: %s, exiting.", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	server.Shutdown(ctx)
}
