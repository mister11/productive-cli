package track

import (
	"time"

	"github.com/mister11/productive-cli/internal/client"
	"github.com/mister11/productive-cli/internal/config"
)

type foodTracker interface {
	TrackFood(days ...time.Time)
}

type httpFoodTracker struct {
	client client.TrackingClient
	config config.ConfigManager
}

func newHTTPFoodTracker(
	client client.TrackingClient,
	config config.ConfigManager,
) *httpFoodTracker {
	return &httpFoodTracker{
		client: client,
		config: config,
	}
}

func (tracker *httpFoodTracker) TrackFood(days ...time.Time) {
	userID := tracker.config.GetUserID()
	for _, day := range days {
		tracker.client.CreateFoodTimeEntry(day, userID)
	}
}
