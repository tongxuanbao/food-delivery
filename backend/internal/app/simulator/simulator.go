package simulator

import (
	"fmt"
	"time"

	"math/rand"

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

	t := time.NewTicker(3 * time.Second)
	defer func() {
		fmt.Println("Simulate defering")
		t.Stop()
	}()

	for {
		<-t.C
		s.addOrder()
	}
}

func (s *Service) addOrder() {
	// Get random restaurant
	restaurants := s.restaurantService.GetRestaurants()
	randomRestaurantID := rand.Intn(len(restaurants))

	// // TODO: when adding user service
	randomUserID := 1
	s.restaurantService.AddOrder(randomUserID, randomRestaurantID)
}
