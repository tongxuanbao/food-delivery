package driver

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"slices"
	"sync"
	"time"

	"github.com/tongxuanbao/food-delivery/backend/internal/pkg/assert"
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
	RouteIndex      int          `json:"routeIndex"`
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

func newDriverInitMessage(drivers []Driver) DriverInitMessage {
	return DriverInitMessage{
		Event:   "init_drivers",
		Drivers: drivers,
	}
}

type DriverUpdateMessage struct {
	Event  string `json:"event"`
	Driver Driver `json:"driver"`
}

func newDriverUpdateMessage(driver *Driver) DriverUpdateMessage {
	return DriverUpdateMessage{
		Event:  "driver",
		Driver: *driver,
	}
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
		for i := range numOfDrivers - len(s.Drivers) {
			randomDriver := s.generateRandomDriver()
			randomDriver.ID = len(s.Drivers) + 1
			randomDriver.Route = geo.GetRouteByIndex(i*100 + i).Route
			slices.Reverse(randomDriver.Route)
			fmt.Printf("restaurant: %d, customer: %d\n", geo.GetRouteByIndex(i*100+i).RestaurantID, geo.GetRouteByIndex(i*100+i).CustomerID)
			s.Drivers = append(s.Drivers, randomDriver)
		}
	} else {
		s.Drivers = s.Drivers[:numOfDrivers]
	}

	jsonBytes, err := json.Marshal(newDriverInitMessage(s.Drivers))
	if err == nil {
		message := fmt.Sprint(string(jsonBytes))
		s.Broker.Publish("driver", message)
	}
}

func (s *Service) getClosestAvailableDriver(location Coordinate) *Driver {
	var currentDriver *Driver
	for _, driver := range s.Drivers {
		// Not available
		if driver.Status != DRIVER_STATUS_WAITING {
			continue
		}

		// First available driver
		if currentDriver == nil {
			currentDriver = &driver
		}

		// Closer available driver
		distanceToCurrent := location.DistanceTo(currentDriver.Coordinate)
		distanceToDriver := location.DistanceTo(driver.Coordinate)
		if distanceToCurrent > distanceToDriver {
			currentDriver = &driver
		}
	}
	return currentDriver
}

// Route is sorted by index, cusId + resID*100
func (s *Service) FindDriver(restaurantID int, customerID int) *Driver {
	assert.Assert(restaurantID >= 0 && restaurantID < 100, "Invalid restaurant ID: must be between 0 and 99")
	assert.Assert(customerID >= 0 && customerID < 100, "Invalid customer ID: must be between 0 and 99"))

	// Get main route
	routeIndex := customerID + restaurantID*100
	mainRoute := geo.GetRouteByIndex(routeIndex)
	assert.Assert(
		mainRoute.RestaurantID == restaurantID && mainRoute.CustomerID == customerID,
		"Mismatch: Retrieved route does not correspond to the provided restaurant and customer IDs",
	)

	// Get the closest available driver
	driver := s.getClosestAvailableDriver(mainRoute.Route[0])
	// There's should be enough available Driver
	assert.Assert(driver != nil, "No available driver found")
	return driver
}

// Get driver to traverse to the end of the route, updating event in the meanwhile
func (s *Service) DriverTraverseToTheEnd(driver *Driver) {
	for i := range len(driver.Route) - 1 {
		driver.CurrentPosition = i
		driver.Coordinate = driver.Route[driver.CurrentPosition]

		jsonBytes, err := json.Marshal(newDriverUpdateMessage(driver))
		if err == nil {
			message := fmt.Sprint(string(jsonBytes))
			s.Broker.Publish("driver", message)
		}

		cl := clamp(driver.CurrentPosition+1, 0, len(driver.Route))
		distance := driver.Coordinate.DistanceTo(driver.Route[cl])

		time.Sleep(time.Duration(math.Sqrt(distance)) * time.Millisecond * 10)
	}
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

// func (s *Service) testDriver(index int) {
// 	driver := s.Drivers[index]
// 	for i := range len(driver.Route) - 1 {
// 		driver.CurrentPosition = i
// 		driver.Coordinate = driver.Route[driver.CurrentPosition]

// 		jsonBytes, err := json.Marshal(newDriverUpdateMessage(driver))
// 		if err == nil {
// 			message := fmt.Sprint(string(jsonBytes))
// 			s.Broker.Publish("driver", message)
// 		}

// 		cl := clamp(driver.CurrentPosition+1, 0, len(driver.Route))
// 		distance := driver.Coordinate.DistanceTo(driver.Route[cl])

// 		time.Sleep(time.Duration(math.Sqrt(distance)) * time.Millisecond * 10)
// 	}
// }
