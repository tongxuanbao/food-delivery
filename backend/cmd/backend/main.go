package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/tongxuanbao/food-delivery/backend/internal/app/restaurant"
	"github.com/tongxuanbao/food-delivery/backend/internal/app/simulator"
	"github.com/tongxuanbao/food-delivery/backend/internal/pkg/geo"
)

type RateResponse struct {
	Id       int              `json:"id"`
	Route    []geo.Coordinate `json:"route"`
	Position geo.Coordinate   `json:"position"`
}

func route(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers to allow all origins. You may want to restrict this to specific origins in a production environment.
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Type")

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	coordList := geo.GetCoordinateList()

	list := geo.GetRoute(coordList[0], coordList[1000])
	jsonBytes, err := json.Marshal(restaurant.RestaurantList)
	if err == nil {
		fmt.Fprintf(w, "event: restaurant\ndata: %s\n\n", fmt.Sprint(string(jsonBytes)))
	}
	w.(http.Flusher).Flush()

	ctx := r.Context()
	for idx, coord := range list {
		select {
		case <-ctx.Done():
			// The client has disconnected, break the loop
			fmt.Println("Client disconnected, stopping stream.")
			return
		default:
			// Response
			response := RateResponse{Id: 1, Route: list[idx:], Position: coord}
			// Send out coord
			jsonBytes, err := json.Marshal(response)
			if err != nil {
				continue
			}
			fmt.Fprintf(w, "event: test\ndata: %s\n\n", fmt.Sprint(string(jsonBytes)))
			time.Sleep(1000 * time.Millisecond)
			w.(http.Flusher).Flush()
		}
	}
}

func simulate(restaurantService *restaurant.Service) {
	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				restaurantService.AddOrder(1, 1)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func handleRoute(restaurantService *restaurant.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers to allow all origins. You may want to restrict this to specific origins in a production environment.
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Type")

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		// Send initial data
		fmt.Fprintf(w, "event: initial\ndata: {\"name\": \"Pasta Place\", \"rating\": 4.5}\n\n")
		w.(http.Flusher).Flush()

		// Subscribe to services
		restaurantSubscriber := restaurantService.Broker.Subscribe("restaurant")

		go func() {
			for {
				select {
				case msg, ok := <-restaurantSubscriber.Channel:
					if !ok {
						fmt.Println("Subscriber channel closed.")
						return
					}
					// Send data
					fmt.Fprintf(w, "event: restaurant\ndata: %s\n\n", msg)
					w.(http.Flusher).Flush()
				case <-restaurantSubscriber.Unsubscribe:
					fmt.Println("Unsubscribed.")
					return
				}
			}
		}()

		ctx := r.Context()
		<-ctx.Done()
	}
}

func NewServer(restaurantService *restaurant.Service) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/route", handleRoute(restaurantService))
	return mux
}

func main() {
	// Initialize services
	restaurantService := restaurant.New()
	serverHandler := NewServer(restaurantService)

	server := &http.Server{
		Addr:    ":8080",
		Handler: serverHandler,
	}
	go func() {
		log.Println("Starting backend server on port 8080")

		err := server.ListenAndServe()
		if err != nil {
			log.Printf("Error starting backend server: %s\n", err)
			os.Exit(1)
		}
	}()

	simulatorService := simulator.New(restaurantService)
	go simulatorService.Simulate()

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
