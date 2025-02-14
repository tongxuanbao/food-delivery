package customer

import (
	_ "embed"
	"encoding/json"
	"log"

	"github.com/tongxuanbao/food-delivery/backend/internal/pkg/broker"
	"github.com/tongxuanbao/food-delivery/backend/internal/pkg/geo"
)

// Coordinate is an alias for geo.Coordinate
type Coordinate = geo.Coordinate

//go:embed customers.json
var customersData []byte

type Customer struct {
	Id         int        `json:"id"`
	Coordinate Coordinate `json:"coordinate"`
}

var CustomerList []Customer

func init() {
	// Unmarshal the JSON into the connections variable
	var listFromJson []Coordinate
	err := json.Unmarshal(customersData, &listFromJson)
	if err != nil {
		log.Fatal(err)
	}
	for idx, customer := range listFromJson {
		if customer.X < 0 || customer.X > 1920 {
			continue
		}
		if customer.Y < 0 || customer.Y > 1080 {
			continue
		}
		coordinate := Coordinate{X: customer.X * 3.125, Y: customer.Y * 3.125}
		CustomerList = append(CustomerList, Customer{Id: idx, Coordinate: coordinate})
	}
}

type Service struct {
	Broker *broker.Broker
}

func New() *Service {
	s := &Service{Broker: broker.New()}
	return s
}

func (s *Service) GetCustomers() []Customer {
	return CustomerList
}
