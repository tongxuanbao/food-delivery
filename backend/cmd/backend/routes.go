package main

import (
	"net/http"

	"github.com/tongxuanbao/food-delivery/backend/internal/app/restaurant"
)

func addRoutes(
	mux *http.ServeMux,
	restaurantService *restaurant.Service,
) {
	mux.HandleFunc("/route", handleRoute(restaurantService))
}
