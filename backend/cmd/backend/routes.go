package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tongxuanbao/food-delivery/backend/internal/app/customer"
	"github.com/tongxuanbao/food-delivery/backend/internal/app/driver"
	"github.com/tongxuanbao/food-delivery/backend/internal/app/restaurant"
)

func addRoutes(
	mux *http.ServeMux,
	restaurantService *restaurant.Service,
	driverService *driver.Service,
	customerService *customer.Service,
) {
	mux.HandleFunc("/route", handleRoute(restaurantService, driverService, customerService))
	mux.HandleFunc("POST /customers", handleCustomers(restaurantService, driverService, customerService))
	mux.Handle("/metrics", promhttp.Handler())
}
