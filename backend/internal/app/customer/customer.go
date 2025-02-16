package customer

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

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

type CustomerInitMessage struct {
	Event     string     `json:"event"`
	Customers []Customer `json:"customers"`
}

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
	mutex     sync.Mutex
	Broker    *broker.Broker
	Customers []Customer
}

func New() *Service {
	s := &Service{Broker: broker.New()}
	return s
}

func (s *Service) GetCustomers() []Customer {
	return CustomerList
}

func (s *Service) SetCustomers(numOfCustomers int) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	fmt.Printf("%s CUSTOMER SetCustomers(%d) from %d \n", time.Now().Format("2006-01-02 15:04:05"), numOfCustomers, len(s.Customers))

	s.Customers = CustomerList[:numOfCustomers]

	jsonBytes, err := json.Marshal(CustomerInitMessage{Event: "init_customers", Customers: s.Customers})
	if err == nil {
		message := fmt.Sprint(string(jsonBytes))
		s.Broker.Publish("customer", message)
	}
}
