package geo

import (
	_ "embed"
	"encoding/json"
	"log"
	"math"
	"math/rand"
)

//go:embed restaurantCustomerCleanPixel.json
var routeData []byte

type Coordinate struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Route struct {
	RestaurantID int          `json:"restaurantId"`
	CustomerID   int          `json:"customerId"`
	Route        []Coordinate `json:"route"`
}

var RouteList []Route

func init() {
	err := json.Unmarshal(routeData, &RouteList)
	if err != nil {
		log.Fatal(err)
	}
	for i := range len(RouteList) {
		for j := range len(RouteList[i].Route) {
			RouteList[i].Route[j].X *= 3.125
			RouteList[i].Route[j].Y *= 3.125
		}
	}
}

func GetRandomRoute() Route {
	return RouteList[rand.Intn(len(RouteList))]
}

func GetRandomCoordinate() Coordinate {
	randomRoute := GetRandomRoute().Route
	return randomRoute[rand.Intn(len(randomRoute))]
}

func clamp(value int, min int, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func GetRouteByIndex(index int) Route {
	return RouteList[clamp(index, 0, len(RouteList)-1)]
}

func (c *Coordinate) DistanceTo(coord Coordinate) float64 {
	x := math.Abs(c.X - coord.X)
	y := math.Abs(c.Y - coord.Y)

	return x*x + y*y
}
