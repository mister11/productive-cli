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
	var day time.Time
	if trackProjectRequest.Day != "" {
		day = dateTimeProvider.ToISOTime(trackProjectRequest.Day)
	} else {
		day = dateTimeProvider.Now()
	}
	dayFormatted := dateTimeProvider.Format(day)
	if existingProject != nil {
		project := existingProject.(config.Project)
		trackSavedProject(productiveClient, stdin, configManager, project, dayFormatted)
	} else {
		trackNewProject(productiveClient, stdin, configManager, dayFormatted)
	}
}

func trackSavedProject(
	productiveClient client.TrackingClient,
	stdin stdin.Stdin,
	configManager config.ConfigManager,
	project config.Project,
	dayFormatted string,
) {
	configManager.RemoveExistingProject(project)
	deal, service := findProjectInfo(productiveClient, project, dayFormatted)
	duration := utils.ParseTime(stdin.Input("Time"))
	notes := createNotes(stdin)
	timeEntry := model.NewTimeEntry(notes, duration, configManager.GetUserID(), service, dayFormatted)
	productiveClient.CreateTimeEntry(timeEntry)
	configManager.SaveProject(config.NewProject(*deal, *service))
}

func trackNewProject(
	productiveClient client.TrackingClient,
	stdin stdin.Stdin,
	configManager config.ConfigManager,
	dayFormatted string,
) {
	selectedDeal := searchNewDeal(productiveClient, stdin, dayFormatted)
	selectedService := searchNewService(productiveClient, stdin, selectedDeal, dayFormatted)

	duration := utils.ParseTime(stdin.Input("Time"))
	notes := createNotes(stdin)
	timeEntry := model.NewTimeEntry(notes, duration, configManager.GetUserID(), selectedService, dayFormatted)
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

func findProjectInfo(productiveClient client.TrackingClient, existingProject config.Project, dayFormatted string) (*model.Deal, *model.Service) {
	deals := productiveClient.SearchDeals(existingProject.DealName, dayFormatted)
	deal := deals[0].(*model.Deal)
	services := productiveClient.SearchService(existingProject.ServiceName, deal.ID, dayFormatted)
	service := services[0].(*model.Service)
	return deal, service
}

func searchNewDeal(productiveClient client.TrackingClient, stdin stdin.Stdin, dayFormatted string) *model.Deal {
	dealQuery := stdin.Input("Search project")
	deals := productiveClient.SearchDeals(dealQuery, dayFormatted)
	return stdin.SelectOne("Select project", deals).(*model.Deal)
}

func searchNewService(productiveClient client.TrackingClient, stdin stdin.Stdin, deal *model.Deal, dayFormatted string) *model.Service {
	serviceQuery := stdin.Input("Search service")
	services := productiveClient.SearchService(serviceQuery, deal.ID, dayFormatted)
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
