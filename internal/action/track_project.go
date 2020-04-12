package action

import (
	"github.com/mister11/productive-cli/internal/utils"
	"strings"
	"time"

	"github.com/mister11/productive-cli/internal/client"
	"github.com/mister11/productive-cli/internal/client/model"
	"github.com/mister11/productive-cli/internal/config"
	"github.com/mister11/productive-cli/internal/datetime"
	"github.com/mister11/productive-cli/internal/stdin"
)

func TrackProject(
	productiveClient client.TrackingClient,
	stdin stdin.Stdin,
	configManager config.ConfigManager,
	dateTimeProvider datetime.DateTimeProvider,
	trackProjectRequest TrackProjectRequest,
) {
	existingProject := selectExistingProject(stdin, configManager)
	var date time.Time
	if trackProjectRequest.Day != "" {
		date = dateTimeProvider.ToISOTime(trackProjectRequest.Day)
	} else {
		date = dateTimeProvider.Now()
	}
	if existingProject != nil {
		project := existingProject.(config.Project)
		trackSavedProject(productiveClient, stdin, configManager, dateTimeProvider, project, date)
	} else {
		trackNewProject(productiveClient, stdin, configManager, dateTimeProvider, date)
	}
}

func trackSavedProject(
	productiveClient client.TrackingClient,
	stdin stdin.Stdin,
	configManager config.ConfigManager,
	dateTimeProvider datetime.DateTimeProvider,
	project config.Project,
	date time.Time,
) {
	configManager.RemoveExistingProject(project)
	deal, service := findProjectInfo(productiveClient, project, date)
	duration := utils.ParseTime(stdin.Input("Time"))
	notes := createNotes(stdin)
	timeEntry := model.NewTimeEntry(notes, duration, configManager.GetUserID(), service, dateTimeProvider.Format(date))
	productiveClient.CreateTimeEntry(timeEntry)
	configManager.SaveProject(config.NewProject(*deal, *service))
}

func trackNewProject(
	productiveClient client.TrackingClient,
	stdin stdin.Stdin,
	configManager config.ConfigManager,
	dateTimeProvider datetime.DateTimeProvider,
	date time.Time,
) {
	selectedDeal := searchNewDeal(productiveClient, stdin, date)
	selectedService := searchNewService(productiveClient, stdin, selectedDeal, date)

	duration := utils.ParseTime(stdin.Input("Time"))
	notes := createNotes(stdin)
	timeEntry := model.NewTimeEntry(notes, duration, configManager.GetUserID(), selectedService, dateTimeProvider.Format(date))
	productiveClient.CreateTimeEntry(timeEntry)

	configManager.SaveProject(config.NewProject(*selectedDeal, *selectedService))
}

func selectExistingProject(stdin stdin.Stdin, configManager config.ConfigManager) interface{} {
	savedProjects := configManager.GetSavedProjects()
	selectedProject := stdin.SelectOneWithSearch(
		"Select project",
		savedProjects,
		searchProjectFunction(savedProjects),
	)
	return selectedProject
}

func findProjectInfo(productiveClient client.TrackingClient, existingProject config.Project, day time.Time) (*model.Deal, *model.Service) {
	deals := productiveClient.SearchDeals(existingProject.DealName, day)
	deal := deals[0].(*model.Deal)
	services := productiveClient.SearchService(existingProject.ServiceName, deal.ID, day)
	service := services[0].(*model.Service)
	return deal, service
}

func searchNewDeal(productiveClient client.TrackingClient, stdin stdin.Stdin, day time.Time) *model.Deal {
	dealQuery := stdin.Input("Search project")
	deals := productiveClient.SearchDeals(dealQuery, day)
	return stdin.SelectOne("Select project", deals).(*model.Deal)
}

func searchNewService(productiveClient client.TrackingClient, stdin stdin.Stdin, deal *model.Deal, day time.Time) *model.Service {
	serviceQuery := stdin.Input("Search service")
	services := productiveClient.SearchService(serviceQuery, deal.ID, day)
	return stdin.SelectOne("Select service", services).(*model.Service)
}

func createNotes(stdin stdin.Stdin) string {
	notes := stdin.InputMultiple("Note")
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
