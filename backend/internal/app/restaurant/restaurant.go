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

type Restaurant struct {
	Id         int        `json:"id"`
	Coordinate Coordinate `json:"coordinate"`
}

type Event struct {
	Event      string     `json:"event"`
	Restaurant Restaurant `json:"restaurant"`
}

//go:embed restaurants.json
var restaurantsData []byte

var RestaurantList []Restaurant

func init() {
	// Unmarshal the JSON into the connections variable
	var listFromJson []Coordinate
	err := json.Unmarshal(restaurantsData, &listFromJson)
	if err != nil {
		log.Fatal(err)
	}
	RestaurantList = make([]Restaurant, 0)
	for idx, restaurant := range listFromJson {
		if len(RestaurantList) == 10 {
			break
		}
		if restaurant.X < 0 || restaurant.X > 1920 {
			continue
		}
		if restaurant.Y < 0 || restaurant.Y > 1080 {
			continue
		}
		coordinate := Coordinate{X: restaurant.X * 3.125, Y: restaurant.Y * 3.125}
		RestaurantList = append(RestaurantList, Restaurant{Id: idx, Coordinate: coordinate})
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

func (s *Service) GetRestaurants() []Restaurant {
	return RestaurantList
}

func (s *Service) getRestaurant(id int) Restaurant {
	return s.GetRestaurants()[id]
}

func (s *Service) constructEventMessage(eventName string, restaurantID int) string {
	// Get restaurant
	restaurant := s.getRestaurant(restaurantID)

	// Construct event
	e := Event{Event: eventName, Restaurant: restaurant}

	// Turn event into json string
	jsonBytes, err := json.Marshal(e)
	if err != nil {
		fmt.Printf("%s RESTAURANT Generating %s Got error: %v \n", time.Now().Format("2006-01-02 15:04:05"), eventName, err)
		return ""
	}
	message := fmt.Sprint(string(jsonBytes))

	return message
}

// Add an order to restaurant
func (s *Service) AddOrder(customerID int, restaurantID int) {
	// Log
	fmt.Printf("%s RESTAURANT AddOrder(customerID: %d, restaurantID: %d) \n", time.Now().Format("2006-01-02 15:04:05"), customerID, restaurantID)

	// Publish action
	message := s.constructEventMessage("order_received", restaurantID)
	s.Broker.Publish("restaurant", message)

	// Schedule prepared action
	go func() {
		time.Sleep(time.Duration(rand.Intn(5)+5) * time.Second)
		s.orderPrepared(customerID, restaurantID)
	}()
}

func (s *Service) orderPrepared(customerID int, restaurantID int) {
	// Log
	fmt.Printf("%s RESTAURANT orderPrepared(customerID: %d, restaurantID: %d) \n", time.Now().Format("2006-01-02 15:04:05"), customerID, restaurantID)

	// Publish action
	message := s.constructEventMessage("order_prepared", restaurantID)
	s.Broker.Publish("restaurant", message)
}

// Order picked up
func (s *Service) orderPickedUp(customerId int, restaurantId int) {
	fmt.Println("Order picked up")
	s.Broker.Publish("restaurant", "Order picked up")
}
