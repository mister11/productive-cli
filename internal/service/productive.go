package service

import (
	"errors"
	"github.com/mister11/productive-cli/internal/infrastructure/log"
	"github.com/mister11/productive-cli/internal/productive"
	"github.com/mister11/productive-cli/internal/service/tracking"
	"time"
)

type ProductiveService struct {
	client             *productive.Client
	userSessionManager UserSessionManager
}

type Project struct {
	deal    productive.Deal
	service productive.Service
}

func NewProductiveService(client *productive.Client) *ProductiveService {
	userSessionManager := NewFileUserSessionManager()
	return &ProductiveService{
		client:             client,
		userSessionManager: userSessionManager,
	}
}

func (s *ProductiveService) CreateTimeEntry(entry tracking.FoodEntry) error {
	sessionData, err := s.userSessionManager.GetUserSession()
	if err != nil {
		return err
	}
	foodService, err := s.FindFoodService(entry.Day)
	if err != nil {
		return err
	}
	log.Debug("Creating food time entry for day %v", entry.Day)
	return s.client.TimeEntryService.CreateTimeEntry(
		"", "30", sessionData.PersonID, foodService, entry.Day, sessionData.Token,
	)
}

func (s *ProductiveService) FindFoodService(day time.Time) (*productive.Service, error) {
	sessionData, err := s.userSessionManager.GetUserSession()
	if err != nil {
		return nil, err
	}
	deals, err := s.client.DealService.SearchDeals("Operations general", day, day, sessionData.Token)
	if err != nil {
		return nil, err
	}
	if len(deals) > 1 {
		return nil, errors.New("multiple 'Operations general' deals found")
	}
	if len(deals) == 0 {
		return nil, errors.New("no 'Operations general' deals found")
	}
	deal := deals[0]
	services, err := s.client.ServiceService.SearchServices("Food", deal.ID, day, day, sessionData.Token)
	if err != nil {
		return nil, err
	}
	if len(services) > 1 {
		return nil, errors.New("multiple 'Operations general/Food' services found")
	}
	if len(services) == 0 {
		return nil, errors.New("no 'Operations general/Food' services found")
	}
	return &services[0], nil
}
