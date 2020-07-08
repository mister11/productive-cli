package application

import (
	"github.com/mister11/productive-cli/internal/domain/tracking"
	"github.com/mister11/productive-cli/internal/infrastructure/client"
	"github.com/mister11/productive-cli/internal/infrastructure/input"
	"github.com/mister11/productive-cli/internal/infrastructure/log"

	"github.com/mister11/productive-cli/internal/domain/datetime"
	"github.com/mister11/productive-cli/internal/infrastructure/config"
)

type TrackingService struct {
	foodEntriesCreator  tracking.FoodEntriesCreator
	projectEntryCreator tracking.ProjectEntryCreator
	trackingClient      client.TrackingClient
	prompt              *input.StdinPrompt
}

func NewTrackingService() *TrackingService {
	prompt := input.NewStdinPrompt()
	configManager := config.NewFileConfigManager()
	dateTimeProvider := datetime.NewRealTimeDateProvider()
	trackingClient := client.NewProductiveClient(configManager)
	searcher := tracking.NewProjectSearcher(prompt, trackingClient)

	return &TrackingService{
		foodEntriesCreator:  tracking.NewFoodEntriesCreator(dateTimeProvider, configManager),
		projectEntryCreator: tracking.NewProjectEntryCreator(dateTimeProvider, prompt, configManager, searcher),
		trackingClient:      trackingClient,
		prompt:              prompt,
	}
}

func (service *TrackingService) TrackFood(request tracking.TrackFoodRequest) error {
	if err := service.loginIfNeeded(); err != nil {
		return err
	}
	foodEntries, err := service.foodEntriesCreator.Create(request)
	if err != nil {
		return err
	}
	return service.trackingClient.TrackFood(foodEntries)
}

func (service *TrackingService) TrackProject(request tracking.TrackProjectRequest) error {
	if err := service.loginIfNeeded(); err != nil {
		return err
	}
	projectEntry, err := service.projectEntryCreator.Create(request)
	if err != nil {
		return nil
	}
	return service.trackingClient.TrackProject(projectEntry)
	//var day time.Time
	//if request.Day != "" {
	//	// this also verifies a format
	//	day = service.dateTimeProvider.ToISOTime(request.Day)
	//} else {
	//	day = service.dateTimeProvider.Now()
	//}
	//service.projectTracker.TrackProject(day)
	//factory.config.SaveProject(config.NewProject(*deal, *service))
}

func (service *TrackingService) loginIfNeeded() error {
	loginStatus, err := service.trackingClient.VerifyLogin()
	if err != nil {
		log.Error("Cannot verify login status. Please re-login.")
	}
	if loginStatus != "ok" {
		username, err := service.prompt.Input("Username")
		if err != nil {
			return err
		}
		password, err := service.prompt.InputMasked("Password")
		if err != nil {
			return err
		}
		return service.trackingClient.Login(username, password)
	}
	return nil
}
