package driver

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/tongxuanbao/food-delivery/backend/internal/pkg/broker"
	"github.com/tongxuanbao/food-delivery/backend/internal/pkg/geo"
)

const (
	DRIVER_STATUS_WAITING    = 0
	DRIVER_STATUS_EN_ROUTE   = 1
	DRIVER_STATUS_DELIVERING = 2
)

// Coordinate is an alias for geo.Coordinate
type Coordinate = geo.Coordinate

type Driver struct {
	ID              int          `json:"id"`
	Coordinate      Coordinate   `json:"coordinate"`
	Route           []Coordinate `json:"route"`
	CurrentPosition int          `json:"currentPosition"`
	Speed           int          `json:"speed"`
	Status          int          `json:"status"`
}

type Service struct {
	mutex   sync.Mutex
	Broker  *broker.Broker
	Drivers []Driver
}

type DriverInitMessage struct {
	Event   string   `json:"event"`
	Drivers []Driver `json:"drivers"`
}

func New() *Service {
	s := &Service{Broker: broker.New()}
	return s
}

func (s *Service) generateRandomDriver() Driver {
	coordinate := geo.GetRandomCoordinate()
	speed := rand.Intn(20) + 20
	return Driver{Coordinate: coordinate, Speed: speed}
}

func (s *Service) SetDrivers(numOfDrivers int) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	fmt.Printf("%s DRIVER SetDrivers(%d) from %d \n", time.Now().Format("2006-01-02 15:04:05"), numOfDrivers, len(s.Drivers))

	if numOfDrivers > len(s.Drivers) {
		for range numOfDrivers - len(s.Drivers) {
			randomDriver := s.generateRandomDriver()
			randomDriver.ID = len(s.Drivers) + 1
			s.Drivers = append(s.Drivers, randomDriver)
		}
	} else {
		s.Drivers = s.Drivers[:numOfDrivers]
	}

	jsonBytes, err := json.Marshal(DriverInitMessage{Event: "init_drivers", Drivers: s.Drivers})
	if err == nil {
		message := fmt.Sprint(string(jsonBytes))
		s.Broker.Publish("driver", message)
	}
}

func (s *Service) GetDriver(customerID int, restaurantID int) {
	// Get what driver should deliver (might be closest available)
	// Schedule for driver to get to the restaurant (Current position -> restaurant)
	// Schedule for driver to deliver the (Restaurant -> customer)
	// Finished then go to random point and wait for next order
}
