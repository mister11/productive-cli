package service

import (
	"strings"
	"time"

	"github.com/mister11/productive-cli/internal/log"
	"github.com/mister11/productive-cli/internal/productive"
	"github.com/mister11/productive-cli/internal/service/datetime"
	"github.com/mister11/productive-cli/internal/utils"
)

type ProjectEntry struct {
	Service  *productive.Service
	Day      time.Time
	Duration string
	Notes    []string
}

type ProjectTrackingService struct {
	productiveService *ProductiveService
	prompt            Prompt
	sessionService    *SessionService
	projectStorage    ProjectStorage
	dateTimeProvider  datetime.DateTimeProvider
}

func NewProjectTrackingService(client *productive.Client) *ProjectTrackingService {
	sessionManager := NewFileUserSessionManager()
	stdInPrompt := NewStdinPrompt()
	productiveService := NewProductiveService(client)
	sessionService := NewSessionService(productiveService, stdInPrompt, sessionManager)
	dateTimeProvider := datetime.NewRealTimeDateProvider()

	return &ProjectTrackingService{
		productiveService: productiveService,
		prompt:            stdInPrompt,
		sessionService:    sessionService,
		projectStorage:    NewFileProjectStorage(),
		dateTimeProvider:  dateTimeProvider,
	}
}

func (s *ProjectTrackingService) TrackProject(request TrackProjectRequest) error {
	userSession, err := s.sessionService.ObtainUserSession()
	if err != nil {
		return err
	}
	var day time.Time
	if request.Day == "" {
		day = s.dateTimeProvider.Now()
	} else {
		day = s.dateTimeProvider.ToISOTime(request.Day)
	}
	existingProject := s.selectExistingProject()
	if existingProject != nil {
		trackedProject := s.productiveService.FindSavedProject(existingProject, day, userSession)
		if trackedProject == nil {
			return s.trackNewProject(day, userSession)
		}
		return s.trackExistingProject(trackedProject, day, userSession)
	} else {
		return s.trackNewProject(day, userSession)
	}
}

// we ignore any error here which will then cause fallback to new project search
func (s *ProjectTrackingService) selectExistingProject() *TrackedProject {
	savedProjects, err := s.projectStorage.GetTrackedProjects()
	if err != nil {
		return nil
	}
	if len(savedProjects) == 0 {
		return nil
	}
	selectedProject, err := s.prompt.SelectOneWithSearch(
		"Select project",
		savedProjects,
		searchProjectFunction(savedProjects),
	)
	if err != nil {
		return nil
	}
	return selectedProject.(*TrackedProject)
}

func (s *ProjectTrackingService) trackNewProject(day time.Time, userSession *UserSessionData) error {
	dealQuery, err := s.prompt.Input("Enter project name")
	if err != nil {
		return err
	}
	deals, err := s.productiveService.FindDeals(dealQuery, day, userSession)
	if err != nil {
		return err
	}
	selectedDeal, err := s.prompt.SelectDeal("Select project", deals)
	if err != nil {
		return err
	}
	serviceQuery, err := s.prompt.Input("Enter service name")
	if err != nil {
		return err
	}
	services, err := s.productiveService.FindServices(serviceQuery, selectedDeal, day, userSession)
	if err != nil {
		return err
	}
	selectedService, err := s.prompt.SelectService("Select service", services)
	if err != nil {
		return err
	}
	projectEntry, err := s.createProjectEntry(selectedService, day)
	if err != nil {
		return err
	}
	if err := s.productiveService.CreateProjectTimeEntry(*projectEntry, userSession); err != nil {
		return err
	}
	return s.projectStorage.UpsertTrackedProject(TrackedProject{
		DealID:      selectedDeal.ID,
		DealName:    selectedDeal.Name,
		ServiceID:   selectedService.ID,
		ServiceName: selectedService.Name,
	})
}

func (s *ProjectTrackingService) trackExistingProject(project *TrackedProject, day time.Time, userSession *UserSessionData) error {
	service := &productive.Service{
		ID:   project.ServiceID,
		Name: project.ServiceName,
	}
	projectEntry, err := s.createProjectEntry(service, day)
	if err != nil {
		return err
	}
	return s.productiveService.CreateProjectTimeEntry(*projectEntry, userSession)
}

func (s *ProjectTrackingService) createProjectEntry(
	service *productive.Service,
	day time.Time,
) (*ProjectEntry, error) {
	duration, err := s.prompt.Input("Time")
	for err != nil {
		log.Info("Duration input prompt failed. Please try again.")
		duration, err = s.prompt.Input("Time")
	}
	durationParsed, err := utils.ParseTime(duration)
	for err != nil {
		log.Info("Illegal duration format. Allowed formats: HH:mm or number of minutes. Please try again.")
		duration, err = s.prompt.Input("Time")
		if err != nil {
			log.Info("Duration input prompt failed. Please try again.")
		} else {
			durationParsed, err = utils.ParseTime(duration)
		}
	}
	notes, err := s.prompt.InputMultiline("Note")
	for err != nil {
		log.Info("There's a problem with note input. Please try again")
		notes, err = s.prompt.InputMultiline("Note")
	}
	projectEntry := &ProjectEntry{
		Service:  service,
		Day:      day,
		Duration: durationParsed,
		Notes:    notes,
	}
	return projectEntry, nil
}

func searchProjectFunction(projects []TrackedProject) func(string, int) bool {
	return func(query string, index int) bool {
		project := projects[index]
		return strings.Contains(project.DealName, query) || strings.Contains(project.ServiceName, query)
	}
}