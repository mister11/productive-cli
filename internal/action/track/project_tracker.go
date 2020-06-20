package track

import (
	"strings"
	"time"

	"github.com/mister11/productive-cli/internal/client"
	"github.com/mister11/productive-cli/internal/client/model"
	"github.com/mister11/productive-cli/internal/config"
	"github.com/mister11/productive-cli/internal/datetime"
	"github.com/mister11/productive-cli/internal/log"
	"github.com/mister11/productive-cli/internal/stdin"
	"github.com/mister11/productive-cli/internal/utils"
)

type projectTracker interface {
	TrackProject(day time.Time, project *config.Project)
}

type httpProjectTracker struct {
	client           client.TrackingClient
	stdIn            stdin.Stdin
	config           config.ConfigManager
	dateTimeProvider datetime.DateTimeProvider
}

func newHTTPProjectTracker(
	client           client.TrackingClient,
	stdIn            stdin.Stdin,
	config           config.ConfigManager,
	dateTimeProvider datetime.DateTimeProvider,
) *httpProjectTracker {
	return &httpProjectTracker{
		client:           client,
		stdIn:            stdIn,
		config:           config,
		dateTimeProvider: dateTimeProvider,
	}
}

func (tracker *httpProjectTracker) TrackProject(day time.Time, project *config.Project) {
	if project == nil {
		tracker.trackNewProject(day)
	} else {
		tracker.trackSavedProject(day, *project)
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
	dealQuery := tracker.stdIn.Input("Search project")
	deals := tracker.client.SearchDeals(dealQuery, day)
	return tracker.stdIn.SelectOne("Select project", deals).(*model.Deal)
}

func (tracker *httpProjectTracker) findAndSelectService(day time.Time, deal *model.Deal) *model.Service {
	serviceQuery := tracker.stdIn.Input("Search service")
	services := tracker.client.SearchServices(serviceQuery, deal.ID, day)
	return tracker.stdIn.SelectOne("Select service", services).(*model.Service)
}

func (tracker *httpProjectTracker) trackSelectedProject(deal *model.Deal, service *model.Service, day time.Time) {
	duration := tracker.stdIn.Input("Time")
	durationParsed := utils.ParseTime(duration)
	notes := tracker.stdIn.InputMultiple("Note")
	notesFormatted := createNotes(notes)
	log.Info("Tracking %s - %s with duration %d for %d", deal.Name, service.Name, duration, utils.FormatDate(day))
	tracker.client.CreateProjectTimeEntry(service, day, durationParsed, notesFormatted, tracker.config.GetUserID())
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