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

	"github.com/google/uuid"
	"github.com/tongxuanbao/food-delivery/backend/internal/app/restaurant"
	"github.com/tongxuanbao/food-delivery/backend/internal/app/simulator"
	"github.com/tongxuanbao/food-delivery/backend/internal/pkg/geo"
)

type RateResponse struct {
	Id       int              `json:"id"`
	Route    []geo.Coordinate `json:"route"`
	Position geo.Coordinate   `json:"position"`
}

type contextKey string

const connectionIDKey contextKey = "connectionID"

// Middleware to generate unique connection ID
func connectionIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Generate a new UUID
		connID := uuid.New().String()
		connID = connID[:8]

		// Add the connection ID to the request context
		ctx := context.WithValue(r.Context(), connectionIDKey, connID)

		// Call the next handler with the updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func handleRoute(restaurantService *restaurant.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers to allow all origins. You may want to restrict this to specific origins in a production environment.
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Type")

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		// Get context and connectionID
		ctx := r.Context()
		connectionID := ctx.Value(connectionIDKey)
		fmt.Printf("%s [%s] /route CONNECTED\n", time.Now().Format("2006-01-02 15:04:05"), connectionID)

		// Send initial restaurant data
		jsonBytes, err := json.Marshal(restaurant.RestaurantList)
		if err == nil {
			fmt.Fprintf(w, "event: initial\ndata: %s\n\n", fmt.Sprint(string(jsonBytes)))
		}
		w.(http.Flusher).Flush()

		// Subscribe to services
		restaurantSubscriber := restaurantService.Broker.Subscribe("restaurant")

		// Listen for events
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
			case <-ctx.Done():
				fmt.Printf("%s [%s] /route DISCONNECTED AND UNSUBSCRIBED\n", time.Now().Format("2006-01-02 15:04:05"), connectionID)
				restaurantService.Broker.Unsubscribe("restaurant", restaurantSubscriber)
				return
			}
		}
	}
}

func NewServer(restaurantService *restaurant.Service) http.Handler {
	mux := http.NewServeMux()
	// Add routes
	mux.HandleFunc("/route", handleRoute(restaurantService))
	// Middle where
	var handler http.Handler = mux
	handler = connectionIDMiddleware(handler)
	return handler
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
