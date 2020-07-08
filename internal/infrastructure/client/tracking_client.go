package client

import (
	"github.com/mister11/productive-cli/internal/domain/tracking"
	"time"
)

type TrackingClient interface {
	TrackFood(entries []tracking.FoodEntry) error
	TrackProject(entry *tracking.ProjectEntry) error
	VerifyLogin() (string, error)
	Login(username string, password string) error
	SearchDeals(query string, day time.Time) []interface{}
	SearchServices(query string, dealID string, day time.Time) []interface{}
}
