package action

import (
	"time"

	"gitlab.com/mister11/productive-cli/internal/client"
	"gitlab.com/mister11/productive-cli/internal/client/model"
	"gitlab.com/mister11/productive-cli/internal/config"
	"gitlab.com/mister11/productive-cli/internal/datetime"
	"gitlab.com/mister11/productive-cli/internal/log"
)

func TrackFood(productiveClient client.ProductiveClient, trackFoodRequest TrackFoodRequest) {
	if !trackFoodRequest.IsValid() {
		log.Error("You've provided both week and day tracking so I don't know what to do.")
		return
	}

	if trackFoodRequest.IsWeekTracking {
		trackFood(productiveClient, getWeekDays()...)
	} else if trackFoodRequest.Day != "" {
		date := datetime.ToISODate(trackFoodRequest.Day)
		trackFood(productiveClient, []time.Time{date}...)
	} else {
		trackFood(productiveClient, []time.Time{datetime.Now()}...)
	}
}

func trackFood(productiveClient client.ProductiveClient, days ...time.Time) {
	userID := config.GetUserID()
	for _, day := range days {
		log.Info("Tracking food for " + datetime.Format(day))
		service := findFoodService(productiveClient, day)
		timeEntry := model.NewTimeEntry("", 30, userID, service, day)
		productiveClient.CreateFoodTimeEntry(timeEntry)
	}
}

func findFoodService(productiveClient client.ProductiveClient, day time.Time) *model.Service {
	deal := productiveClient.SearchDeals("Operations general", day)[0].(*model.Deal)
	service := productiveClient.SearchService("Food", deal.ID, day)[0].(*model.Service)
	return service
}

func getWeekDays() []time.Time {
	var days []time.Time

	start := datetime.WeekStart()
	end := datetime.WeekEnd()
	for day := start; day.Day() <= end.Day(); day = day.AddDate(0, 0, 1) {
		days = append(days, day)
	}
	return days
}
