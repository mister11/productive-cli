package track

import (
	"github.com/mister11/productive-cli/internal/action"
	"github.com/mister11/productive-cli/internal/client"
	"github.com/mister11/productive-cli/internal/client/model"
	"github.com/mister11/productive-cli/internal/config"
	"github.com/mister11/productive-cli/internal/datetime"
	"github.com/mister11/productive-cli/internal/log"
	"github.com/mister11/productive-cli/internal/stdin"
	"github.com/mister11/productive-cli/internal/utils"
	"strings"
	"time"
)

type TrackingManager struct {
	trackingClient   client.TrackingClient
	stdIn            stdin.Stdin
	config           config.ConfigManager
	dateTimeProvider datetime.DateTimeProvider
}

func NewTrackingManager(
	trackingClient client.TrackingClient,
	stdIn stdin.Stdin,
	config config.ConfigManager,
	dateTimeProvider datetime.DateTimeProvider,
) *TrackingManager {
	return &TrackingManager{
		trackingClient:   trackingClient,
		stdIn:            stdIn,
		config:           config,
		dateTimeProvider: dateTimeProvider,
	}
}

func (manager *TrackingManager) TrackFood(trackFoodRequest action.TrackFoodRequest) {
	if !trackFoodRequest.IsValid() {
		log.Error("You've provided both week and day tracking so I don't know what to do.")
		return
	}

	if trackFoodRequest.IsWeekTracking {
		manager.trackFood(manager.dateTimeProvider.GetWeekDays()...)
	} else if trackFoodRequest.Day != "" {
		date := manager.dateTimeProvider.ToISOTime(trackFoodRequest.Day)
		manager.trackFood([]time.Time{date}...)
	} else {
		manager.trackFood([]time.Time{manager.dateTimeProvider.Now()}...)
	}
}

func (manager *TrackingManager) trackFood(days ...time.Time) {
	userID := manager.config.GetUserID()
	for _, day := range days {
		dayFormatted := manager.dateTimeProvider.Format(day)
		log.Info("Tracking food for " + dayFormatted)
		service := manager.findFoodService(dayFormatted)
		timeEntry := model.NewTimeEntry("", 30, userID, service, dayFormatted)
		manager.trackingClient.CreateTimeEntry(timeEntry)
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
	dayFormatted := manager.dateTimeProvider.Format(day)
	if existingProject != nil {
		project := existingProject.(config.Project)
		manager.trackSavedProject(project, dayFormatted)
	} else {
		manager.trackNewProject(dayFormatted)
	}
}

func (manager *TrackingManager) trackSavedProject(project config.Project, dayFormatted string) {
	manager.config.RemoveExistingProject(project)
	deal, service := manager.findProjectInfo(project.DealName, project.ServiceName, dayFormatted)
	manager.trackSelectedProject(deal, service, dayFormatted)
	manager.config.SaveProject(config.NewProject(*deal, *service))
}

func (manager *TrackingManager) trackNewProject(dayFormatted string) {
	deal := manager.findAndSelectDeal(dayFormatted)
	service := manager.findAndSelectService(deal, dayFormatted)
	manager.trackSelectedProject(deal, service, dayFormatted)
	manager.config.SaveProject(config.NewProject(*deal, *service))
}

func (manager *TrackingManager) trackSelectedProject(deal *model.Deal, service *model.Service, dayFormatted string) {
	duration := manager.stdIn.Input("Time")
	durationParsed := utils.ParseTime(duration)
	notes := manager.createNotes()
	log.Info("Tracking %s - %s with duration %d for %d", deal.Name, service.Name, durationParsed, dayFormatted)
	timeEntry := model.NewTimeEntry(notes, durationParsed, manager.config.GetUserID(), service, dayFormatted)
	manager.trackingClient.CreateTimeEntry(timeEntry)
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

func (manager *TrackingManager) findProjectInfo(
	dealName string,
	serviceName string,
	dayFormatted string,
) (*model.Deal, *model.Service) {
	deal := manager.trackingClient.SearchDeals(dealName, dayFormatted)[0].(*model.Deal)
	service := manager.trackingClient.SearchServices(serviceName, deal.ID, dayFormatted)[0].(*model.Service)
	return deal, service
}

func (manager *TrackingManager) findFoodService(dayFormatted string) *model.Service {
	_, service := manager.findProjectInfo("Operations general", "Food", dayFormatted)
	return service
}

func (manager *TrackingManager) findAndSelectDeal(dayFormatted string) *model.Deal {
	dealQuery := manager.stdIn.Input("Search project")
	deals := manager.trackingClient.SearchDeals(dealQuery, dayFormatted)
	return manager.stdIn.SelectOne("Select project", deals).(*model.Deal)
}

func (manager *TrackingManager) findAndSelectService(deal *model.Deal, dayFormatted string) *model.Service {
	serviceQuery := manager.stdIn.Input("Search service")
	services := manager.trackingClient.SearchServices(serviceQuery, deal.ID, dayFormatted)
	return manager.stdIn.SelectOne("Select service", services).(*model.Service)
}

func (manager *TrackingManager) createNotes() string {
	notes := manager.stdIn.InputMultiple("Note")
	if len(notes) == 0 {
		return ""
	}
	var notesHTML strings.Builder
	notesHTML.WriteString("<ul>")
	for _, note := range notes {
		notesHTML.WriteString("<li>" + note + "</li>")
	}
	notesHTML.WriteString("</ul>")
	return notesHTML.String()
}

func searchProjectFunction(projects []config.Project) func(string, int) bool {
	return func(query string, index int) bool {
		project := projects[index]
		return strings.Contains(project.DealName, query) || strings.Contains(project.ServiceName, query)
	}
}
