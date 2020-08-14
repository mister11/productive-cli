package service

import (
	"errors"
	"github.com/mister11/productive-cli/internal/productive"
	"github.com/mister11/productive-cli/internal/service/datetime"
	"time"
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
	sessionService := NewSessionService(productiveClient, sessionManager)
	stdInPrompt := NewStdinPrompt()
	dateTimeProvider := datetime.NewRealTimeDateProvider()

	return &FoodTrackingService{
		productiveService: NewProductiveService(productiveClient),
		prompt:            stdInPrompt,
		sessionService:    sessionService,
		dateTimeProvider:  dateTimeProvider,
	}
}

func (service *FoodTrackingService) TrackFood(request TrackFoodRequest) error {
	if err := service.loginIfNeeded(); err != nil {
		return err
	}
	if !request.IsValid() {
		return errors.New("invalid track food request")
	}
	days := service.getTrackingDays(request)
	var entries []FoodEntry
	for _, day := range days {
		entry := FoodEntry{Day: day}
		entries = append(entries, entry)
	}
	for _, entry := range entries {
		if err := service.productiveService.CreateFoodTimeEntry(entry); err != nil {
			return err
		}
	}
	return nil
}

func (service *FoodTrackingService) getTrackingDays(request TrackFoodRequest) []time.Time {
	if request.IsWeekTracking {
		return service.dateTimeProvider.GetWeekDays()
	} else if request.Day != "" {
		return []time.Time{service.dateTimeProvider.ToISOTime(request.Day)}
	} else {
		return []time.Time{service.dateTimeProvider.Now()}
	}
}

func (service *FoodTrackingService) loginIfNeeded() error {
	isSessionValid, err := service.sessionService.IsSessionValid()
	if err != nil {
		return err
	}
	// session is valid, we don't need login and there's no error
	if isSessionValid {
		return nil
	}
	username, err := service.prompt.Input("E-mail")
	if err != nil {
		return err
	}
	password, err := service.prompt.InputMasked("Password")
	if err != nil {
		return err
	}
	return service.sessionService.Login(username, password)
}
