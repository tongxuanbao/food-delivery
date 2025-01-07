package restaurant

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/tongxuanbao/food-delivery/backend/internal/pkg/geo"
)

// Coordinate is an alias for geo.Coordinate
type Coordinate = geo.Coordinate

//go:embed restaurants.json
var restaurantsData []byte

var RestaurantList []Coordinate

func init() {
	// Unmarshal the JSON into the connections variable
	err := json.Unmarshal(restaurantsData, &RestaurantList)
	if err != nil {
		log.Fatal(err)
	}
	for idx, restaurant := range RestaurantList {
		RestaurantList[idx] = Coordinate{X: restaurant.X * 3.125, Y: restaurant.Y * 3.125}
	}
}

type Service struct {
	quit chan bool
}

func NewRestaurantService() *Service {
	s := Service{make(chan bool)}

	return s
}

func (s *Service) run() {
	for {
		select {
		case <-s.quit:
			fmt.Println("finishing task")
			time.Sleep(time.Second)
			fmt.Println("task done")
			s.quit <- true
			return
		case <-time.After(time.Second):
			fmt.Println("running task")
		}
	}
}

// Add an order to restaurant
func (s *Service) addOrder(customerId int, restaurantId int) {
	fmt.Printf("Restaurant is cooking. CustomerId: %d, restaurantId %d\n", customerId, restaurantId)
	time.Sleep(time.Duration(rand.Intn(3)+2) * time.Second)
	fmt.Printf("Restaurant cooking is FINISHED. CustomerId: %d, restaurantId %d\n", customerId, restaurantId)
}

// Order picked up
func (s *Service) orderPickedUp(customerId int, restaurantId int) {
	fmt.Println("Order picked up")
}

// Subscribe to all restaurant
func (s *Service) subscribeEvents(ch chan string) {
}

//
