package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/tongxuanbao/food-delivery/backend/internal/app/restaurant"
)

// "/route"
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
