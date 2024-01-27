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

func main() {
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "pong\n")
	})

	http.HandleFunc("/clicked", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "<button>Clicked</button>")
	})

	http.HandleFunc("/colors", func(w http.ResponseWriter, r *http.Request) {
		color := r.URL.Query().Get("color")
		if color == "" || color == "green" {
			color = "red"
		} else if color == "red" {
			color = "blue"
		} else {
			color = "green"
		}
		fmt.Fprintf(w, `<div id="color-demo" class="smooth" style="color:%[1]s" hx-get="/colors?color=%[1]s" hx-swap="outerHTML" hx-trigger="every 1s"> <h1>Hello World :D</h1> </div>`, color)
	})

	http.HandleFunc("/square", func(w http.ResponseWriter, r *http.Request) {
		x := r.URL.Query().Get("x")
		y := r.URL.Query().Get("y")

		if x == "50" && y == "50" {
			y = "150"
		} else if x == "50" && y == "150" {
			x = "150"
		} else if x == "150" && y == "150" {
			y = "50"
		} else { 
			x = "50"
		}
			
		fmt.Fprintf(w, `
			<div id="dot-demo" class="dot smooth" style="top:%[1]spx;left:%[2]spx;"
            hx-get="/square?x=%[1]s&y=%[2]s" hx-swap="outerHTML" hx-trigger="every 2s"></div>
		`, x, y)
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
