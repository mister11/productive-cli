package tracking

import (
	"errors"
	"github.com/mister11/productive-cli/internal/domain/datetime"
	"github.com/mister11/productive-cli/internal/interactive"
	"github.com/mister11/productive-cli/internal/productive"
	"github.com/mister11/productive-cli/internal/service"
	"time"
)

type FoodEntry struct {
	Day time.Time
}

type FoodTrackingService struct {
	productiveService *service.ProductiveService
	prompt            interactive.Prompt
	sessionService    *service.SessionService
	dateTimeProvider  datetime.DateTimeProvider
}

func NewFoodTrackingService(productiveClient *productive.Client) *FoodTrackingService {
	sessionManager := service.NewFileUserSessionManager()
	sessionService := service.NewSessionService(productiveClient, sessionManager)
	stdInPrompt := interactive.NewStdinPrompt()
	dateTimeProvider := datetime.NewRealTimeDateProvider()

	return &FoodTrackingService{
		productiveService: service.NewProductiveService(productiveClient),
		prompt:            stdInPrompt,
		sessionService:    sessionService,
		dateTimeProvider:  dateTimeProvider,
	}
}

func (service *FoodTrackingService) TrackFood(request service.TrackFoodRequest) error {
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
		if err := service.productiveService.CreateTimeEntry(entry); err != nil {
			return err
		}
	}
	return nil
}

func (service *FoodTrackingService) getTrackingDays(request service.TrackFoodRequest) []time.Time {
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
	username, err := service.prompt.Input("Username")
	if err != nil {
		return err
	}
	password, err := service.prompt.InputMasked("Password")
	if err != nil {
		return err
	}
	return service.sessionService.Login(username, password)
}
