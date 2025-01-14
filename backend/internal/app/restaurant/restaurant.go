package restaurant

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/tongxuanbao/food-delivery/backend/internal/pkg/broker"
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

type Subscriber interface {
	Trigger(data string)
}

type Service struct {
	Broker *broker.Broker
}

func New() *Service {
	s := &Service{
		Broker: broker.New(),
	}
	return s
}

// Add an order to restaurant
func (s *Service) AddOrder(customerId int, restaurantId int) {
	fmt.Printf("Restaurant: Add order. customerId: %d, restaurantId: %d\n", customerId, restaurantId)
	message := fmt.Sprintf("Restaurant is cooking. CustomerId: %d, restaurantId %d", customerId, restaurantId)
	s.Broker.Publish("restaurant", message)
	go func() {
		time.Sleep(time.Duration(rand.Intn(5)+10) * time.Second)
		s.orderPrepared(customerId, restaurantId)
	}()
}

func (s *Service) orderPrepared(customerId int, restaurantId int) {
	message := fmt.Sprintf("Restaurant cooking is FINISHED. CustomerId: %d, restaurantId %d", customerId, restaurantId)
	s.Broker.Publish("restaurant", message)
}

// Order picked up
func (s *Service) orderPickedUp(customerId int, restaurantId int) {
	fmt.Println("Order picked up")
	s.Broker.Publish("restaurant", "Order picked up")
}
