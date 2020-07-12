package application

import (
	"github.com/mister11/productive-cli/internal/domain"
	"github.com/mister11/productive-cli/internal/domain/tracking"
	"github.com/mister11/productive-cli/internal/infrastructure/client"
	"github.com/mister11/productive-cli/internal/infrastructure/input"
	"github.com/mister11/productive-cli/internal/infrastructure/session"

	"github.com/mister11/productive-cli/internal/domain/datetime"
)

type TrackingService struct {
	foodEntriesCreator    tracking.FoodEntriesCreator
	projectEntryCreator   tracking.ProjectEntryCreator
	trackedProjectManager domain.TrackedProjectManager
	trackingClient        tracking.TrackingClient
	prompt                *input.StdinPrompt
	loginManager          domain.LoginManager
}

func NewTrackingService() *TrackingService {
	prompt := input.NewStdinPrompt()
	userConfigManager := client.NewFileUserSessionManager()
	trackedProjectManager := domain.NewFileTrackedProjectsManager()
	dateTimeProvider := datetime.NewRealTimeDateProvider()
	trackingClient := client.NewProductiveClient(userConfigManager)
	loginManager := session.NewProductiveLoginManager(trackingClient, userConfigManager)

	return &TrackingService{
		foodEntriesCreator:    tracking.NewFoodEntriesCreator(dateTimeProvider),
		projectEntryCreator:   tracking.NewProjectEntryCreator(dateTimeProvider, prompt, trackedProjectManager, trackingClient),
		trackedProjectManager: trackedProjectManager,
		trackingClient:        trackingClient,
		prompt:                prompt,
		loginManager:          loginManager,
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
		return err
	}
	if err := service.trackingClient.TrackProject(projectEntry); err != nil {
		return err
	}
	return service.trackedProjectManager.UpsertTrackedProject(domain.TrackedProject{
		DealID:      projectEntry.Deal.ID,
		DealName:    projectEntry.Deal.Name,
		ServiceID:   projectEntry.Service.ID,
		ServiceName: projectEntry.Service.Name,
	})
}

func (service *TrackingService) loginIfNeeded() error {
	isSessionValid, err := service.loginManager.IsSessionValid()
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
	return service.loginManager.Login(username, password)
}
