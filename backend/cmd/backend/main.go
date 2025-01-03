package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/tongxuanbao/food-delivery/backend/pkg/geo"
)

func route(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers to allow all origins. You may want to restrict this to specific origins in a production environment.
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Type")

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	list := geo.GetCoordinateList()

	starting := list[0]
	ctx := r.Context()
	for {
		select {
		case <-ctx.Done():
			// The client has disconnected, break the loop
			fmt.Println("Client disconnected, stopping stream.")
			return
		default:
			fmt.Printf("%+v\n", starting)
			// Send out coord
			jsonBytes, err := json.Marshal(starting)
			if err != nil {
				continue
			}
			fmt.Fprintf(w, "data: %s\n\n", fmt.Sprint(string(jsonBytes)))
			time.Sleep(100 * time.Millisecond)
			w.(http.Flusher).Flush()

			// get a random neighbors
			neighbors := starting.GetNeighbors()
			next := neighbors[rand.Intn(len(neighbors))]

			// moving to the random neighbor
			starting = next
		}
	}
}

func main() {
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "pong\n")
	})

	http.HandleFunc("/route", route)

	server := &http.Server{Addr: ":8080"}

	go func() {
		log.Println("Starting backend server on port 8080")

		err := server.ListenAndServe()
		if err != nil {
			log.Printf("Error starting backend server: %s\n", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interrupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	// test := geo.GetAdjacentList()

	// fmt.Printf("%+v\n", test)

	// Block until a signal is received.
	sig := <-c
	log.Printf("Got signal: %s, exiting.", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	server.Shutdown(ctx)
}
