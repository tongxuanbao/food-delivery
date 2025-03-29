package simulator

import (
	"fmt"

	"github.com/tongxuanbao/food-delivery/backend/internal/app/driver"
	"github.com/tongxuanbao/food-delivery/backend/internal/app/restaurant"
)

type Service struct {
	restaurantService *restaurant.Service
	driverService     *driver.Service
}

func New(restaurantService *restaurant.Service, driverService *driver.Service) *Service {
	return &Service{
		restaurantService: restaurantService,
		driverService:     driverService,
	}
}

func (s *Service) Simulate() {
	// t := time.NewTicker(3 * time.Second)
	// defer func() {
	// 	fmt.Println("Simulate deferring")
	// 	t.Stop()
	// }()

	list := len(s.restaurantService.GetRestaurants())

	for range list {
		s.addOrder()
	}
}

// Get random restaurant, driver and customer
// Setting up route for driver
// Free restaurant, driver and customer
func (s *Service) addOrder() {
	fmt.Println("Add new order")
	// Get Random Restaurant
	restaurantID := 1
	// Get Random Customer
	customerID := 1

	// Get A driver
	driver := s.driverService.FindDriver(restaurantID, customerID)
	s.driverService.DriverTraverseToTheEnd(driver)
}
