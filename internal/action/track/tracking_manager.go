package track

import (
	"strings"
	"time"

	"github.com/mister11/productive-cli/internal/action"
	"github.com/mister11/productive-cli/internal/client"
	"github.com/mister11/productive-cli/internal/config"
	"github.com/mister11/productive-cli/internal/datetime"
	"github.com/mister11/productive-cli/internal/log"
	"github.com/mister11/productive-cli/internal/stdin"
	"github.com/mister11/productive-cli/internal/stdin/promptui"
)

type TrackingManager struct {
	foodTracker      foodTracker
	projectTracker   projectTracker
	stdIn            stdin.Stdin
	config           config.ConfigManager
	dateTimeProvider datetime.DateTimeProvider
}

func NewTrackingManager() *TrackingManager {
	stdIn := promptui.NewPromptUiStdin()
	configManager := config.NewFileConfigManager()
	dateTimeProvider := datetime.NewRealTimeDateProvider()
	productiveClient := client.NewProductiveClient(configManager)

	return &TrackingManager{
		foodTracker:      newHTTPFoodTracker(productiveClient, configManager),
		projectTracker:   newHTTPProjectTracker(productiveClient, stdIn, configManager, dateTimeProvider),
		stdIn:            stdIn,
		config:           configManager,
		dateTimeProvider: dateTimeProvider,
	}
}

func (manager *TrackingManager) TrackFood(trackFoodRequest action.TrackFoodRequest) {
	if !trackFoodRequest.IsValid() {
		log.Error("You've provided both week and day tracking so I don't know what to do.", nil)
		return
	}

	if trackFoodRequest.IsWeekTracking {
		manager.foodTracker.TrackFood(manager.dateTimeProvider.GetWeekDays()...)
	} else if trackFoodRequest.Day != "" {
		date := manager.dateTimeProvider.ToISOTime(trackFoodRequest.Day)
		manager.foodTracker.TrackFood([]time.Time{date}...)
	} else {
		manager.foodTracker.TrackFood([]time.Time{manager.dateTimeProvider.Now()}...)
	}
}

func (manager *TrackingManager) TrackProject(request action.TrackProjectRequest) {
	existingProject := manager.selectExistingProject()
	var day time.Time
	if request.Day != "" {
		// this also verifies a format
		day = manager.dateTimeProvider.ToISOTime(request.Day)
	} else {
		day = manager.dateTimeProvider.Now()
	}
	if existingProject != nil {
		project := existingProject.(config.Project)
		manager.projectTracker.TrackProject(day, &project)
	} else {
		manager.projectTracker.TrackProject(day, nil)
	}
}

func (manager *TrackingManager) selectExistingProject() interface{} {
	savedProjects := manager.config.GetSavedProjects()
	if len(savedProjects) == 0 {
		return nil
	}
	selectedProject := manager.stdIn.SelectOneWithSearch(
		"Select project",
		savedProjects,
		searchProjectFunction(savedProjects),
	)
	return selectedProject
}

func searchProjectFunction(projects []config.Project) func(string, int) bool {
	return func(query string, index int) bool {
		project := projects[index]
		return strings.Contains(project.DealName, query) || strings.Contains(project.ServiceName, query)
	}
}
