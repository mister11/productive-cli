package tracking

import (
	"github.com/mister11/productive-cli/internal/domain/config"
	"github.com/mister11/productive-cli/internal/domain/input"
	"strings"
	"time"

	"github.com/mister11/productive-cli/internal/domain/datetime"
	"github.com/mister11/productive-cli/internal/infrastructure/client"
	"github.com/mister11/productive-cli/internal/infrastructure/client/model"
	"github.com/mister11/productive-cli/internal/infrastructure/log"
	"github.com/mister11/productive-cli/internal/utils"
)

type ProjectTracker interface {
	TrackProject(day time.Time)
}

type httpProjectTracker struct {
	client           client.TrackingClient
	stdIn            input.Prompt
	config           config.Manager
	dateTimeProvider datetime.DateTimeProvider
}

func NewHTTPProjectTracker(
	client client.TrackingClient,
	stdIn input.Prompt,
	config config.Manager,
	dateTimeProvider datetime.DateTimeProvider,
) *httpProjectTracker {
	return &httpProjectTracker{
		client:           client,
		stdIn:            stdIn,
		config:           config,
		dateTimeProvider: dateTimeProvider,
	}
}

func (tracker *httpProjectTracker) TrackProject(day time.Time) {
	existingProject := tracker.selectExistingProject()
	if existingProject != nil {
		project := existingProject.(config.Project)
		tracker.trackSavedProject(day, project)
	} else {
		tracker.trackNewProject(day)
	}
}

func (tracker *httpProjectTracker) trackNewProject(day time.Time) {
	deal := tracker.findAndSelectDeal(day)
	service := tracker.findAndSelectService(day, deal)
	tracker.trackSelectedProject(deal, service, day)
	tracker.config.SaveProject(config.NewProject(*deal, *service))
}

func (tracker *httpProjectTracker) trackSavedProject(day time.Time, project config.Project) {
	tracker.config.RemoveExistingProject(project)
	deal, service := tracker.client.FindProjectInfo(project.DealName, project.ServiceName, day)
	tracker.trackSelectedProject(deal, service, day)
	tracker.config.SaveProject(config.NewProject(*deal, *service))
}

func (tracker *httpProjectTracker) findAndSelectDeal(day time.Time) *model.Deal {
	dealQuery, _ := tracker.stdIn.Input("Search project")
	deals := tracker.client.SearchDeals(dealQuery, day)
	return tracker.stdIn.SelectOne("Select project", deals).(*model.Deal)
}

func (tracker *httpProjectTracker) findAndSelectService(day time.Time, deal *model.Deal) *model.Service {
	serviceQuery, _ := tracker.stdIn.Input("Search service")
	services := tracker.client.SearchServices(serviceQuery, deal.ID, day)
	return tracker.stdIn.SelectOne("Select service", services).(*model.Service)
}

func (tracker *httpProjectTracker) trackSelectedProject(deal *model.Deal, service *model.Service, day time.Time) {
	duration, err := tracker.stdIn.Input("Time")
	if err != nil {
		log.Error("Duration input prompt failed. Please try again.")
		tracker.trackSelectedProject(deal, service, day)
	}
	if err := verifyDuration(duration); err != nil {
		log.Error("Illegal duration format. Allowed formats: HH:mm or number of minutes. Please try again.")
		tracker.trackSelectedProject(deal, service, day)
	}
	notes := tracker.stdIn.InputMultiline("Note")
	notesFormatted := createNotes(notes)
	log.Info("Tracking %s - %s with duration %d for %d", deal.Name, service.Name, duration, utils.FormatDate(day))
	tracker.client.CreateProjectTimeEntry(service, day, duration, notesFormatted, tracker.config.GetUserID())
}

func (tracker *httpProjectTracker) selectExistingProject() interface{} {
	savedProjects := tracker.config.GetSavedProjects()
	if len(savedProjects) == 0 {
		return nil
	}
	selectedProject := tracker.stdIn.SelectOneWithSearch(
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