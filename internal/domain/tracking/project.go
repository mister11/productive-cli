package tracking

import (
	"errors"
	"github.com/mister11/productive-cli/internal/domain"
	"github.com/mister11/productive-cli/internal/interactive"
	"github.com/mister11/productive-cli/internal/service"
	"strings"
	"time"

	"github.com/mister11/productive-cli/internal/domain/datetime"
	"github.com/mister11/productive-cli/internal/infrastructure/log"
	"github.com/mister11/productive-cli/internal/utils"
)

type ProjectEntry struct {
	Service  *domain.Service
	Deal     *domain.Deal
	Day      time.Time
	Duration string
	Notes    []string
}

type ProjectEntryCreator interface {
	Create(request service.TrackProjectRequest) (*ProjectEntry, error)
}

type projectEntryFactory struct {
	dateTimeProvider     datetime.DateTimeProvider
	prompt               interactive.Prompt
	projectConfigManager domain.TrackedProjectManager
	trackingClient       TrackingClient
}

func NewProjectEntryCreator(
	dateTimeProvider datetime.DateTimeProvider,
	prompt interactive.Prompt,
	projectConfigManager domain.TrackedProjectManager,
	trackingClient TrackingClient,
) ProjectEntryCreator {
	return &projectEntryFactory{
		dateTimeProvider:     dateTimeProvider,
		prompt:               prompt,
		projectConfigManager: projectConfigManager,
		trackingClient:       trackingClient,
	}
}

func (factory *projectEntryFactory) Create(request service.TrackProjectRequest) (*ProjectEntry, error) {
	day := factory.dateTimeProvider.ToISOTime(request.Day)
	existingProject := factory.selectExistingProject()
	if existingProject != nil {
		project := existingProject.(domain.TrackedProject)
		projectEntry, err := factory.getSavedProject(day, project)
		if projectEntry == nil {
			return factory.getNewProject(day)
		}
		return projectEntry, err
	} else {
		return factory.getNewProject(day)
	}
}

func (factory *projectEntryFactory) getNewProject(day time.Time) (*ProjectEntry, error) {
	projectQuery, err := factory.prompt.Input("Enter project name")
	if err != nil {
		return nil, err
	}
	deals, err := factory.trackingClient.SearchDeals(projectQuery, day)
	if err != nil {
		return nil, err
	}
	deal, err := factory.prompt.SelectOne("Select project", deals)
	if err != nil {
		return nil, err
	}
	service, err := factory.findAndSelectService(day, deal.(*domain.Deal))
	if err != nil {
		return nil, err
	}
	return factory.createProjectEntry(deal.(*domain.Deal), service, day)
}

func (factory *projectEntryFactory) findAndSelectProject(day time.Time) (*domain.Deal, error) {
	dealQuery, err := factory.prompt.Input("Enter deal name")
	if err != nil {
		return nil, err
	}
	deals, err := factory.trackingClient.SearchDeals(dealQuery, day)
	if err != nil {
		return nil, err
	}
	deal, err := factory.prompt.SelectOne("Select project", deals)
	if err != nil {
		return nil, err
	}
	return deal.(*domain.Deal), nil
}

func (factory *projectEntryFactory) findAndSelectService(day time.Time, deal *domain.Deal) (*domain.Service, error) {
	serviceQuery, err := factory.prompt.Input("Enter service name")
	if err != nil {
		return nil, err
	}
	services, err := factory.trackingClient.SearchServices(serviceQuery, deal.ID, day)
	if err != nil {
		return nil, err
	}
	service, err := factory.prompt.SelectOne("Select service", services)
	if err != nil {
		return nil, err
	}
	return service.(*domain.Service), nil
}

func (factory *projectEntryFactory) createProjectEntry(deal *domain.Deal, service *domain.Service, day time.Time) (*ProjectEntry, error) {
	duration, err := factory.prompt.Input("Time")
	if err != nil {
		log.Error("Duration input prompt failed. Please try again.")
		return nil, err
	}
	if err := validateDurationFormat(duration); err != nil {
		log.Error("Illegal duration format. Allowed formats: HH:mm or number of minutes. Please try again.")
		return nil, err
	}
	notes, err := factory.prompt.InputMultiline("Note")
	if err != nil {
		return nil, err
	}
	projectEntry := &ProjectEntry{
		Service:  service,
		Deal:     deal,
		Day:      day,
		Duration: duration,
		Notes:    notes,
	}
	return projectEntry, nil
}

func (factory *projectEntryFactory) getSavedProject(day time.Time, project domain.TrackedProject) (*ProjectEntry, error) {
	services, err := factory.trackingClient.SearchServices(project.ServiceName, project.DealID, day)
	if err != nil {
		return nil, err
	}
	// sometimes, Productive will return multiple matches even if there's a exact match
	// so, we try and find exact match
	var exactMatchService *domain.Service
	for _, service := range services {
		if service.(*domain.Service).Name == project.ServiceName {
			exactMatchService = service.(*domain.Service)
		}
	}
	if exactMatchService == nil {
		log.Info("Cannot find exact match of saved service. Falling back to global search.")
		return nil, nil
	}
	deal := &domain.Deal{
		ID:   project.DealID,
		Name: project.DealName,
	}
	return factory.createProjectEntry(deal, exactMatchService, day)
}

func (factory *projectEntryFactory) selectExistingProject() interface{} {
	savedProjects, err := factory.projectConfigManager.GetTrackedProjects()
	if err != nil {
		return nil
	}
	if len(savedProjects) == 0 {
		return nil
	}
	selectedProject := factory.prompt.SelectOneWithSearch(
		"Select project",
		savedProjects,
		searchProjectFunction(savedProjects),
	)
	return selectedProject
}

func searchProjectFunction(projects []domain.TrackedProject) func(string, int) bool {
	return func(query string, index int) bool {
		project := projects[index]
		return strings.Contains(project.DealName, query) || strings.Contains(project.ServiceName, query)
	}
}

func validateDurationFormat(duration string) error {
	matches := utils.TimeRegex.FindStringSubmatch(duration)
	if len(matches) != 3 {
		return errors.New("time format error (only minutes or HH:mm allowed)")
	}
	return nil
}
