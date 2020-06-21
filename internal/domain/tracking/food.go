package tracking

import (
	config2 "github.com/mister11/productive-cli/internal/domain/config"
	"github.com/mister11/productive-cli/internal/domain/datetime"
	"github.com/mister11/productive-cli/internal/infrastructure/client"
	"github.com/mister11/productive-cli/internal/infrastructure/log"
	"time"
)

type FoodTracker interface {
	TrackFood(trackFoodRequest TrackFoodRequest)
}

type httpFoodTracker struct {
	client           client.TrackingClient
	config           config2.Manager
	dateTimeProvider datetime.DateTimeProvider
}

func NewHTTPFoodTracker(
	client client.TrackingClient,
	config config2.Manager,
	dateTimeProvider datetime.DateTimeProvider,
) *httpFoodTracker {
	return &httpFoodTracker{
		client:           client,
		config:           config,
		dateTimeProvider: dateTimeProvider,
	}
}

func (tracker *httpFoodTracker) TrackFood(trackFoodRequest TrackFoodRequest) {
	if !trackFoodRequest.IsValid() {
		log.Error("You've provided both week and day tracking so I don't know what to do.", nil)
		return
	}
	days := tracker.getTrackingDays(trackFoodRequest)
	userID := tracker.config.GetUserID()
	for _, day := range days {
		tracker.client.CreateFoodTimeEntry(day, userID)
	}
}

func (tracker *httpFoodTracker) getTrackingDays(request TrackFoodRequest) []time.Time {
	if request.IsWeekTracking {
		return tracker.dateTimeProvider.GetWeekDays()
	} else if request.Day != "" {
		return []time.Time{tracker.dateTimeProvider.ToISOTime(request.Day)}
	} else {
		return []time.Time{tracker.dateTimeProvider.Now()}
	}
}
