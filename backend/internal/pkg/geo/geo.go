package geo

import (
	_ "embed"
	"encoding/json"
	"log"
	"math"
	"math/rand"
)

//go:embed adjacentListPixel.json
var pixelData []byte

//go:embed restaurantCustomerCleanPixel.json
var routeData []byte

type Coordinate struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Route struct {
	RestaurantId int          `json:"restaurantId"`
	CustomerId   int          `json:"customerId"`
	Route        []Coordinate `json:"route"`
}

var adjacentList = make(map[Coordinate][]Coordinate)
var coordinateList []Coordinate

var RouteList []Route

func (c Coordinate) GetNeighbors() []Coordinate {
	return adjacentList[c]
}

func init() {
	// Unmarshal the JSON into the connections variable
	var connections [][][2]float64
	err := json.Unmarshal(pixelData, &connections)
	if err != nil {
		log.Fatal(err)
	}

	// Create a map to store
	adjacentList = make(map[Coordinate][]Coordinate)
	for _, conn := range connections {
		// Extract two points
		pointA := Coordinate{X: conn[0][0] * 3.125, Y: conn[0][1] * 3.125}
		pointB := Coordinate{X: conn[1][0] * 3.125, Y: conn[1][1] * 3.125}

		// Create map of those points
		if _, exists := adjacentList[pointA]; !exists {
			adjacentList[pointA] = make([]Coordinate, 0)
			coordinateList = append(coordinateList, pointA)
		}
		if _, exists := adjacentList[pointB]; !exists {
			adjacentList[pointB] = make([]Coordinate, 0)
			coordinateList = append(coordinateList, pointB)
		}

		// Add connecting points to their list
		adjacentList[pointA] = append(adjacentList[pointA], pointB)
		adjacentList[pointB] = append(adjacentList[pointB], pointA)
	}

	// Unmarshal the JSON into the connections variable
	err = json.Unmarshal(routeData, &RouteList)
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

func GetAdjacentList() map[Coordinate][]Coordinate {
	return adjacentList
}

func GetRandomCoordinate() Coordinate {
	return coordinateList[rand.Intn(len(coordinateList))]
}

func visitCoordinate(current Coordinate, end Coordinate, visited map[Coordinate]bool) []Coordinate {
	visited[current] = true

	if current == end {
		return []Coordinate{current}
	}

	for _, neighbor := range current.GetNeighbors() {
		if visited[neighbor] {
			continue
		}

		route := visitCoordinate(neighbor, end, visited)
		if route != nil {
			return append([]Coordinate{current}, route...)
		}
	}

	return nil
}

func GetRoute(start Coordinate, end Coordinate) []Coordinate {
	return visitCoordinate(start, end, make(map[Coordinate]bool))
}

func GetCoordinateList() []Coordinate {
	return coordinateList
}

func (c *Coordinate) DistanceTo(coord Coordinate) float64 {
	x := math.Abs(c.X - coord.X)
	y := math.Abs(c.Y - coord.Y)

	return x*x + y*y

}
