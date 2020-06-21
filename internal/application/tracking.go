package application

import (
	"github.com/mister11/productive-cli/internal/domain/tracking"
	"github.com/mister11/productive-cli/internal/infrastructure/input"
	"time"

	"github.com/mister11/productive-cli/internal/config"
	"github.com/mister11/productive-cli/internal/domain/datetime"
	"github.com/mister11/productive-cli/internal/infrastructure/client"
)

type TrackingService struct {
	foodTracker      tracking.FoodTracker
	projectTracker   tracking.ProjectTracker
	dateTimeProvider datetime.DateTimeProvider
}

func NewTrackingService() *TrackingService {
	stdIn := input.NewStdinPrompt()
	configManager := config.NewFileConfigManager()
	dateTimeProvider := datetime.NewRealTimeDateProvider()
	productiveClient := client.NewProductiveClient(configManager)

	return &TrackingService{
		foodTracker:      tracking.NewHTTPFoodTracker(productiveClient, configManager),
		projectTracker:   tracking.NewHTTPProjectTracker(productiveClient, stdIn, configManager, dateTimeProvider),
		dateTimeProvider: dateTimeProvider,
	}
}

func (service *TrackingService) TrackFood(trackFoodRequest tracking.TrackFoodRequest) {
	service.foodTracker.TrackFood(trackFoodRequest)
}

func (service *TrackingService) TrackProject(request tracking.TrackProjectRequest) {
	var day time.Time
	if request.Day != "" {
		// this also verifies a format
		day = service.dateTimeProvider.ToISOTime(request.Day)
	} else {
		day = service.dateTimeProvider.Now()
	}
	service.projectTracker.TrackProject(day)
}
