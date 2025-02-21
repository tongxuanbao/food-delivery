package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/tongxuanbao/food-delivery/backend/internal/app/customer"
	"github.com/tongxuanbao/food-delivery/backend/internal/app/driver"
	"github.com/tongxuanbao/food-delivery/backend/internal/app/restaurant"
	"github.com/tongxuanbao/food-delivery/backend/internal/app/simulator"
	"github.com/tongxuanbao/food-delivery/backend/internal/pkg/geo"
)

type RateResponse struct {
	Id       int              `json:"id"`
	Route    []geo.Coordinate `json:"route"`
	Position geo.Coordinate   `json:"position"`
}

func NewServer(
	restaurantService *restaurant.Service,
	driverService *driver.Service,
	customerService *customer.Service,
) http.Handler {
	mux := http.NewServeMux()

	// Add routes
	addRoutes(mux, restaurantService, driverService, customerService)

	// Middlewares
	var handler http.Handler = mux
	handler = connectionIDMiddleware(handler)
	return handler
}

func main() {
	// Initialize services
	restaurantService := restaurant.New()
	driverService := driver.New()
	customerService := customer.New()

	restaurantService.SetRestaurants(1)
	driverService.SetDrivers(1)
	customerService.SetCustomers(1)

	serverHandler := NewServer(restaurantService, driverService, customerService)

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

	simulatorService := simulator.New(restaurantService, driverService)
	go simulatorService.Simulate()

	driverService.Test()

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
