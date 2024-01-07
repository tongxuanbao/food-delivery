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

	http.HandleFunc("/new", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello world, it's order1 service\n")
		go func() {
			log.Println("Order placed")
			time.Sleep(5 * time.Second)
			log.Println("Order accepted")
			time.Sleep(5 * time.Second)
			log.Println("Order cooking")
			time.Sleep(5 * time.Second)
			log.Println("Order cooked and waiting for shipper")
			time.Sleep(5 * time.Second)
			log.Println("Order on the way")
			time.Sleep(5 * time.Second)
			log.Println("Order delivered")
		}()
	})

	server := &http.Server{Addr: ":8080"}

	go func() {
		log.Println("Starting ORDER server on port 8080")

		err := server.ListenAndServe()
		if err != nil {
			log.Printf("Error starting ORDER server: %s\n", err)
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
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	server.Shutdown(ctx)
}
