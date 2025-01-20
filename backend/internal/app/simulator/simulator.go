package simulator

import (
	"time"

	"math/rand"

	"github.com/tongxuanbao/food-delivery/backend/internal/app/restaurant"
)

type Service struct {
	restaurantService *restaurant.Service
}

func New(restaurantService *restaurant.Service) *Service {
	return &Service{
		restaurantService: restaurantService,
	}
}

func (s *Service) Simulate() {
	t := time.NewTicker(time.Second)
	defer t.Stop()

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
