package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/tongxuanbao/food-delivery/backend/internal/app/customer"
	"github.com/tongxuanbao/food-delivery/backend/internal/app/driver"
	"github.com/tongxuanbao/food-delivery/backend/internal/app/restaurant"
)

type InitialMessage struct {
	Restaurants []restaurant.Restaurant `json:"restaurants"`
	Drivers     []driver.Driver         `json:"drivers"`
	Customers   []customer.Customer     `json:"customers"`
}

// "/route"
func handleRoute(restaurantService *restaurant.Service, driverService *driver.Service, customerService *customer.Service) func(http.ResponseWriter, *http.Request) {
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
		restaurantService.SetRestaurants(1)
		driverService.SetDrivers(1)
		customerService.SetCustomers(1)
		initialMessage := InitialMessage{
			Restaurants: restaurantService.Restaurants,
			Drivers:     driverService.Drivers,
			Customers:   customerService.Customers,
		}
		jsonBytes, err := json.Marshal(initialMessage)
		if err == nil {
			fmt.Fprintf(w, "event: initial\ndata: %s\n\n", fmt.Sprint(string(jsonBytes)))
		}
		w.(http.Flusher).Flush()

		// Subscribe to services
		restaurantSubscriber := restaurantService.Broker.Subscribe("restaurant")
		driverSubscriber := driverService.Broker.Subscribe("driver")
		// customerSubscriber := customerService.Broker.Subscribe("customer")

		// Listen for events
		for {
			select {
			case msg, ok := <-restaurantSubscriber.Channel:
				if !ok {
					fmt.Println("Subscriber channel closed.")
				}
				// Send data
				fmt.Fprintf(w, "event: restaurant\ndata: %s\n\n", msg)
				w.(http.Flusher).Flush()
			case msg := <-driverSubscriber.Channel:
				fmt.Fprintf(w, "event: driver\ndata: %s\n\n", msg)
				w.(http.Flusher).Flush()
			case <-ctx.Done():
				fmt.Printf("%s [%s] /route DISCONNECTED AND UNSUBSCRIBED\n", time.Now().Format("2006-01-02 15:04:05"), connectionID)
				restaurantService.Broker.Unsubscribe("restaurant", restaurantSubscriber)
				driverService.Broker.Unsubscribe("driver", driverSubscriber)
				return
			}
		}
	}
}

type CustomersRequestBody struct {
	NumberOfCustomers int `json:"numberOfCustomers"`
}

// "POST /drivers"
func handleCustomers(restaurantService *restaurant.Service, driverService *driver.Service, customerService *customer.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CustomersRequestBody
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			errMessage := fmt.Sprintf("Invalid JSON: %v", err)
			http.Error(w, errMessage, http.StatusBadRequest)
			return
		}
		if req.NumberOfCustomers < 1 {
			http.Error(w, "numberOfCustomers must be greater than 0", http.StatusBadRequest)
			return
		}

		restaurantService.SetRestaurants(req.NumberOfCustomers)
		driverService.SetDrivers(req.NumberOfCustomers)
		customerService.SetCustomers(req.NumberOfCustomers)

		response := map[string]any{
			"message":    "Number of customers has updated",
			"numDrivers": req.NumberOfCustomers,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
