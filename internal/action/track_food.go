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
		trackFood(productiveClient, configManager, dateTimeProvider, dateTimeProvider.GetWeekDays()...)
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
		dayFormatted := dateTimeProvider.Format(day)
		log.Info("Tracking food for " + dayFormatted)
		service := findFoodService(productiveClient, dayFormatted)
		timeEntry := model.NewTimeEntry("", 30, userID, service, dayFormatted)
		productiveClient.CreateTimeEntry(timeEntry)
	}
}

func findFoodService(productiveClient client.TrackingClient, dayFormatted string) *model.Service {
	deal := productiveClient.SearchDeals("Operations general", dayFormatted)[0].(*model.Deal)
	service := productiveClient.SearchService("Food", deal.ID, dayFormatted)[0].(*model.Service)
	return service
}
