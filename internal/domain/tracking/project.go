package tracking

import (
	"errors"
	"github.com/mister11/productive-cli/internal/domain/config"
	"github.com/mister11/productive-cli/internal/domain/input"
	"strings"
	"time"

	"github.com/mister11/productive-cli/internal/domain/datetime"
	"github.com/mister11/productive-cli/internal/infrastructure/client/model"
	"github.com/mister11/productive-cli/internal/infrastructure/log"
	"github.com/mister11/productive-cli/internal/utils"
)

type ProjectEntry struct {
	service  *model.Service
	day      time.Time
	duration string
	notes    []string
	userID   string
}

type ProjectEntryCreator interface {
	Create(request TrackProjectRequest) (*ProjectEntry, error)
}

type projectEntryFactory struct {
	dateTimeProvider datetime.DateTimeProvider
	prompt           input.Prompt
	config           config.Manager
	searcher         Searcher
}

func NewProjectEntryCreator(
	dateTimeProvider datetime.DateTimeProvider,
	prompt input.Prompt,
	config config.Manager,
	searcher Searcher,
) ProjectEntryCreator {
	return &projectEntryFactory{
		dateTimeProvider: dateTimeProvider,
		prompt:           prompt,
		config:           config,
		searcher:         searcher,
	}
}

func (factory *projectEntryFactory) Create(request TrackProjectRequest) (*ProjectEntry, error) {
	day := factory.dateTimeProvider.ToISOTime(request.Day)
	existingProject := factory.selectExistingProject()
	if existingProject != nil {
		project := existingProject.(config.Project)
		return factory.getSavedProject(day, project)
	} else {
		return factory.getNewProject(day)
	}
}

func (factory *projectEntryFactory) getNewProject(day time.Time) (*ProjectEntry, error) {
	deal, err := factory.findAndSelectProject(day)
	if err != nil {
		return nil, err
	}
	service, err := factory.findAndSelectService(day, deal)
	if err != nil {
		return nil, err
	}
	return factory.createProjectEntry(service, day)
}

func (factory *projectEntryFactory) findAndSelectProject(day time.Time) (*model.Deal, error) {
	project, err := factory.searcher.SearchProjects(day)
	if err != nil {
		log.Error("Project search failed. Please try again.")
		return nil, err
	}
	return project, nil
}

func (factory *projectEntryFactory) findAndSelectService(day time.Time, deal *model.Deal) (*model.Service, error) {
	service, err := factory.searcher.SearchServices(day, deal)
	if err != nil {
		log.Error("Service search failed. Please try again.")
		return nil, err
	}
	return service, nil
}

func (factory *projectEntryFactory) createProjectEntry(service *model.Service, day time.Time) (*ProjectEntry, error) {
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
	//notesFormatted := createNotes(notes)
	projectEntry := &ProjectEntry{
		service:  service,
		day:      day,
		duration: duration,
		notes:    notes,
		userID:   factory.config.GetUserID(),
	}
	return projectEntry, nil
}

func (factory *projectEntryFactory) getSavedProject(day time.Time, project config.Project) (*ProjectEntry, error) {
	factory.config.RemoveExistingProject(project)
	service, err := factory.searcher.SearchService(project.ServiceName, project.DealID, day)
	if err != nil {
		return nil, err
	}
	return factory.createProjectEntry(service, day)
}

func (factory *projectEntryFactory) selectExistingProject() interface{} {
	savedProjects := factory.config.GetSavedProjects()
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

func searchProjectFunction(projects []config.Project) func(string, int) bool {
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

func createNotes(notes []string) string {
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
