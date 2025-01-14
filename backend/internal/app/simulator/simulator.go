package simulator

import (
	"time"

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
	t := time.NewTicker(10 * time.Second)
	defer t.Stop()

	for {
		<-t.C
		s.restaurantService.AddOrder(1, 1)
	}
}
