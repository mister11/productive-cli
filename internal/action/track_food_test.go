package action

import (
	"github.com/mister11/productive-cli/internal/client/model"
	"github.com/mister11/productive-cli/mocks"
	"testing"
	"time"
)

func TestTrackFood(t *testing.T) {
	client := new(mocks.TrackingClient)
	configManger := new(mocks.ConfigManager)
	dateTimeProvider := new(mocks.DateTimeProvider)

	request := TrackFoodRequest{
		IsWeekTracking: false,
		Day:            "",
	}

	configManger.On("GetUserID").Return("101")

	timeNow, _ := time.Parse("2006-01-02", "2020-02-20")
	dateTimeProvider.On("Now").Return(timeNow)
	dateTimeProvider.On("Format", timeNow).Return("2020-02-20")

	deal := &model.Deal{
		ID:      "10",
		Name:    "Deal 1",
		EndDate: "2020-03-13",
	}

	service := &model.Service{
		ID:   "20",
		Name: "Service 1",
	}

	client.
		On("SearchDeals", "Operations general", timeNow).
		Return([]interface{}{deal})

	client.
		On("SearchService", "Food", "10", timeNow).
		Return([]interface{}{service})

	timeEntry := model.NewTimeEntry(
		"", 30, "101", service, "2020-02-20",
	)

	client.On("CreateTimeEntry", timeEntry).Return()

	TrackFood(client, configManger, dateTimeProvider, request)

	client.AssertExpectations(t)
	configManger.AssertExpectations(t)
	dateTimeProvider.AssertExpectations(t)
}
