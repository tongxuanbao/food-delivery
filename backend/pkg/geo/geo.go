package geo

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
)

//go:embed adjacentListPixel.json
var pixelData []byte

type Coordinate struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

var adjacentList = make(map[Coordinate][]Coordinate)
var coordinateList []Coordinate

func (c Coordinate) GetNeighbors() []Coordinate {
	// Access the adjacentList directly
	return adjacentList[c]
}

func init() {
	var connections [][][2]float64

	// Unmarshal the JSON into the connections variable
	err := json.Unmarshal(pixelData, &connections)
	if err != nil {
		log.Fatal(err)
	}

	// Create a map to store
	adjacentList = make(map[Coordinate][]Coordinate)
	for _, conn := range connections {
		// Extract two points
		pointA := Coordinate{X: conn[0][0], Y: conn[0][1]}
		pointB := Coordinate{X: conn[1][0], Y: conn[1][1]}

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

	// Print the result
	// fmt.Printf("%+v\n", coordinateList)
	fmt.Println(len(connections))
}

func GetAdjacentList() map[Coordinate][]Coordinate {
	return adjacentList
}

func GetCoordinateList() []Coordinate {
	return coordinateList
}
