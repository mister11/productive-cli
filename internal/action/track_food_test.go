package action

import (
	"github.com/mister11/productive-cli/internal/client/model"
	"github.com/mister11/productive-cli/mocks"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestTrackFoodInvalid(t *testing.T) {
	client := new(mocks.TrackingClient)
	configManger := new(mocks.ConfigManager)
	dateTimeProvider := new(mocks.DateTimeProvider)

	request := TrackFoodRequest{
		IsWeekTracking: true,
		Day:            "2020-02-12",
	}

	TrackFood(client, configManger, dateTimeProvider, request)

	client.AssertNotCalled(t, mock.Anything)
	configManger.AssertNotCalled(t, mock.Anything)
	dateTimeProvider.AssertNotCalled(t, mock.Anything)
}

func TestTrackFoodWeek(t *testing.T) {
	client := new(mocks.TrackingClient)
	configManger := new(mocks.ConfigManager)
	dateTimeProvider := new(mocks.DateTimeProvider)

	request := TrackFoodRequest{
		IsWeekTracking: true,
		Day:            "",
	}

	configManger.On("GetUserID").Return("101")

	mockDateRangeAndClient(dateTimeProvider, client)

	TrackFood(client, configManger, dateTimeProvider, request)

	client.AssertExpectations(t)
	configManger.AssertExpectations(t)
	dateTimeProvider.AssertExpectations(t)
	dateTimeProvider.AssertNumberOfCalls(t, "Format", 5)
}

func TestTrackFoodDay(t *testing.T) {
	client := new(mocks.TrackingClient)
	configManger := new(mocks.ConfigManager)
	dateTimeProvider := new(mocks.DateTimeProvider)

	request := TrackFoodRequest{
		IsWeekTracking: false,
		Day:            "2020-02-20",
	}

	configManger.On("GetUserID").Return("101")

	expectedTime, _ := time.Parse("2006-01-02", "2020-02-20")
	mockServiceSearch(client, expectedTime)
	dateTimeProvider.On("ToISOTime", "2020-02-20").Return(expectedTime)
	dateTimeProvider.On("Format", expectedTime).Return("2020-02-20").Once()

	TrackFood(client, configManger, dateTimeProvider, request)

	client.AssertExpectations(t)
	configManger.AssertExpectations(t)
	dateTimeProvider.AssertExpectations(t)
}

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
	dateTimeProvider.On("Format", timeNow).Return("2020-02-20").Once()

	mockServiceSearch(client, timeNow)

	TrackFood(client, configManger, dateTimeProvider, request)

	client.AssertExpectations(t)
	configManger.AssertExpectations(t)
	dateTimeProvider.AssertExpectations(t)
}

func mockServiceSearch(client *mocks.TrackingClient, time time.Time) {
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
		On("SearchDeals", "Operations general", time).
		Return([]interface{}{deal})

	client.
		On("SearchService", "Food", "10", time).
		Return([]interface{}{service})

	timeEntry := model.NewTimeEntry(
		"", 30, "101", service, time.Format("2006-01-02"),
	)

	client.On("CreateTimeEntry", timeEntry).Return()
}

func mockDateRangeAndClient(dateTimeProvider *mocks.DateTimeProvider, client *mocks.TrackingClient) {
	dateStrings := []string {
		"2020-01-01", "2020-01-02", "2020-01-03", "2020-01-04", "2020-01-05",
	}
	var dates []time.Time
	for _, dateString := range dateStrings {
		date, _ := time.Parse("2006-01-02", dateString)
		mockServiceSearch(client, date)
		dates = append(dates, date)
		dateTimeProvider.
			On("Format", date).
			Return(dateString)
	}
	dateTimeProvider.
		On("GetWeekDays").
		Return(dates)
}
