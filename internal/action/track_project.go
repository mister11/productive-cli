package action

import (
	"github.com/mister11/productive-cli/internal/utils"
	"strconv"
	"strings"
	"time"

	"github.com/mister11/productive-cli/internal/client"
	"github.com/mister11/productive-cli/internal/client/model"
	"github.com/mister11/productive-cli/internal/config"
	"github.com/mister11/productive-cli/internal/datetime"
	"github.com/mister11/productive-cli/internal/prompt"
)

func TrackProject(productiveClient client.ProductiveClient, trackProjectRequest TrackProjectRequest) {
	existingProject := selectExistingProject()
	var date time.Time
	if trackProjectRequest.Day != "" {
		date = datetime.ToISODate(trackProjectRequest.Day)
	} else {
		date = datetime.Now()
	}
	if existingProject != nil {
		project := existingProject.(config.Project)
		trackSavedProject(productiveClient, project, date)
	} else {
		trackNewProject(productiveClient, date)
	}
}

func trackSavedProject(productiveClient client.ProductiveClient, project config.Project, date time.Time) {
	config.RemoveExistingProject(project)
	deal, service := findProjectInfo(productiveClient, project, date)
	duration := utils.ParseTime(prompt.Input("Time"))
	notes := createNotes()
	timeEntry := model.NewTimeEntry(notes, duration, config.GetUserID(), service, date)
	productiveClient.CreateTimeEntry(timeEntry)
	config.SaveProjectToConfig(config.NewProject(*deal, *service))
}

func trackNewProject(productiveClient client.ProductiveClient, date time.Time) {
	selectedDeal := searchNewDeal(productiveClient, date)
	selectedService := searchNewService(productiveClient, selectedDeal, date)

	duration := utils.ParseTime(prompt.Input("Time"))
	notes := createNotes()
	timeEntry := model.NewTimeEntry(notes, duration, config.GetUserID(), selectedService, date)
	productiveClient.CreateTimeEntry(timeEntry)

	config.SaveProjectToConfig(config.NewProject(*selectedDeal, *selectedService))
}

func selectExistingProject() interface{} {
	savedProjects := config.GetSavedProjects()
	selectedProject := prompt.SelectOneWithSearch(
		"Select project",
		savedProjects,
		searchProjectFunction(savedProjects),
	)
	return selectedProject
}

func findProjectInfo(productiveClient client.ProductiveClient, existingProject config.Project, day time.Time) (*model.Deal, *model.Service) {
	deals := productiveClient.SearchDeals(existingProject.DealName, day)
	deal := deals[0].(*model.Deal)
	services := productiveClient.SearchService(existingProject.ServiceName, deal.ID, day)
	service := services[0].(*model.Service)
	return deal, service
}

func searchNewDeal(productiveClient client.ProductiveClient, day time.Time) *model.Deal {
	dealQuery := prompt.Input("Search project")
	deals := productiveClient.SearchDeals(dealQuery, day)
	return prompt.SelectOne("Select project", deals).(*model.Deal)
}

func searchNewService(productiveClient client.ProductiveClient, deal *model.Deal, day time.Time) *model.Service {
	serviceQuery := prompt.Input("Search service")
	services := productiveClient.SearchService(serviceQuery, deal.ID, day)
	return prompt.SelectOne("Select service", services).(*model.Service)
}

func createNotes() string {
	index := 1
	var notes []string
	for isEnd := false; !isEnd; {
		note := prompt.Input("Enter note " + strconv.Itoa(index) + " (empty to finish)")
		if len(note) == 0 {
			isEnd = true
			continue
		}
		notes = append(notes, note)
		index++
	}
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
