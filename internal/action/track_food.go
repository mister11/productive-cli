package action

import (
	"time"

	"github.com/mister11/productive-cli/internal/client"
	"github.com/mister11/productive-cli/internal/client/model"
	"github.com/mister11/productive-cli/internal/config"
	"github.com/mister11/productive-cli/internal/datetime"
	"github.com/mister11/productive-cli/internal/log"
)

func TrackFood(
	productiveClient client.TrackingClient,
	configManager config.ConfigManager,
	dateTimeProvider datetime.DateTimeProvider,
	trackFoodRequest TrackFoodRequest,
) {
	if !trackFoodRequest.IsValid() {
		log.Error("You've provided both week and day tracking so I don't know what to do.")
		return
	}

	if trackFoodRequest.IsWeekTracking {
		trackFood(productiveClient, configManager, dateTimeProvider, getWeekDays(dateTimeProvider)...)
	} else if trackFoodRequest.Day != "" {
		date := dateTimeProvider.ToISOTime(trackFoodRequest.Day)
		trackFood(productiveClient, configManager, dateTimeProvider, []time.Time{date}...)
	} else {
		trackFood(productiveClient, configManager, dateTimeProvider, []time.Time{dateTimeProvider.Now()}...)
	}
}

func trackFood(
	productiveClient client.TrackingClient,
	configManager config.ConfigManager,
	dateTimeProvider datetime.DateTimeProvider,
	days ...time.Time,
) {
	userID := configManager.GetUserID()
	for _, day := range days {
		log.Info("Tracking food for " + dateTimeProvider.Format(day))
		service := findFoodService(productiveClient, day)
		timeEntry := model.NewTimeEntry("", 30, userID, service, dateTimeProvider.Format(day))
		productiveClient.CreateTimeEntry(timeEntry)
	}
}

func findFoodService(productiveClient client.TrackingClient, day time.Time) *model.Service {
	deal := productiveClient.SearchDeals("Operations general", day)[0].(*model.Deal)
	service := productiveClient.SearchService("Food", deal.ID, day)[0].(*model.Service)
	return service
}

func getWeekDays(dateTimeProvider datetime.DateTimeProvider) []time.Time {
	var days []time.Time

	start := dateTimeProvider.WeekStart()
	end := dateTimeProvider.WeekEnd()
	for day := start; day.Day() <= end.Day(); day = day.AddDate(0, 0, 1) {
		days = append(days, day)
	}
	return days
}
