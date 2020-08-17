package service

import (
	"errors"
	"time"

	"github.com/mister11/productive-cli/internal/productive"
	"github.com/mister11/productive-cli/internal/service/datetime"
)

type FoodEntry struct {
	Day time.Time
}

type FoodTrackingService struct {
	productiveService *ProductiveService
	prompt            Prompt
	sessionService    *SessionService
	dateTimeProvider  datetime.DateTimeProvider
}

func NewFoodTrackingService(productiveClient *productive.Client) *FoodTrackingService {
	sessionManager := NewFileUserSessionManager()
	productiveService := NewProductiveService(productiveClient)
	stdInPrompt := NewStdinPrompt()
	sessionService := NewSessionService(productiveService, stdInPrompt, sessionManager)
	dateTimeProvider := datetime.NewRealTimeDateProvider()

	return &FoodTrackingService{
		productiveService: productiveService,
		prompt:            stdInPrompt,
		sessionService:    sessionService,
		dateTimeProvider:  dateTimeProvider,
	}
}

func (s *FoodTrackingService) TrackFood(request TrackFoodRequest) error {
	userSession, err := s.sessionService.ObtainUserSession()
	if err != nil {
		return err
	}
	if !request.IsValid() {
		return errors.New("invalid track food request")
	}
	days := s.getTrackingDays(request)
	var entries []FoodEntry
	for _, day := range days {
		entry := FoodEntry{Day: day}
		entries = append(entries, entry)
	}
	for _, entry := range entries {
		if err := s.productiveService.CreateFoodTimeEntry(entry, userSession); err != nil {
			return err
		}
	}
	return nil
}

func (s *FoodTrackingService) getTrackingDays(request TrackFoodRequest) []time.Time {
	if request.IsWeekTracking {
		return s.dateTimeProvider.GetWeekDays()
	} else if request.Day != "" {
		return []time.Time{s.dateTimeProvider.ToISOTime(request.Day)}
	} else {
		return []time.Time{s.dateTimeProvider.Now()}
	}
}
