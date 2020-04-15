package track

import (
	"github.com/mister11/productive-cli/internal/action"
	"github.com/mister11/productive-cli/internal/client/model"
	"github.com/mister11/productive-cli/internal/config"
	"github.com/mister11/productive-cli/mocks"
	"testing"
	"time"
)

func TestTrackProjectWithDay(t *testing.T) {
	client := new(mocks.TrackingClient)
	stdIn := new(mocks.Stdin)
	configManger := new(mocks.ConfigManager)
	dateTimeProvider := new(mocks.DateTimeProvider)

	trackingManger := NewTrackingManager(client, stdIn, configManger, dateTimeProvider)

	day := "2020-02-20"
	dayTime, _ := time.Parse("2016-01-02", day)

	// no saved projects
	configManger.On("GetSavedProjects").Return(nil).Once()

	dateTimeProvider.On("ToISOTime", day).Return(dayTime).Once()
	dateTimeProvider.On("Format", dayTime).Return(day).Once()

	deals := []*model.Deal {
		{
			ID:      "0",
			Name:    "Deal 1",
			EndDate: "2021-01-21",
		},
		{
			ID:      "1",
			Name:    "Deal 2",
			EndDate: "2021-03-21",
		},
	}
	stdIn.On("Input", "Search project").Return("dealQuery").Once()
	// Go in fucking amazing...
	// Good thing we have generics so we can extract that to a helper function... Oh, wait...
	var g []interface{}
	for _, d := range deals {
		g = append(g, d)
	}
	client.On("SearchDeals", "dealQuery", day).Return(g).Once()
	stdIn.On("SelectOne", "Select project", g).Return(g[0]).Once()

	services := []*model.Service {
		{
			ID:   "10",
			Name: "Service 1",
		},
		{
			ID:      "11",
			Name:    "Service 2",
		},
	}
	stdIn.On("Input", "Search service").Return("serviceQuery").Once()
	// Go in fucking amazing...
	// Good thing we have generics so we can extract that to a helper function... Oh, wait...
	var h []interface{}
	for _, s := range services {
		h = append(h, s)
	}
	client.On("SearchServices", "serviceQuery", "0", day).Return(h).Once()
	stdIn.On("SelectOne", "Select service", h).Return(h[1]).Once()

	stdIn.On("Input", "Time").Return("8:00").Once()
	stdIn.On("InputMultiple", "Note").Return([]string {"Task 1", "Task 2"}).Once()
	configManger.On("GetUserID").Return("69").Once()
	expectedTimeEntry := &model.TimeEntry{
		ID:     "0",
		Date:   "2020-02-20",
		Note:   getFormattedNotes(),
		Time:   "480",
		User:   &model.Person{ID: "69"},
		Budget: services[1],
	}
	expectedConfigProject := config.Project{
		DealID:      "0",
		DealName:    "Deal 1",
		ServiceID:   "11",
		ServiceName: "Service 2",
	}

	client.On("CreateTimeEntry", expectedTimeEntry).Return().Once()
	configManger.On("SaveProject", expectedConfigProject).Return().Once()
	
	request := action.TrackProjectRequest{Day: day}
	trackingManger.TrackProject(request)

	client.AssertExpectations(t)
	stdIn.AssertExpectations(t)
	configManger.AssertExpectations(t)
	dateTimeProvider.AssertExpectations(t)
}

func getFormattedNotes() string {
	return "<ul>" +
		"<li>Task 1</li>" +
		"<li>Task 2</li>" +
		"</ul>"
}