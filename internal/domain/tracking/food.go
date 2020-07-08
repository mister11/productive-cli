package tracking

import (
	"errors"
	"github.com/mister11/productive-cli/internal/domain/config"
	"github.com/mister11/productive-cli/internal/domain/datetime"
	"time"
)

type FoodEntry struct {
	day    time.Time
	userID string
}

type FoodEntriesCreator interface {
	Create(request TrackFoodRequest) ([]FoodEntry, error)
}

type foodEntriesFactory struct {
	dateTimeProvider datetime.DateTimeProvider
	config           config.Manager
}

func NewFoodEntriesCreator(
	provider datetime.DateTimeProvider,
	manager config.Manager,
) FoodEntriesCreator {
	return &foodEntriesFactory{
		dateTimeProvider: provider,
		config:           manager,
	}
}

func (factory *foodEntriesFactory) Create(trackFoodRequest TrackFoodRequest) ([]FoodEntry, error) {
	if !trackFoodRequest.IsValid() {
		return nil, errors.New("invalid track food request")
	}
	days := factory.getTrackingDays(trackFoodRequest)
	userID := factory.config.GetUserID()
	var entries []FoodEntry
	for _, day := range days {
		entry := FoodEntry{
			day:    day,
			userID: userID,
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

func (factory *foodEntriesFactory) getTrackingDays(request TrackFoodRequest) []time.Time {
	if request.IsWeekTracking {
		return factory.dateTimeProvider.GetWeekDays()
	} else if request.Day != "" {
		return []time.Time{factory.dateTimeProvider.ToISOTime(request.Day)}
	} else {
		return []time.Time{factory.dateTimeProvider.Now()}
	}
}
